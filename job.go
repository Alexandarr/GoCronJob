package job

// Jober represent a job
type Jober interface {
	Name() string
	Start()
	CanStart() bool
	Stop()
	CanStop() bool
	Status() interface{}
}
