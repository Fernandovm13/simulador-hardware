package application

import (
	"log"
	"simulador-hard/ports"
)

// SimulatorService coordina los simuladores de hardware
type SimulatorService struct {
	esp32Simulators []ports.ESP32Simulator
	usbSimulator    ports.USBSimulator
	publisher       ports.DataPublisher
}

// NewSimulatorService crea un nuevo servicio de simulación
func NewSimulatorService(
	esp32s []ports.ESP32Simulator,
	usb ports.USBSimulator,
	publisher ports.DataPublisher,
) *SimulatorService {
	return &SimulatorService{
		esp32Simulators: esp32s,
		usbSimulator:    usb,
		publisher:       publisher,
	}
}

// StartAll inicia todos los simuladores
func (s *SimulatorService) StartAll() {
	log.Println("Iniciando simuladores...")

	// Iniciar ESP32s
	for _, sim := range s.esp32Simulators {
		sim.Start()
	}
	log.Printf("%d ESP32 simulados iniciados", len(s.esp32Simulators))

	// Iniciar USB
	s.usbSimulator.Start()
	log.Println("Sensores USB Direct iniciados")
}

// StopAll detiene todos los simuladores
func (s *SimulatorService) StopAll() {
	log.Println("Deteniendo simuladores...")

	// Detener ESP32s
	for _, sim := range s.esp32Simulators {
		sim.Stop()
	}

	// Detener USB
	s.usbSimulator.Stop()

	// Desconectar MQTT
	if s.publisher != nil {
		s.publisher.Disconnect()
	}

	log.Println("Simuladores detenidos correctamente")
}

// GetESP32Simulators retorna los simuladores ESP32
func (s *SimulatorService) GetESP32Simulators() []ports.ESP32Simulator {
	return s.esp32Simulators
}

// GetUSBSimulator retorna el simulador USB
func (s *SimulatorService) GetUSBSimulator() ports.USBSimulator {
	return s.usbSimulator
}

// IsMQTTConnected verifica si MQTT está conectado
func (s *SimulatorService) IsMQTTConnected() bool {
	if s.publisher == nil {
		return false
	}
	return s.publisher.IsConnected()
}