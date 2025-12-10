package hardware

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"simulador-hard/domain"
	"simulador-hard/ports"
)

type USBHardwareSimulator struct {
	publisher  ports.DataPublisher
	stopChan   chan struct{}
	motionChan chan string
	mu         sync.RWMutex
	lastMotion       domain.MotionReading
	lastCamera       domain.CameraReading
	lastCameraStream domain.CameraStreamReading
}

func NewUSBSimulator(publisher ports.DataPublisher) *USBHardwareSimulator {
	return &USBHardwareSimulator{
		publisher:  publisher,
		stopChan:   make(chan struct{}),
		motionChan: make(chan string, 10),
	}
}

func (s *USBHardwareSimulator) Start() {
	go s.simulatePIRSensor()        // Goroutine 1: PIR cada 2.5s
	go s.simulateWebcamCapture()    // Goroutine 2: Captura solo con movimiento
	go s.simulateCameraStream()     // Goroutine 3: Stream cada 1s
}

// PIR: SIEMPRE publica (detectado o no)
func (s *USBHardwareSimulator) simulatePIRSensor() {
	ticker := time.NewTicker(2500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			detected := rand.Float64() > 0.65

			intensity := 0.0
			if detected {
				intensity = 40.0 + rand.Float64()*60.0
			} else {
				intensity = rand.Float64() * 20.0
			}

			motionID := uuid.New().String()

			reading := domain.MotionReading{
				ID:             motionID,
				SensorID:       "USB-PIR-RPi",
				SystemID:       0,
				MotionDetected: detected,
				Intensity:      intensity,
				Timestamp:      time.Now(),
			}

			s.mu.Lock()
			s.lastMotion = reading
			s.mu.Unlock()

			// SIEMPRE publicar (detectado o no)
			if s.publisher != nil && s.publisher.IsConnected() {
				if err := s.publisher.Publish("vigiltech/sensors/usb/motion", reading); err != nil {
					log.Printf("ERROR publishing motion: %v", err)
				}
			}

			// Solo enviar por canal si hay movimiento
			if detected {
				select {
				case s.motionChan <- motionID:
					log.Printf("🚨 MOVIMIENTO DETECTADO - Motion ID: %s - Intensidad: %.1f%%", motionID, intensity)
				default:
				}
			} else {
				log.Printf("⚪ Sin movimiento - Intensidad: %.1f%%", intensity)
			}
		}
	}
}

// CAMERA CAPTURE: Solo cuando hay movimiento (camera_capture con motion_id)
func (s *USBHardwareSimulator) simulateWebcamCapture() {
	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()

	var currentMotionID string

	for {
		select {
		case <-s.stopChan:
			return

		case motionID := <-s.motionChan:
			currentMotionID = motionID
			log.Printf("📸 Cámara lista para capturar - Motion ID recibido: %s", motionID)

		case <-ticker.C:
			if currentMotionID != "" {
				photoURL := fmt.Sprintf("https://picsum.photos/seed/%d/640/480", time.Now().UnixNano())
				latency := 10 + rand.Intn(40)

				reading := domain.CameraReading{
					ID:        uuid.New().String(),
					SensorID:  "USB-WEBCAM-RPi",
					SystemID:  0,
					ImagePath: photoURL,
					MotionID:  currentMotionID,
					LatencyMs: latency,
					Timestamp: time.Now(),
				}

				s.mu.Lock()
				s.lastCamera = reading
				s.mu.Unlock()

				if s.publisher != nil && s.publisher.IsConnected() {
					log.Printf("📸 FOTO CAPTURADA - Motion ID: %s - URL: %s - Latencia: %dms",
						currentMotionID, photoURL, latency)

					if err := s.publisher.Publish("vigiltech/sensors/usb/camera", reading); err != nil {
						log.Printf("ERROR publishing camera: %v", err)
					}
				}

				currentMotionID = ""
			}
		}
	}
}

// CAMERA STREAM: SIEMPRE envía imágenes cada 1 segundo (camera_stream sin motion_id)
func (s *USBHardwareSimulator) simulateCameraStream() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			photoURL := fmt.Sprintf("https://picsum.photos/seed/%d/640/480", time.Now().UnixNano())
			latency := 5 + rand.Intn(15)

			reading := domain.CameraStreamReading{
				ID:        uuid.New().String(),
				SensorID:  "USB-WEBCAM-RPi-STREAM",
				SystemID:  0,
				ImagePath: photoURL,
				LatencyMs: latency,
				Timestamp: time.Now(),
			}

			s.mu.Lock()
			s.lastCameraStream = reading
			s.mu.Unlock()

			if s.publisher != nil && s.publisher.IsConnected() {
				if err := s.publisher.Publish("vigiltech/sensors/usb/camera_stream", reading); err != nil {
					log.Printf("ERROR publishing camera stream: %v", err)
				}
			}
		}
	}
}

func (s *USBHardwareSimulator) Stop() {
	close(s.stopChan)
}

func (s *USBHardwareSimulator) GetState() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &domain.USBState{
		LastMotion:       s.lastMotion,
		LastCamera:       s.lastCamera,
		LastCameraStream: s.lastCameraStream,
	}
}

func (s *USBHardwareSimulator) GetMotionReading() domain.MotionReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastMotion
}

func (s *USBHardwareSimulator) GetCameraReading() domain.CameraReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastCamera
}

func (s *USBHardwareSimulator) GetCameraStreamReading() domain.CameraStreamReading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastCameraStream
}
