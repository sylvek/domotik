package broker

type Log struct {
	Topic string
	Name  string
	Unit  string
	Value float64
}

type BrokerClient interface {
	ConnectAndListen()
	Logs() chan Log
	Disconnect()
	Publish(log Log)
}
