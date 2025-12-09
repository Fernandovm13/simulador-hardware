package hardware

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"simulador-hard/domain"
	"simulador-hard/ports"
)

// USBHardwareSimulator implementa el simulador USB Direct
type USBHardwareSimulator struct {
	publisher  ports.DataPublisher
	stopChan   chan struct{}
	motionChan chan bool
	mu         sync.RWMutex
	lastMotion domain.MotionReading
	lastCamera domain.CameraReading
}

// NewUSBSimulator crea un nuevo simulador USB
func NewUSBSimulator(publisher ports.DataPublisher) *USBHardwareSimulator {
	return &USBHardwareSimulator{
		publisher:  publisher,
		stopChan:   make(chan struct{}),
		motionChan: make(chan bool, 10), // Buffer de 10
	}
}

// Start inicia la simulación USB
func (s *USBHardwareSimulator) Start() {
	// Goroutine 1: Sensor PIR
	go s.simulatePIRSensor()

	// Goroutine 2: Webcam (coordinada con PIR mediante canal)
	go s.simulateWebcam()
}

// simulatePIRSensor simula el sensor PIR HC-SR501
func (s *USBHardwareSimulator) simulatePIRSensor() {
	ticker := time.NewTicker(2500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			// Generar detección de movimiento
			detected := rand.Float64() > 0.65

			reading := domain.MotionReading{
				SensorID:  "USB-PIR-RPi",
				Location:  "usb_direct",
				Detected:  detected,
				Timestamp: time.Now(),
			}

			// Actualizar estado interno
			s.mu.Lock()
			s.lastMotion = reading
			s.mu.Unlock()

			// Enviar señal al canal (PATRÓN PIPELINE)
			select {
			case s.motionChan <- detected:
			default:
				// Canal lleno, descartar
			}

			// Publicar a MQTT
			if s.publisher != nil && s.publisher.IsConnected() {
				if err := s.publisher.Publish("vigiltech/sensors/usb/motion", reading); err != nil {
					log.Printf("ERROR publishing motion reading: %v", err)
				}
			}
		}
	}
}

// simulateWebcam simula la webcam USB
func (s *USBHardwareSimulator) simulateWebcam() {
	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()

	lastMotion := false

	for {
		select {
		case <-s.stopChan:
			return

		case motion := <-s.motionChan:
			// Recibir señal del PIR (PATRÓN PIPELINE)
			lastMotion = motion

		case <-ticker.C:
			var reading domain.CameraReading

			// Procesar solo si hay movimiento detectado
if lastMotion {
    humanDetected := rand.Float64() > 0.2
    var url string
    if humanDetected {
        // Genera URL pública de prueba (picsum.photos). Seed único para variar.
        url = fmt.Sprintf("https://picsum.photos/seed/%d/640/480", time.Now().UnixNano())
    }

    reading = domain.CameraReading{
        SensorID:      "USB-WEBCAM-RPi",
        Location:      "usb_direct",
        HumanDetected: humanDetected,
        PhotoTaken:    humanDetected,
        ImageURL:      url,
        Timestamp:     time.Now(),
    }
} else {
    // Sin movimiento → no debe haber foto ni persona detectada
    reading = domain.CameraReading{
        SensorID:      "USB-WEBCAM-RPi",
        Location:      "usb_direct",
        HumanDetected: false,
        PhotoTaken:    false,
        ImageURL:      "",
        Timestamp:     time.Now(),
    }
}

			// Actualizar estado interno
			s.mu.Lock()
			s.lastCamera = reading
			s.mu.Unlock()

			// Publicar a MQTT (con logging)
			if s.publisher != nil && s.publisher.IsConnected() {
				log.Printf("Publishing camera reading: sensor=%s photo_taken=%v human_detected=%v image_url=%q",
					reading.SensorID, reading.PhotoTaken, reading.HumanDetected, reading.ImageURL)

				if err := s.publisher.Publish("vigiltech/sensors/usb/camera", reading); err != nil {
					log.Printf("ERROR publishing camera reading: %v", err)
				}
			}
		}
	}
}

// Stop detiene la simulación
func (s *USBHardwareSimulator) Stop() {
	close(s.stopChan)
}

// GetState retorna el estado actual
func (s *USBHardwareSimulator) GetState() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &domain.USBState{
		LastMotion: s.lastMotion,
		LastCamera: s.lastCamera,
	}
}

// GetMotionReading retorna la última lectura de movimiento
func (s *USBHardwareSimulator) GetMotionReading() domain.MotionReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastMotion
}

// GetCameraReading retorna la última lectura de cámara
func (s *USBHardwareSimulator) GetCameraReading() domain.CameraReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastCamera
}
