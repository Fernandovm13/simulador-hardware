package ports

import "simulador-hard/domain"

//define el contrato para simular sensores
type SensorSimulator interface {
	Start()
	Stop()
	GetState() interface{}
}

//define el contrato para simuladores ESP32
type ESP32Simulator interface {
	SensorSimulator
	GetMesaID() int
	GetGasReading() domain.GasReading
	GetParticleReading() domain.ParticleReading
}

//define el contrato para simulador USB
type USBSimulator interface {
	SensorSimulator
	GetMotionReading() domain.MotionReading
	GetCameraReading() domain.CameraReading
}