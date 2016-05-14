package veino

type Scheduler interface {
	Start()
	Add(string, string, func()) error
	Remove(string)
	Stop()
}
