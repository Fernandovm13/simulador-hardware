package domain

import "time"

// GasReading representa una lectura del sensor de gas
type GasReading struct {
	SensorID  string    `json:"sensor_id"`
	MesaID    int       `json:"mesa_id"`
	Level     float64   `json:"level"`
	Alert     bool      `json:"alert"`
	Timestamp time.Time `json:"timestamp"`
}

// ParticleReading representa una lectura del sensor de partículas
type ParticleReading struct {
	SensorID  string    `json:"sensor_id"`
	MesaID    int       `json:"mesa_id"`
	PM25      float64   `json:"pm25"`
	PM10      float64   `json:"pm10"`
	Alert     bool      `json:"alert"`
	Timestamp time.Time `json:"timestamp"`
}

// MotionReading representa una lectura del sensor PIR
type MotionReading struct {
	SensorID  string    `json:"sensor_id"`
	Location  string    `json:"location"`
	Detected  bool      `json:"detected"`
	Timestamp time.Time `json:"timestamp"`
}

// CameraReading representa una lectura de la webcam
type CameraReading struct {
	SensorID      string    `json:"sensor_id"`
	Location      string    `json:"location"`
	HumanDetected bool      `json:"human_detected"`
	PhotoTaken    bool      `json:"photo_taken"`
	Timestamp     time.Time `json:"timestamp"`
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
	LastMotion MotionReading
	LastCamera CameraReading
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