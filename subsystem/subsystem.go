package subsystem

type ResourceConfig struct {
	MemoryLimit string
	CPUShare    string
	CPUSet      string
}

type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

var (
	SubsystemIns = []Subsystem{
		//&CPUSetSubSystem{},
		&MemorySubSystem{},
		//&CPUSubSystem{},
	}
)
