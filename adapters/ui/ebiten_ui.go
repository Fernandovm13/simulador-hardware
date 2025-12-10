package ui

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"simulador-hard/domain"
	"simulador-hard/ports"
)

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 700
)

type EbitenUI struct {
	esp32Simulators []ports.ESP32Simulator
	usbSimulator    ports.USBSimulator
	mqttConnected   bool
	time            float64
	images          *imageCache
}

func NewEbitenUI(esp32s []ports.ESP32Simulator, usb ports.USBSimulator, mqttConnected bool) *EbitenUI {
	return &EbitenUI{
		esp32Simulators: esp32s,
		usbSimulator:    usb,
		mqttConnected:   mqttConnected,
		images:          newImageCache(),
	}
}

func (ui *EbitenUI) Update() error {
	ui.time += 1.0 / 60.0
	return nil
}

func (ui *EbitenUI) Draw(screen *ebiten.Image) {
	for y := 0; y < SCREEN_HEIGHT; y++ {
		intensity := uint8(25 + float32(y)/float32(SCREEN_HEIGHT)*30)
		vector.DrawFilledRect(screen, 0, float32(y), SCREEN_WIDTH, 1,
			color.RGBA{intensity, intensity + 8, intensity + 18, 255}, false)
	}

	ui.drawHeader(screen)
	ui.drawGrid(screen)
	ui.drawRaspberryPi(screen, 360, 60)

	ui.drawESP32Module(screen, 60, 240, 1)
	ui.drawESP32Module(screen, 360, 240, 2)
	ui.drawESP32Module(screen, 660, 240, 3)
	ui.drawESP32Module(screen, 960, 240, 4)

	ui.drawUSBModule(screen, 980, 60)
	ui.drawStatusPanel(screen, 1050, 420)
}

func (ui *EbitenUI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func (ui *EbitenUI) drawHeader(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, SCREEN_WIDTH, 50, color.RGBA{15, 20, 40, 240}, false)
	vector.DrawFilledRect(screen, 0, 48, SCREEN_WIDTH, 2, color.RGBA{0, 200, 255, 200}, false)

	title := "VIGILTECH - SIMULADOR DE HARDWARE | Arquitectura Hexagonal"
	if ui.mqttConnected {
		title += " | MQTT ACTIVO"
	}
	ebitenutil.DebugPrintAt(screen, title, 20, 10)

	subtitle := "4 ESP32 (Gas+PM) + USB (PIR+Webcam) | 10 Goroutines | Pipeline Pattern"
	ebitenutil.DebugPrintAt(screen, subtitle, 20, 28)

	timestamp := time.Now().Format("15:04:05")
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("⏱ %s", timestamp), 1140, 20)
}

func (ui *EbitenUI) drawGrid(screen *ebiten.Image) {
	gridColor := color.RGBA{35, 45, 65, 35}
	for x := 0; x < SCREEN_WIDTH; x += 50 {
		vector.StrokeLine(screen, float32(x), 50, float32(x), SCREEN_HEIGHT, 1, gridColor, false)
	}
	for y := 50; y < SCREEN_HEIGHT; y += 50 {
		vector.StrokeLine(screen, 0, float32(y), SCREEN_WIDTH, float32(y), 1, gridColor, false)
	}
}

