package stats

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

type gopsutilsCollector struct {
}

func (g *gopsutilsCollector) CollectProcessesRecords() ([]ProcessRecord, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	return convertProcessesStatsToRecords(processes)
}

func convertProcessesStatsToRecords(processes []*process.Process) ([]ProcessRecord, error) {
	result := make([]ProcessRecord, len(processes))
	for i, processStat := range processes {
		name, err := processStat.Name()
		if err != nil {
			log.Error(err)
		}
		username, err := processStat.Username()
		if err != nil {
			log.Error(err)
		}
		cmdline, err := processStat.Cmdline()
		if err != nil {
			log.Debug(err)
		}
		cpuPercent, err := processStat.CPUPercent()
		if err != nil {
			log.Error(err)
		}
		memoryPercent, err := processStat.MemoryPercent()
		if err != nil {
			log.Error(err)
		}
		cwd, err := processStat.Cwd()
		if err != nil {
			log.Debug(err)
		}
		exe, err := processStat.Exe()
		if err != nil {
			log.Debug(errors.Wrap(err, "Collecting Processes Exe stat error"))
		}
		parentPid, err := processStat.Ppid()
		if err != nil {
			log.Error(err)
		}
		status, err := processStat.Status()
		if err != nil {
			log.Error(err)
		}
		result[i] = convertProcessStatToRecord(
			processStat.Pid,
			parentPid,
			name,
			username,
			convertToProcessStatusEnum(status),
			cpuPercent,
			memoryPercent,
			cmdline,
			cwd,
			exe,
		)
	}
	return result, nil
}

func convertToProcessStatusEnum(status string) ProcessStatusRecord {
	switch status {
	case "D":
		return UninterruptibleSleep
	case "R":
		return Running
	case "S":
		return Sleep
	case "T":
		return Stop
	case "t":
		return StoppedByDebugger
	case "I":
		return Idle
	case "Z":
		return Zombie
	case "W":
		return Wait
	case "L":
		return Lock
	default:
		log.Error("Undefined process status: " + status)
		return Undefined
	}
}

func convertProcessStatToRecord(
	pid int32,
	parentPid int32,
	name string,
	username string,
	status ProcessStatusRecord,
	cpuPercent float64,
	memoryPercent float32,
	cmdline string,
	cwd string,
	exe string,
) ProcessRecord {
	return ProcessRecord{
		Pid:                     pid,
		ParentPid:               parentPid,
		Status:                  status,
		Name:                    name,
		Username:                username,
		CpuPercent:              cpuPercent,
		MemoryPercent:           memoryPercent,
		Cmdline:                 cmdline,
		CurrentWorkingDirectory: cwd,
		ExecutablePath:          exe,
	}
}

func NewGopsutilsCollector() *gopsutilsCollector {
	return &gopsutilsCollector{}
}

func (g *gopsutilsCollector) CollectCpuRecords() (*[]HostCpuRecord, error) {
	times, err := cpu.Times(true)
	if err != nil {
		return nil, err
	}
	percent, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}
	return convertCpuStatsToRecords(times, percent), nil
}

func convertCpuStatsToRecords(stats []cpu.TimesStat, percent []float64) *[]HostCpuRecord {
	result := make([]HostCpuRecord, len(stats))
	for i, stat := range stats {
		result[i] = convertCpuStatToRecord(&stat, percent[i])
	}
	return &result
}

func convertCpuStatToRecord(stat *cpu.TimesStat, percent float64) HostCpuRecord {
	record := HostCpuRecord{
		CPU:       stat.CPU,
		Percent:   percent,
		User:      stat.User,
		System:    stat.System,
		Idle:      stat.Idle,
		Nice:      stat.Nice,
		Iowait:    stat.Iowait,
		Irq:       stat.Irq,
		Softirq:   stat.Softirq,
		Steal:     stat.Steal,
		Guest:     stat.Guest,
		GuestNice: stat.GuestNice,
	}
	return record
}
