package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"simulador-hard/adapters/hardware"
	"simulador-hard/adapters/mqtt"
	"simulador-hard/adapters/ui"
	"simulador-hard/application"
	"simulador-hard/ports"
)

// Configuración
const (
	MQTT_BROKER  = "tcp://52.45.244.182:1883"
	MQTT_ENABLED = true 
	NUM_MESAS    = 4
)

func main() {
	rand.Seed(time.Now().UnixNano())

	printBanner()

	//Configurar MQTT Publisher 
	var publisher ports.DataPublisher
	mqttConnected := false

	if MQTT_ENABLED {
		mqttPub := mqtt.NewMQTTPublisher(MQTT_BROKER, "vigiltech-hardware-simulator")

		if err := mqttPub.Connect(); err != nil {
			log.Printf("No se pudo conectar a MQTT: %v", err)
			log.Println("Continuando solo con visualización...")
			publisher = nil
		} else {
			publisher = mqttPub
			mqttConnected = true
			log.Println("MQTT conectado - Publicando datos")
		}
	} else {
		log.Println("MQTT deshabilitado - Solo visualización")
	}

	//Crear simuladores ESP32 (Adaptadores Primarios)
	esp32Simulators := make([]ports.ESP32Simulator, NUM_MESAS)
	for i := 1; i <= NUM_MESAS; i++ {
		esp32Simulators[i-1] = hardware.NewESP32Simulator(i, publisher)
	}

	//Crear simulador USB Direct (Adaptador Primario)
	usbSimulator := hardware.NewUSBSimulator(publisher)


	//Crear servicio de simulación
	simulatorService := application.NewSimulatorService(
		esp32Simulators,
		usbSimulator,
		publisher,
	)

	//Iniciar todos los simuladores
	simulatorService.StartAll()

	log.Println("========================================")
	log.Println("Iniciando visualización gráfica...")
	log.Println("========================================")

	//Crear interfaz Ebiten
	game := ui.NewEbitenUI(esp32Simulators, usbSimulator, mqttConnected)

	//Configurar ventana
	ebiten.SetWindowSize(1280, 700)
	ebiten.SetWindowTitle("VigiTech - Simulador de Hardware | Arquitectura Hexagonal")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	//Cleanup al finalizar
	defer simulatorService.StopAll()

	// Ejecutar aplicación
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func printBanner() {
	log.Println("========================================")
	log.Println("  VIGILTECH - Simulador de Hardware")
	log.Println("  Arquitectura Hexagonal")
	log.Println("========================================")
	log.Println("")
}