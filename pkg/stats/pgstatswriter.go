package stats

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var insertCpuMetricsStmtQuery = `
INSERT INTO 
cpu_metrics(
  time,
  cpu,
  percent,
  time_user,
  time_system,
  time_idle,
  time_nice,
  time_iowait,
  time_irq,
  time_softirq,
  time_steal,
  time_guest,
  time_guest_nice
) VALUES (
  now(), -- Start time of the current transaction in Postgres
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12
)
`

var insertProcessesMetricsStmtQuery = `
INSERT INTO 
processes_metrics(
  time,
  pid,
  parent_pid,
  name,
  status,
  username,
  cpu_percent,
  memory_percent,
  cmd_line,
  current_working_directory,
  executable_path
) VALUES (
  now(), -- Start time of the current transaction in Postgres
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10
)
`

type pgWriter struct {
	db                         *sql.DB
	connectionString           string
	insertCpuMetricsStmt       *sql.Stmt
	insertProcessesMetricsStmt *sql.Stmt
	prepared                   bool
}

func NewPgWriter(
	connectionString string,
) (*pgWriter, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "Creating connecting error")
	}
	return &pgWriter{
		db:       db,
		prepared: false,
	}, nil
}

func (p *pgWriter) createPreparingStatements() error {
	stmt, err := p.db.Prepare(insertCpuMetricsStmtQuery)
	if err != nil {
		return err
	}
	p.insertCpuMetricsStmt = stmt

	stmt, err = p.db.Prepare(insertProcessesMetricsStmtQuery)
	if err != nil {
		return err
	}
	p.insertProcessesMetricsStmt = stmt

	return nil
}

func (p *pgWriter) WriteCpuRecords(records *[]HostCpuRecord) error {
	txn, err := p.db.Begin()
	if err != nil {
		return err
	}

	for _, record := range *records {
		_, err := txn.Stmt(p.insertCpuMetricsStmt).Exec(
			record.CPU,
			record.Percent,
			record.User,
			record.System,
			record.Idle,
			record.Nice,
			record.Iowait,
			record.Irq,
			record.Softirq,
			record.Steal,
			record.Guest,
			record.GuestNice,
		)
		if err != nil {
			err := txn.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *pgWriter) WriteProcessesRecords(records []ProcessRecord) error {
	txn, err := p.db.Begin()
	if err != nil {
		return err
	}

	for _, record := range records {
		_, err := txn.Stmt(p.insertProcessesMetricsStmt).Exec(
			record.Pid,
			record.ParentPid,
			record.Name,
			record.Status,
			record.Username,
			record.CpuPercent,
			record.MemoryPercent,
			record.Cmdline,
			record.CurrentWorkingDirectory,
			record.ExecutablePath,
		)
		if err != nil {
			err := txn.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *pgWriter) Prepared() bool {
	return p.prepared
}

func (p *pgWriter) Prepare() error {
	err := p.createPreparingStatements()
	if err != nil {
		return errors.Wrap(err, "Creating prepared statement error")
	}
	p.prepared = true
	return nil
}

func (p *pgWriter) PingConnect() error {
	if p.db == nil {
		return errors.New("connection hasn't open yet. use PrepareWorker first")
	}
	return errors.Wrap(p.db.Ping(), "Database ping error")
}
