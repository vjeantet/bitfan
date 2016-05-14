package veino

type Processor interface {
	Configure(map[string]interface{}) error
	Start(IPacket) error
	Tick(IPacket) error
	Receive(IPacket) error
	Stop(IPacket) error
}

type ProcessorFactory func(Logger) Processor
