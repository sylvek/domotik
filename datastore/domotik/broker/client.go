package broker

type Log struct {
	Topic string
	Name  string
	Unit  string
	Value float32
}

type BrokerClient interface {
	ConnectAndListen(logs chan Log)
	Disconnect()
}
