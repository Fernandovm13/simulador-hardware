package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTPublisher implementa el adaptador MQTT
type MQTTPublisher struct {
	client    mqtt.Client
	broker    string
	clientID  string
	connected bool
}

// NewMQTTPublisher crea un nuevo publicador MQTT
func NewMQTTPublisher(broker, clientID string) *MQTTPublisher {
	return &MQTTPublisher{
		broker:   broker,
		clientID: clientID,
	}
}

// Connect establece conexión con el broker MQTT
func (p *MQTTPublisher) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(p.broker)
	opts.SetClientID(p.clientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Conexión MQTT perdida: %v", err)
		p.connected = false
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Conectado al broker MQTT")
		p.connected = true
	})

	p.client = mqtt.NewClient(opts)

	token := p.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	p.connected = true
	return nil
}

// Publish publica un mensaje en un topic
func (p *MQTTPublisher) Publish(topic string, payload interface{}) error {
	if !p.IsConnected() {
		return nil // Silenciar si no está conectado
	}

	// Serializar payload a JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("MQTT marshal error topic=%s: %v", topic, err)
		return err
	}

	// Publicar con QoS 1
	token := p.client.Publish(topic, 1, false, data)
	token.Wait()
	if token.Error() != nil {
		log.Printf("MQTT publish error topic=%s: %v", topic, token.Error())
		return token.Error()
	}

	log.Printf("MQTT published topic=%s len=%d", topic, len(data))
	return nil
}

// IsConnected verifica si está conectado
func (p *MQTTPublisher) IsConnected() bool {
	return p.connected && p.client != nil && p.client.IsConnected()
}

// Disconnect cierra la conexión MQTT
func (p *MQTTPublisher) Disconnect() {
	if p.client != nil && p.client.IsConnected() {
		p.client.Disconnect(250)
		p.connected = false
		log.Println("Desconectado de MQTT")
	}
}
