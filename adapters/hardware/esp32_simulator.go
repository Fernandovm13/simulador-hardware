package hardware

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"simulador-hard/domain"
	"simulador-hard/ports"
)

type ESP32HardwareSimulator struct {
	mesaID       int
	publisher    ports.DataPublisher
	stopChan     chan struct{}
	mu           sync.RWMutex
	lastGas      domain.GasReading
	lastParticle domain.ParticleReading
}

func NewESP32Simulator(mesaID int, publisher ports.DataPublisher) *ESP32HardwareSimulator {
	return &ESP32HardwareSimulator{
		mesaID:    mesaID,
		publisher: publisher,
		stopChan:  make(chan struct{}),
	}
}

func (s *ESP32HardwareSimulator) Start() {
	go s.simulateGasSensor()
	go s.simulateParticleSensor()
}

func (s *ESP32HardwareSimulator) simulateGasSensor() {
	ticker := time.NewTicker(1800 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			baseLPG := 150.0 + rand.Float64()*250
			baseCO := 100.0 + rand.Float64()*200
			baseSmoke := 120.0 + rand.Float64()*230

			if rand.Float64() > 0.85 {
				spikeType := rand.Intn(3)
				spike := rand.Float64() * 400

				switch spikeType {
				case 0:
					baseLPG += spike
				case 1:
					baseCO += spike
				case 2:
					baseSmoke += spike
				}
			}

			reading := domain.GasReading{
				ID:        uuid.New().String(),  // Generar UUID
				SensorID:  fmt.Sprintf("ESP32-MESA-%d-GAS", s.mesaID),
				SystemID:  s.mesaID,
				LPG:       baseLPG,
				CO:        baseCO,
				Smoke:     baseSmoke,
				Timestamp: time.Now(),
			}

			s.mu.Lock()
			s.lastGas = reading
			s.mu.Unlock()

			if s.publisher != nil && s.publisher.IsConnected() {
				topic := fmt.Sprintf("vigiltech/sensors/mesa%d/gas", s.mesaID)
				s.publisher.Publish(topic, reading)
			}
		}
	}
}

func (s *ESP32HardwareSimulator) simulateParticleSensor() {
	ticker := time.NewTicker(2200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			pm10 := 10.0 + rand.Float64()*50
			pm25 := pm10 + 5.0 + rand.Float64()*30
			pm100 := pm25 + 10.0 + rand.Float64()*35

			if rand.Float64() > 0.8 {
				contaminationFactor := 1.5 + rand.Float64()*1.5
				pm10 *= contaminationFactor
				pm25 *= contaminationFactor
				pm100 *= contaminationFactor
			}

			reading := domain.ParticleReading{
				ID:        uuid.New().String(),  // Generar UUID
				SensorID:  fmt.Sprintf("ESP32-MESA-%d-PM", s.mesaID),
				SystemID:  s.mesaID,
				PM10:      pm10,
				PM25:      pm25,
				PM100:     pm100,
				Timestamp: time.Now(),
			}

			s.mu.Lock()
			s.lastParticle = reading
			s.mu.Unlock()

			if s.publisher != nil && s.publisher.IsConnected() {
				topic := fmt.Sprintf("vigiltech/sensors/mesa%d/particles", s.mesaID)
				s.publisher.Publish(topic, reading)
			}
		}
	}
}

func (s *ESP32HardwareSimulator) Stop() {
	close(s.stopChan)
}

func (s *ESP32HardwareSimulator) GetState() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &domain.ESP32State{
		MesaID:       s.mesaID,
		LastGas:      s.lastGas,
		LastParticle: s.lastParticle,
	}
}

func (s *ESP32HardwareSimulator) GetMesaID() int {
	return s.mesaID
}

func (s *ESP32HardwareSimulator) GetGasReading() domain.GasReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastGas
}

func (s *ESP32HardwareSimulator) GetParticleReading() domain.ParticleReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastParticle
}
