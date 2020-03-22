package stats

type Collector interface {
	CollectCpuRecords() (*[]HostCpuRecord, error)
	CollectProcessesRecords() ([]ProcessRecord, error)
}

type Writer interface {
	WriteCpuRecords(*[]HostCpuRecord) error
	WriteProcessesRecords([]ProcessRecord) error
	PingConnect() error
	Prepared() bool
	Prepare() error
}

type HostCpuRecord struct {
	CPU       string  `json:"cpu"`
	Percent   float64 `json:"percent"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}

type IOCountersRecord struct {
	ReadCount  uint64 `json:"readCount"`
	WriteCount uint64 `json:"writeCount"`
	ReadBytes  uint64 `json:"readBytes"`
	WriteBytes uint64 `json:"writeBytes"`
}

type ProcessStatusRecord string

const (
	Undefined            ProcessStatusRecord = "Undefined"
	UninterruptibleSleep                     = "Uninterruptible sleep"
	Running                                  = "Running"
	Sleep                                    = "Sleep"
	Stop                                     = "Stop"
	StoppedByDebugger                        = "Stopped by debugger"
	Idle                                     = "Idle"
	Zombie                                   = "Zombie"
	Wait                                     = "Wait"
	Lock                                     = "Lock"
)

type ProcessRecord struct {
	Pid                     int32
	ParentPid               int32
	Name                    string
	Status                  ProcessStatusRecord
	Username                string
	CpuPercent              float64
	MemoryPercent           float32
	Cmdline                 string
	CurrentWorkingDirectory string
	ExecutablePath          string
}
