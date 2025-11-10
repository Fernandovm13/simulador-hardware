package hardware

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"simulador-hard/domain"
	"simulador-hard/ports"
)

//implementa el simulador de ESP32
type ESP32HardwareSimulator struct {
	mesaID       int
	publisher    ports.DataPublisher
	stopChan     chan struct{}
	mu           sync.RWMutex
	lastGas      domain.GasReading
	lastParticle domain.ParticleReading
}

//crea un nuevo simulador ESP32
func NewESP32Simulator(mesaID int, publisher ports.DataPublisher) *ESP32HardwareSimulator {
	return &ESP32HardwareSimulator{
		mesaID:    mesaID,
		publisher: publisher,
		stopChan:  make(chan struct{}),
	}
}

//inicia la simulación del ESP32
func (s *ESP32HardwareSimulator) Start() {
	// Goroutine 1: Sensor de Gas
	go s.simulateGasSensor()

	// Goroutine 2: Sensor de Partículas
	go s.simulateParticleSensor()
}

//simula el sensor MQ-135
func (s *ESP32HardwareSimulator) simulateGasSensor() {
	ticker := time.NewTicker(1800 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:		
			baseLevel := 200.0 + rand.Float64()*300
			spike := 0.0
			if rand.Float64() > 0.85 {
				spike = rand.Float64() * 400
			}

			level := baseLevel + spike
			alert := level > 700

			reading := domain.GasReading{
				SensorID:  fmt.Sprintf("ESP32-MESA-%d-GAS", s.mesaID),
				MesaID:    s.mesaID,
				Level:     level,
				Alert:     alert,
				Timestamp: time.Now(),
			}

			// Actualizar estado interno
			s.mu.Lock()
			s.lastGas = reading
			s.mu.Unlock()

			// Publicar a MQTT
			if s.publisher != nil && s.publisher.IsConnected() {
				topic := fmt.Sprintf("vigiltech/sensors/mesa%d/gas", s.mesaID)
				s.publisher.Publish(topic, reading)
			}
		}
	}
}

//simula el sensor PMS5003
func (s *ESP32HardwareSimulator) simulateParticleSensor() {
	ticker := time.NewTicker(2200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			// Generar lectura de partículas
			pm25 := 15.0 + rand.Float64()*60
			if rand.Float64() > 0.8 {
				pm25 += rand.Float64() * 30
			}

			pm10 := pm25 + rand.Float64()*25 + 10
			alert := pm25 > 75

			reading := domain.ParticleReading{
				SensorID:  fmt.Sprintf("ESP32-MESA-%d-PM", s.mesaID),
				MesaID:    s.mesaID,
				PM25:      pm25,
				PM10:      pm10,
				Alert:     alert,
				Timestamp: time.Now(),
			}

			// Actualizar estado interno
			s.mu.Lock()
			s.lastParticle = reading
			s.mu.Unlock()

			// Publicar a MQTT
			if s.publisher != nil && s.publisher.IsConnected() {
				topic := fmt.Sprintf("vigiltech/sensors/mesa%d/particles", s.mesaID)
				s.publisher.Publish(topic, reading)
			}
		}
	}
}

//detiene la simulación
func (s *ESP32HardwareSimulator) Stop() {
	close(s.stopChan)
}

//retorna el estado actual
func (s *ESP32HardwareSimulator) GetState() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &domain.ESP32State{
		MesaID:       s.mesaID,
		LastGas:      s.lastGas,
		LastParticle: s.lastParticle,
	}
}

//retorna el ID de la mesa
func (s *ESP32HardwareSimulator) GetMesaID() int {
	return s.mesaID
}

//retorna la última lectura de gas
func (s *ESP32HardwareSimulator) GetGasReading() domain.GasReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastGas
}

//retorna la última lectura de partículas
func (s *ESP32HardwareSimulator) GetParticleReading() domain.ParticleReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastParticle
}