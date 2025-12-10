package domain

import "time"

// GasReading representa una lectura del sensor de gas MQ-135
// Mapea exactamente a: gas_sensor (id, timestamp, lpg, co, smoke, system_id)
type GasReading struct {
	ID        string    `json:"id"`
	SensorID  string    `json:"sensor_id"`
	SystemID  int       `json:"system_id"`
	LPG       float64   `json:"lpg"`
	CO        float64   `json:"co"`
	Smoke     float64   `json:"smoke"`
	Timestamp time.Time `json:"timestamp"`
}

// ParticleReading representa una lectura del sensor de partículas PMS5003
// Mapea exactamente a: particle_sensor (id, timestamp, pm1_0, pm2_5, pm10, system_id)
type ParticleReading struct {
	ID        string    `json:"id"`
	SensorID  string    `json:"sensor_id"`
	SystemID  int       `json:"system_id"`
	PM10      float64   `json:"pm1_0"`
	PM25      float64   `json:"pm2_5"`
	PM100     float64   `json:"pm10"`
	Timestamp time.Time `json:"timestamp"`
}

// MotionReading representa una lectura del sensor PIR HC-SR501
// Mapea exactamente a: motion_sensors (id, timestamp, motion_detected, intensity, system_id)
type MotionReading struct {
	ID             string    `json:"id"`
	SensorID       string    `json:"sensor_id"`
	SystemID       int       `json:"system_id"`
	MotionDetected bool      `json:"motion_detected"`
	Intensity      float64   `json:"intensity"`
	Timestamp      time.Time `json:"timestamp"`
}

// CameraReading representa una captura de cámara cuando hay movimiento
// Mapea exactamente a: camera_capture (id, timestamp, image_path, motion_id, latency_ms, system_id)
type CameraReading struct {
	ID        string    `json:"id"`
	SensorID  string    `json:"sensor_id"`
	SystemID  int       `json:"system_id"`
	ImagePath string    `json:"image_path"`
	MotionID  string    `json:"motion_id"`
	LatencyMs int       `json:"latency_ms"`
	Timestamp time.Time `json:"timestamp"`
}

// CameraStreamReading representa el stream continuo de cámara (cada segundo)
// Mapea exactamente a: camera_stream (id, timestamp, image_path, system_id, latency_ms)
type CameraStreamReading struct {
	ID        string    `json:"id"`
	SensorID  string    `json:"sensor_id"`
	SystemID  int       `json:"system_id"`
	ImagePath string    `json:"image_path"`
	LatencyMs int       `json:"latency_ms"`
	Timestamp time.Time `json:"timestamp"`
}

// SystemState representa el estado del sistema completo
type SystemState struct {
	ESP32States   map[int]*ESP32State
	USBState      *USBState
	Time          float64
	MQTTConnected bool
}

// ESP32State representa el estado de un ESP32
type ESP32State struct {
	MesaID       int
	LastGas      GasReading
	LastParticle ParticleReading
}

// USBState representa el estado de los sensores USB
type USBState struct {
	LastMotion       MotionReading
	LastCamera       CameraReading
	LastCameraStream CameraStreamReading
}

// NewSystemState crea un nuevo estado del sistema
func NewSystemState(numMesas int) *SystemState {
	state := &SystemState{
		ESP32States: make(map[int]*ESP32State),
		USBState:    &USBState{},
	}
	for i := 1; i <= numMesas; i++ {
		state.ESP32States[i] = &ESP32State{
			MesaID: i,
		}
	}
	return state
}
