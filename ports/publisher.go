package ports

//define la interfaz para publicar datos de sensores
type DataPublisher interface {
	Publish(topic string, payload interface{}) error
	IsConnected() bool
	Connect() error
	Disconnect()
}