func (ui *EbitenUI) drawRaspberryPi(screen *ebiten.Image, x, y float32) {
	vector.DrawFilledRect(screen, x, y, 300, 140, color.RGBA{80, 20, 80, 230}, false)
	vector.StrokeRect(screen, x, y, 300, 140, 3, color.RGBA{200, 100, 200, 255}, false)

	vector.DrawFilledRect(screen, x+5, y+5, 290, 25, color.RGBA{120, 40, 120, 255}, false)
	ebitenutil.DebugPrintAt(screen, "RASPBERRY PI 5 - CENTRAL", int(x+15), int(y+12))

	ebitenutil.DebugPrintAt(screen, "Broadcom BCM2712 | 8GB RAM", int(x+15), int(y+38))
	ebitenutil.DebugPrintAt(screen, "ARM Cortex-A76 @ 2.4GHz", int(x+15), int(y+52))

	vector.DrawFilledRect(screen, x+10, y+65, 280, 30, color.RGBA{20, 30, 50, 200}, false)
	vector.StrokeRect(screen, x+10, y+65, 280, 30, 2, color.RGBA{100, 200, 255, 255}, false)
	ebitenutil.DebugPrintAt(screen, "MQTT Broker: Mosquitto", int(x+20), int(y+73))

	statusText := "DESCONECTADO"
	statusColor := color.RGBA{255, 50, 50, 255}
	if ui.mqttConnected {
		statusText = "CONECTADO"
		statusColor = color.RGBA{0, 255, 100, 255}
	}
	ebitenutil.DebugPrintAt(screen, statusText, int(x+20), int(y+85))
	vector.DrawFilledCircle(screen, x+260, y+88, 6, statusColor, false)

	for i := 0; i < 4; i++ {
		ledOn := (int(ui.time*4)+i)%4 == 0 && ui.mqttConnected
		ledColor := color.RGBA{80, 80, 80, 255}
		if ledOn {
			ledColor = color.RGBA{0, 255, 0, 255}
			vector.DrawFilledCircle(screen, x+180+float32(i*25), y+108, 8,
				color.RGBA{0, 255, 0, 60}, false)
		}
		vector.DrawFilledCircle(screen, x+180+float32(i*25), y+108, 5, ledColor, false)
	}

	msgs := "0"
	if ui.mqttConnected {
		msgs = fmt.Sprintf("%.0f", 120+math.Sin(ui.time*2)*20)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MQTT Msg/s: %s", msgs), int(x+15), int(y+125))
}

func (ui *EbitenUI) drawESP32Module(screen *ebiten.Image, x, y float32, mesaID int) {
	vector.DrawFilledRect(screen, x, y, 280, 260, color.RGBA{139, 90, 43, 255}, false)
	vector.StrokeRect(screen, x, y, 280, 260, 3, color.RGBA{101, 67, 33, 255}, false)

	for i := 0; i < 13; i++ {
		alpha := uint8(40 + i*8)
		vector.StrokeLine(screen, x+5, y+10+float32(i*20), x+275, y+10+float32(i*20),
			1, color.RGBA{120, 80, 40, alpha}, false)
	}

	vector.DrawFilledRect(screen, x+10, y+10, 260, 25, color.RGBA{20, 25, 35, 200}, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mesa %d - ESP32 + Sensores", mesaID), int(x+20), int(y+17))

	vector.DrawFilledRect(screen, x+20, y+45, 240, 80, color.RGBA{30, 30, 50, 240}, false)
	vector.StrokeRect(screen, x+20, y+45, 240, 80, 2, color.RGBA{100, 150, 255, 255}, false)

	vector.DrawFilledRect(screen, x+35, y+55, 60, 60, color.RGBA{50, 50, 70, 255}, false)
	vector.StrokeRect(screen, x+35, y+55, 60, 60, 2, color.RGBA{150, 150, 200, 255}, false)
	ebitenutil.DebugPrintAt(screen, "ESP32", int(x+45), int(y+78))
	ebitenutil.DebugPrintAt(screen, "WROOM", int(x+40), int(y+92))

	for i := 0; i < 15; i++ {
		vector.DrawFilledRect(screen, x+35+float32(i*4), y+118, 2, 6,
			color.RGBA{200, 200, 0, 255}, false)
	}

	vector.DrawFilledRect(screen, x+105, y+55, 145, 35, color.RGBA{20, 30, 50, 200}, false)
	vector.StrokeRect(screen, x+105, y+55, 145, 35, 1, color.RGBA{100, 200, 255, 255}, false)
	ebitenutil.DebugPrintAt(screen, "WiFi -> MQTT", int(x+115), int(y+63))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mesa%d/*", mesaID), int(x+115), int(y+76))

	for i := 0; i < 3; i++ {
		ledOn := (int(ui.time*3)+i+mesaID)%3 == 0 && ui.mqttConnected
		ledColor := color.RGBA{60, 60, 60, 255}
		if ledOn {
			ledColor = color.RGBA{0, 255, 100, 255}
			vector.DrawFilledCircle(screen, x+115+float32(i*20), y+100, 6,
				color.RGBA{0, 255, 100, 80}, false)
		}
		vector.DrawFilledCircle(screen, x+115+float32(i*20), y+100, 4, ledColor, false)
	}

	var gasReading domain.GasReading
	var pmReading domain.ParticleReading

	if mesaID > 0 && mesaID <= len(ui.esp32Simulators) {
		gasReading = ui.esp32Simulators[mesaID-1].GetGasReading()
		pmReading = ui.esp32Simulators[mesaID-1].GetParticleReading()
	}

	ui.drawGasSensor(screen, x+30, y+140, gasReading)
	ui.drawParticleSensor(screen, x+160, y+140, pmReading)

	if ui.mqttConnected {
		for i := 0; i < 3; i++ {
			phase := ui.time*2 - float64(i)*0.3 - float64(mesaID)*0.2
			offset := math.Sin(phase) * 15
			alpha := uint8((math.Sin(phase)+1)/2*150 + 50)
			vector.DrawFilledCircle(screen, x+140+float32(offset), y-float32(i*15)-10, 3,
				color.RGBA{100, 200, 255, alpha}, false)
		}

		if int(ui.time*2)%4 == mesaID%4 {
			vector.DrawFilledRect(screen, x+125, y-60, 45, 15, color.RGBA{20, 30, 50, 200}, false)
			ebitenutil.DebugPrintAt(screen, "MQTT", int(x+130), int(y-57))
		}
	}
}

func (ui *EbitenUI) drawGasSensor(screen *ebiten.Image, x, y float32, reading domain.GasReading) {
	vector.DrawFilledCircle(screen, x+25, y+25, 22, color.RGBA{80, 80, 100, 255}, false)
	vector.StrokeCircle(screen, x+25, y+25, 22, 2, color.RGBA{150, 150, 180, 255}, false)

	hasAlert := reading.LPG > 700 || reading.CO > 700 || reading.Smoke > 700

	gasColor := color.RGBA{0, 255, 100, 255}
	if hasAlert {
		gasColor = color.RGBA{255, 50, 50, 255}
		pulse := math.Sin(ui.time*8)*5 + 23
		vector.StrokeCircle(screen, x+25, y+25, float32(pulse), 2,
			color.RGBA{255, 50, 50, 150}, false)
	}
	vector.DrawFilledCircle(screen, x+25, y+25, 12, gasColor, false)

	ebitenutil.DebugPrintAt(screen, "GAS", int(x+10), int(y+55))
	maxGas := reading.LPG
	if reading.CO > maxGas {
		maxGas = reading.CO
	}
	if reading.Smoke > maxGas {
		maxGas = reading.Smoke
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", maxGas), int(x+5), int(y+68))
}

func (ui *EbitenUI) drawParticleSensor(screen *ebiten.Image, x, y float32, reading domain.ParticleReading) {
	vector.DrawFilledRect(screen, x, y, 50, 50, color.RGBA{60, 60, 80, 255}, false)
	vector.StrokeRect(screen, x, y, 50, 50, 2, color.RGBA{120, 120, 150, 255}, false)

	for i := 0; i < 5; i++ {
		vector.StrokeLine(screen, x+5, y+10+float32(i*8), x+45, y+10+float32(i*8),
			1, color.RGBA{100, 100, 120, 255}, false)
	}

	hasAlert := reading.PM25 > 75

	particleColor := color.RGBA{255, 180, 100, 255}
	if hasAlert {
		particleColor = color.RGBA{255, 100, 0, 255}
	}

	for i := 0; i < 8; i++ {
		offset := math.Sin(ui.time*3+float64(i)*0.5) * 3
		size := 1.0 + math.Sin(ui.time*4+float64(i))*0.5
		vector.DrawFilledCircle(screen, x+10+float32(i*5), y+25+float32(offset),
			float32(size), particleColor, false)
	}

	ebitenutil.DebugPrintAt(screen, "PM", int(x+18), int(y+55))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1f", reading.PM25), int(x+10), int(y+68))
}

func (ui *EbitenUI) drawUSBModule(screen *ebiten.Image, x, y float32) {
	vector.DrawFilledRect(screen, x, y, 280, 180, color.RGBA{200, 200, 210, 255}, false)
	vector.StrokeRect(screen, x, y, 280, 180, 4, color.RGBA{150, 150, 160, 255}, false)

	vector.StrokeLine(screen, x+140, y, x+140, y+180, 3, color.RGBA{180, 180, 190, 255}, false)
	vector.StrokeLine(screen, x, y+90, x+280, y+90, 3, color.RGBA{180, 180, 190, 255}, false)

	vector.DrawFilledRect(screen, x+10, y+10, 260, 25, color.RGBA{40, 40, 60, 220}, false)
	ebitenutil.DebugPrintAt(screen, "ESQUINA - USB DIRECTO A RPI", int(x+30), int(y+17))

	vector.DrawFilledRect(screen, x+100, y+40, 80, 10, color.RGBA{80, 80, 90, 255}, false)
	vector.StrokeRect(screen, x+100, y+40, 80, 10, 2, color.RGBA{120, 120, 130, 255}, false)

	vector.DrawFilledRect(screen, x+80, y+55, 120, 110, color.RGBA{240, 240, 245, 255}, false)
	vector.StrokeRect(screen, x+80, y+55, 120, 110, 3, color.RGBA{100, 100, 120, 255}, false)

	vector.DrawFilledRect(screen, x+85, y+60, 110, 35, color.RGBA{80, 20, 80, 230}, false)
	vector.StrokeRect(screen, x+85, y+60, 110, 35, 2, color.RGBA{150, 100, 200, 255}, false)
	ebitenutil.DebugPrintAt(screen, "CONEXION USB", int(x+95), int(y+68))
	ebitenutil.DebugPrintAt(screen, "Directo -> RPi", int(x+95), int(y+81))

	for i := 0; i < 3; i++ {
		ledOn := (int(ui.time*4)+i)%3 == 0
		ledColor := color.RGBA{60, 60, 60, 255}
		if ledOn {
			ledColor = color.RGBA{150, 100, 200, 255}
			vector.DrawFilledCircle(screen, x+100+float32(i*20), y+105, 5,
				color.RGBA{150, 100, 200, 80}, false)
		}
		vector.DrawFilledCircle(screen, x+100+float32(i*20), y+105, 3, ledColor, false)
	}

	motionReading := ui.usbSimulator.GetMotionReading()
	cameraReading := ui.usbSimulator.GetCameraReading()

	ui.drawPIRSensor(screen, x+90, y+120, motionReading)
	ui.drawCamera(screen, x+155, y+120, cameraReading)

	vector.StrokeLine(screen, x+140, y, x+140, y-50, 4, color.RGBA{80, 80, 100, 220}, false)
	vector.StrokeLine(screen, x+138, y, x+138, y-50, 2, color.RGBA{200, 200, 220, 200}, false)

	for i := 0; i < 3; i++ {
		offset := math.Mod(ui.time*100+float64(i)*30, 50)
		alpha := uint8(255 - offset*4)
		vector.DrawFilledCircle(screen, x+140, y-float32(offset), 3,
			color.RGBA{150, 100, 200, alpha}, false)
	}

	vector.DrawFilledRect(screen, x+145, y-35, 60, 20, color.RGBA{80, 20, 80, 200}, false)
	ebitenutil.DebugPrintAt(screen, "USB 3.0", int(x+150), int(y-30))
	ebitenutil.DebugPrintAt(screen, "↑ RPi", int(x+155), int(y-18))
}

func (ui *EbitenUI) drawPIRSensor(screen *ebiten.Image, x, y float32, reading domain.MotionReading) {
	vector.DrawFilledCircle(screen, x+25, y+25, 25, color.RGBA{240, 240, 250, 255}, false)
	vector.StrokeCircle(screen, x+25, y+25, 25, 2, color.RGBA{180, 180, 200, 255}, false)

	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4
		dist := 12.0 + math.Sin(ui.time*2+float64(i))*2
		px := x + 25 + float32(math.Cos(angle)*dist)
		py := y + 25 + float32(math.Sin(angle)*dist)
		vector.DrawFilledCircle(screen, px, py, 3, color.RGBA{100, 100, 150, 255}, false)
	}

	if reading.MotionDetected {
		for i := 1; i <= 3; i++ {
			radius := 25 + float32(i*12) + float32(math.Sin(ui.time*6)*3)
			alpha := uint8(200 - i*60)
			vector.StrokeCircle(screen, x+25, y+25, radius, 2,
				color.RGBA{255, 50, 50, alpha}, false)
		}
		vector.DrawFilledCircle(screen, x+25, y+25, 8, color.RGBA{255, 50, 50, 255}, false)
	} else {
		vector.DrawFilledCircle(screen, x+25, y+25, 8, color.RGBA{100, 100, 150, 255}, false)
	}

	ebitenutil.DebugPrintAt(screen, "PIR", int(x+10), int(y+55))
	status := fmt.Sprintf("%.0f%%", reading.Intensity)
	if !reading.MotionDetected {
		status = "OFF"
	}
	ebitenutil.DebugPrintAt(screen, status, int(x+8), int(y+68))
}

func (ui *EbitenUI) drawCamera(screen *ebiten.Image, x, y float32, reading domain.CameraReading) {
	if reading.ImagePath != "" && ui.images != nil {
		ui.images.Load(reading.ImagePath)
	}

	imgX := x + 5
	imgY := y + 5
	imgW := float32(45)
	imgH := float32(40)

	if reading.ImagePath != "" && ui.images != nil {
		if img, ok := ui.images.Get(reading.ImagePath); ok && img != nil {
			w := float64(img.Bounds().Dx())
			h := float64(img.Bounds().Dy())
			if w > 0 && h > 0 {
				scaleX := imgW / float32(w)
				scaleY := imgH / float32(h)
				scale := float32(scaleX)
				if scaleY < scaleX {
					scale = float32(scaleY)
				}
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(float64(scale), float64(scale))
				tx := float64(imgX) + (float64(imgW)-float64(w)*float64(scale))/2.0
				ty := float64(imgY) + (float64(imgH)-float64(h)*float64(scale))/2.0
				op.GeoM.Translate(tx, ty)
				screen.DrawImage(img, op)
				vector.StrokeRect(screen, imgX, imgY, imgW, imgH, 1, color.RGBA{80, 80, 90, 180}, false)

				status := fmt.Sprintf("%dms", reading.LatencyMs)
				ebitenutil.DebugPrintAt(screen, "CAM", int(x+10), int(y+52))
				ebitenutil.DebugPrintAt(screen, status, int(x+8), int(y+65))
				return
			}
		}
	}

	vector.DrawFilledRect(screen, x+5, y+5, 45, 40, color.RGBA{30, 30, 35, 255}, false)
	vector.StrokeRect(screen, x+5, y+5, 45, 40, 2, color.RGBA{100, 100, 120, 255}, false)

	vector.DrawFilledCircle(screen, x+27, y+25, 15, color.RGBA{50, 50, 80, 255}, false)
	vector.StrokeCircle(screen, x+27, y+25, 15, 2, color.RGBA{150, 150, 200, 255}, false)
	vector.DrawFilledCircle(screen, x+27, y+25, 10, color.RGBA{20, 20, 40, 255}, false)
	vector.DrawFilledCircle(screen, x+27, y+25, 5, color.RGBA{100, 100, 200, 120}, false)

	ledColor := color.RGBA{80, 80, 80, 255}
	if reading.ImagePath != "" {
		ledColor = color.RGBA{255, 0, 0, 255}
		vector.DrawFilledCircle(screen, x+15, y+12, 5, color.RGBA{255, 0, 0, 100}, false)
	}
	vector.DrawFilledCircle(screen, x+15, y+12, 3, ledColor, false)

	ebitenutil.DebugPrintAt(screen, "CAM", int(x+10), int(y+52))
	status := "Wait"
	if reading.ImagePath != "" {
		status = "REC!"
	}
	ebitenutil.DebugPrintAt(screen, status, int(x+10), int(y+65))
}

func (ui *EbitenUI) drawStatusPanel(screen *ebiten.Image, x, y float32) {
	vector.DrawFilledRect(screen, x, y, 210, 220, color.RGBA{20, 25, 35, 230}, false)
	vector.StrokeRect(screen, x, y, 210, 220, 2, color.RGBA{100, 150, 200, 200}, false)

	ebitenutil.DebugPrintAt(screen, "ESTADO DEL SISTEMA", int(x+10), int(y+10))
	vector.StrokeLine(screen, x+10, y+25, x+200, y+25, 1,
		color.RGBA{100, 150, 200, 200}, false)

	yOffset := y + 35

	mqttStatus := "MQTT: OFF"
	mqttColor := color.RGBA{255, 100, 100, 255}
	if ui.mqttConnected {
		mqttStatus = "MQTT: ON"
		mqttColor = color.RGBA{0, 255, 100, 255}
	}
	vector.DrawFilledCircle(screen, x+15, yOffset, 3, mqttColor, false)
	ebitenutil.DebugPrintAt(screen, mqttStatus, int(x+25), int(yOffset-5))
	yOffset += 18

	ebitenutil.DebugPrintAt(screen, "Goroutines: 10", int(x+15), int(yOffset))
	yOffset += 18

	for i := 1; i <= 4; i++ {
		vector.DrawFilledCircle(screen, x+15, yOffset, 3, color.RGBA{0, 255, 100, 255}, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ESP32-%d: OK", i), int(x+25), int(yOffset-5))
		yOffset += 18
	}

	vector.DrawFilledCircle(screen, x+15, yOffset, 3, color.RGBA{0, 255, 100, 255}, false)
	ebitenutil.DebugPrintAt(screen, "USB Direct: OK", int(x+25), int(yOffset-5))
	yOffset += 18

	alertCnt := 0
	for _, sim := range ui.esp32Simulators {
		gas := sim.GetGasReading()
		pm := sim.GetParticleReading()
		if gas.LPG > 700 || gas.CO > 700 || gas.Smoke > 700 {
			alertCnt++
		}
		if pm.PM25 > 75 {
			alertCnt++
		}
	}

	alertColor := color.RGBA{0, 255, 100, 255}
	if alertCnt > 0 {
		alertColor = color.RGBA{255, 200, 0, 255}
	}
	vector.DrawFilledCircle(screen, x+15, yOffset, 3, alertColor, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Alertas: %d", alertCnt), int(x+25), int(yOffset-5))
}
