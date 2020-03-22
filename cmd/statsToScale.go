package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	lFile "github.com/snowzach/rotatefilehook"
	"log/syslog"
	"os"
	"os/signal"
	"stats-to-scalse/pkg/stats"
	"strconv"
	"syscall"
)

func main() {
	setupLogging()
	log.Infof("===> Stats-to-scale <===")

	// init collector
	collector = stats.NewGopsutilsCollector()

	setupWriter()

	// start cron job
	cronManager := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(&cronLogrusLogger{})),
	)
	_, err := cronManager.AddFunc("@every 3s", task)
	if err != nil {
		log.WithField("stackTrace", fmt.Sprintf("%+v", err)).Fatal(err)
	}
	cronManager.Start()
	log.Infof("Cron manager started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Infof("Term signals connected")
	<-sigs
	log.Infof("Terminate signal received")
	cronManager.Stop()
	log.Infof("Cron manager stopped")
}

func setupWriter() {
	connectionString := os.Getenv("STS_DATABASE_CONNECTION_STRING")
	pgWriter, err := stats.NewPgWriter(connectionString)
	if err != nil {
		log.WithField("stackTrace", fmt.Sprintf("%+v", err)).Fatal(err)
	}
	writer = pgWriter
}

func setupLogging() {
	log.SetOutput(os.Stdout)
	setupLogLevel()
	setupLoggingToSyslog()
	setupLoggingToFile()
}

func setupLoggingToSyslog() {
	if os.Getenv("STS_LOGGING_TO_SYSLOG_ENABLED") != "true" {
		return
	}
	hook, err := lSyslog.NewSyslogHook(
		os.Getenv("STS_SYSLOG_PROTOCOL"),
		os.Getenv("STS_SYSLOG_ADDRESS"),
		syslog.LOG_INFO,
		"",
	)
	if err != nil {
		log.Fatal(err)
	}
	log.AddHook(hook)
}

func setupLoggingToFile() {
	if os.Getenv("STS_LOGGING_TO_FILE_ENABLED") != "true" {
		return
	}
	hook, err := lFile.NewRotateFileHook(lFile.RotateFileConfig{
		Filename:   os.Getenv("STS_LOGGING_TO_FILE_FILENAME"),
		MaxSize:    orElse(os.Getenv("STS_LOGGING_TO_FILE_FILE_SIZE_MB"), 5),
		MaxBackups: orElse(os.Getenv("STS_LOGGING_TO_FILE_MAX_BACKUPS_FILES_NUMBER"), 1),
		MaxAge:     orElse(os.Getenv("STS_LOGGING_TO_FILE_MAX_BACKUPS_FILE_AGE"), 0),
		Level:      log.DebugLevel,
		Formatter:  &log.TextFormatter{},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.AddHook(hook)
}

func orElse(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	atoi, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return atoi
}

func setupLogLevel() {
	logLevel := os.Getenv("STS_LOGGING_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}

type cronLogrusLogger struct {
}

func (c cronLogrusLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Infof(msg, keysAndValues)
}

func (c cronLogrusLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Error(err)
	log.Errorf(msg, keysAndValues)
}

var collector stats.Collector
var writer stats.Writer

func task() {
	log.Debugf("Task started")
	err := writer.PingConnect()
	if err != nil {
		log.WithField("stackTrace", fmt.Sprintf("%+v", err)).Error(err)
		return
	}

	if !writer.Prepared() {
		err := writer.Prepare()
		if err != nil {
			log.WithField("stackTrace", fmt.Sprintf("%+v", err)).Error(err)
			return
		}
	}

	err = sendCpuMetrics()
	if err != nil {
		log.Error(err)
	}

	err = sendProcessesMetrics()
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Task finished successful")
}

func sendCpuMetrics() error {
	records, err := collector.CollectCpuRecords()
	if err != nil {
		return err
	}

	err = writer.WriteCpuRecords(records)
	if err != nil {
		return err
	}
	return nil
}

func sendProcessesMetrics() error {
	records, err := collector.CollectProcessesRecords()
	if err != nil {
		return err
	}

	err = writer.WriteProcessesRecords(records)
	if err != nil {
		return err
	}
	return nil
}
