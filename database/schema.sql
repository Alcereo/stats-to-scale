CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

--- Cpu metrics table

CREATE TABLE cpu_metrics
(
    time            TIMESTAMPTZ      NOT NULL,
    cpu             text             not null,
    percent         double precision not null,
    time_user       double precision,
    time_system     double precision,
    time_idle       double precision,
    time_nice       double precision,
    time_iowait     double precision,
    time_irq        double precision,
    time_softirq    double precision,
    time_steal      double precision,
    time_guest      double precision,
    time_guest_nice double precision
);

COMMENT ON COLUMN cpu_metrics.cpu IS 'Cpu number name';

SELECT create_hypertable('cpu_metrics', 'time');

--- Processes metrics table

-- man ps - PROCESS STATE CODES
CREATE TYPE process_status as ENUM (
    'Undefined'
        'Uninterruptible sleep'
        'Running',
    'Sleep',
    'Stop',
    'Stopped by debugger',
    'Idle',
    'Zombie',
    'Wait',
    'Lock'
    );

CREATE TABLE processes_metrics
(
    time                      TIMESTAMPTZ      NOT NULL,
    pid                       int2             not null,
    parent_pid                int2,
    name                      text,
    status                    process_status,
    username                  text,
    cpu_percent               double precision not null,
    memory_percent            double precision not null,
    cmd_line                  text,
    current_working_directory text,
    executable_path           text
);

SELECT create_hypertable('processes_metrics', 'time');
