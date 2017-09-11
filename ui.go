package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gizak/termui"
	"github.com/luan/f1-telemetry/f1"
)

type SpeedUnit int

var TeamColors = map[byte]string{
	4:  "fg-mercedes",
	0:  "fg-redbull",
	1:  "fg-ferrari",
	6:  "fg-forceindia",
	7:  "fg-williams",
	2:  "fg-mclaren",
	8:  "fg-tororosso",
	11: "fg-haas",
	3:  "fg-renault",
	5:  "fg-sauber",
}

var TyreColors = map[byte]string{
	0: "fg-ultrasoft",
	1: "fg-supersoft",
	2: "fg-soft",
	3: "fg-medium",
	4: "fg-hard",
	5: "fg-inter",
	6: "fg-wet",
}

var SectorColors = map[byte]string{
	0: "fg-white",
	1: "fg-cyan",
	2: "fg-green",
}

const (
	MPH SpeedUnit = iota
	KPH
)

type UI struct {
	dataChan  <-chan f1.TelemetryData
	speedUnit atomic.Value

	components []termui.Bufferer

	logoPar     *termui.Par
	speedPar    *termui.Par
	throttle    *termui.Gauge
	brake       *termui.Gauge
	driverTable *termui.Table
	lapsTable   *termui.Table

	carPar *termui.Par

	playerLaps [][4]float32
}

func NewUI(dataChan <-chan f1.TelemetryData) *UI {
	err := termui.Init()
	if err != nil {
		log.Fatal(err)
	}

	ui := &UI{
		dataChan: dataChan,

		logoPar:     termui.NewPar(""),
		speedPar:    termui.NewPar(""),
		brake:       termui.NewGauge(),
		throttle:    termui.NewGauge(),
		driverTable: termui.NewTable(),
		lapsTable:   termui.NewTable(),

		carPar: termui.NewPar(""),

		playerLaps: [][4]float32{},
	}

	ui.initColors()

	ui.speedUnit.Store(KPH)
	ui.components = []termui.Bufferer{ui.logoPar, ui.speedPar, ui.brake, ui.throttle, ui.driverTable, ui.lapsTable, ui.carPar}

	ui.logoPar.Height = 7
	ui.logoPar.Width = 100
	ui.logoPar.X = 0
	ui.logoPar.Y = 1
	ui.logoPar.Border = false
	ui.logoPar.Text = `
                          _____   __                             _______
      _    _             / ___/.-' /                             _______
      \'../ |o_..__     / /__   / /                             _\=.o.=/_
    '.,(_)______(_).>  / ___/  / /                             |_|_____|_|
    ~~~~~~~~~~~~~~~~~~/_/~~~~~/_/~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	ui.speedPar.Height = 6
	ui.speedPar.Width = 60
	ui.speedPar.X = 0
	ui.speedPar.Y = 8
	ui.speedPar.Border = false

	ui.brake.Width = 30
	ui.brake.Height = 3
	ui.brake.X = 0
	ui.brake.Y = 14
	ui.brake.Border = true
	ui.brake.BarColor = termui.ColorRed
	ui.brake.BorderLabel = "Brakes"

	ui.throttle.Width = 30
	ui.throttle.Height = 3
	ui.throttle.X = 30
	ui.throttle.Y = 14
	ui.throttle.Border = true
	ui.throttle.BarColor = termui.ColorGreen
	ui.throttle.BorderLabel = "Throttle"

	ui.driverTable.Width = 60
	ui.driverTable.Height = 22
	ui.driverTable.X = 95
	ui.driverTable.Y = 8
	ui.driverTable.BorderFg = termui.ColorWhite
	ui.driverTable.Separator = false

	ui.carPar.Height = 47
	ui.carPar.Width = 33
	ui.carPar.X = 61
	ui.carPar.Y = 8
	ui.carPar.Border = false
	ui.carPar.Text = strings.Join(car, "\n")

	ui.lapsTable.Width = 60
	ui.lapsTable.Height = 22
	ui.lapsTable.X = 0
	ui.lapsTable.Y = 17
	ui.lapsTable.BorderLabel = "Laps"
	ui.lapsTable.BorderFg = termui.ColorWhite
	ui.lapsTable.Separator = false

	ui.setupEvents()

	return ui
}

func (ui *UI) initColors() {
	termui.SetOutputMode(termui.Output256)
	termui.AddColorMap("mercedes", termui.Attribute(1+49))
	termui.AddColorMap("ferrari", termui.Attribute(1+124))
	termui.AddColorMap("redbull", termui.Attribute(1+63))
	termui.AddColorMap("forceindia", termui.Attribute(1+171))
	termui.AddColorMap("williams", termui.Attribute(1+21))
	termui.AddColorMap("tororosso", termui.Attribute(1+27))
	termui.AddColorMap("haas", termui.Attribute(1+52))
	termui.AddColorMap("renault", termui.Attribute(1+226))
	termui.AddColorMap("mclaren", termui.Attribute(1+16))
	termui.AddColorMap("sauber", termui.Attribute(1+39))

	termui.AddColorMap("ultrasoft", termui.Attribute(1+93))
	termui.AddColorMap("supersoft", termui.Attribute(1+124))
	termui.AddColorMap("soft", termui.Attribute(1+226))
	termui.AddColorMap("medium", termui.Attribute(1+21))
	termui.AddColorMap("hard", termui.Attribute(1+16))
	termui.AddColorMap("inter", termui.Attribute(1+40))
	termui.AddColorMap("wet", termui.Attribute(1+39))

	termui.AddColorMap("pitting", termui.Attribute(1+247))
	termui.AddColorMap("inpits", termui.Attribute(1+239))
}

func (ui *UI) setupEvents() {
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/s", func(termui.Event) {
		if ui.speedUnit.Load().(SpeedUnit) == MPH {
			ui.speedUnit.Store(KPH)
		} else {
			ui.speedUnit.Store(MPH)
		}
	})
}

func (ui *UI) render() {
	termui.Render(ui.components...)
}

func (ui *UI) Start() {
	defer termui.Close()

	ui.processTelemetry(f1.TelemetryData{
		Speed:          40,
		Throttle:       1,
		Brake:          1,
		PlayerCarIndex: 0,
		Cars: [20]f1.CarData{
			{
				CurrentLapNum: 1,
				LastlapTime:   0,
				Sector:        0,
				Sector1Time:   0,
				Sector2Time:   0,
			},
		},
	})

	ui.processTelemetry(f1.TelemetryData{
		Speed:          40,
		Throttle:       1,
		Brake:          1,
		PlayerCarIndex: 0,
		Cars: [20]f1.CarData{
			{
				CurrentLapNum: 2,
				LastlapTime:   75.1234,
				Sector:        0,
				Sector1Time:   30,
				Sector2Time:   30,
			},
		},
	})

	ui.processTelemetry(f1.TelemetryData{
		Speed:          40,
		Throttle:       1,
		Brake:          1,
		PlayerCarIndex: 0,
		Cars: [20]f1.CarData{
			{
				CurrentLapNum: 3,
				LastlapTime:   72.4333,
				Sector:        0,
				Sector1Time:   28,
				Sector2Time:   33,
			},
		},
	})

	ui.processTelemetry(f1.TelemetryData{
		Speed:          40,
		Throttle:       1,
		Brake:          1,
		PlayerCarIndex: 0,
		Cars: [20]f1.CarData{
			{
				CurrentLapNum:  3,
				LastlapTime:    72.4333,
				Sector:         2,
				Sector1Time:    29,
				Sector2Time:    33,
				CurrentlapTime: 70,
			},
		},
	})

	ui.render()

	wg := sync.WaitGroup{}
	wg.Add(2)
	signal := make(chan struct{})

	go func() {
		defer wg.Done()
		termui.Loop()
		close(signal)
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case telemetry := <-ui.dataChan:
				ui.processTelemetry(telemetry)
				ui.render()
			case <-signal:
				return
			}
		}
	}()

	wg.Wait()
}

func (ui *UI) processTelemetry(telemetry f1.TelemetryData) {
	ui.speedPar.Text = ui.renderSpeed(telemetry.Speed)
	ui.brake.Percent = int(100 * telemetry.Brake)
	ui.throttle.Percent = int(100 * telemetry.Throttle)

	sortedCars := sortCars(telemetry.Cars)

	ui.processPlayerLap(telemetry)
	ui.renderPlayerLaps()

	ui.renderCar(telemetry)

	ui.renderCars(sortedCars, telemetry.TrackSize, byte(telemetry.SessionType))
}

func (ui *UI) renderCar(telemetry f1.TelemetryData) {
	carCopy := make([]string, len(car))
	copy(carCopy, car)

	ui.carPar.Text = strings.Join(carCopy, "\n")
}

func (ui *UI) renderPlayerLaps() {
	ui.lapsTable.Rows = make([][]string, 2+len(ui.playerLaps))

	ui.lapsTable.Rows[0] = []string{
		"#",
		"Sector 1",
		"Sector 2",
		"Sector 3",
		"Lap Time",
	}

	ui.lapsTable.Rows[1] = []string{
		"--",
		"---------",
		"---------",
		"---------",
		"--------------",
	}

	lowest := []float32{
		float32(math.MaxFloat32),
		float32(math.MaxFloat32),
		float32(math.MaxFloat32),
		float32(math.MaxFloat32),
	}

	for i, lap := range ui.playerLaps {
		for j := range lap {
			if i == len(ui.playerLaps)-1 &&
				(j == 3 || j == 2 || lap[j+1] == 0) {
				continue
			}
			if lap[j] > 0 && lap[j] < lowest[j] {
				lowest[j] = lap[j]
			}
		}
	}

	for i, lap := range ui.playerLaps {
		s := []string{fmt.Sprintf("%2d", i+1), "", "", "", ""}
		for j := range lap {
			if lap[j] > 0 {
				s[j+1] = floatToTime(lap[j])
			}
			if lap[j] == lowest[j] {
				s[j+1] = fmt.Sprintf("[%s](fg-magenta)", s[j+1])
			}
			if i == len(ui.playerLaps)-1 &&
				(j == 3 || j == 2 || lap[j+1] == 0) {
				s[j+1] = fmt.Sprintf("[%s](fg-green)", s[j+1])
			}
		}
		ui.lapsTable.Rows[i+2] = s
	}
}

func (ui *UI) processPlayerLap(telemetry f1.TelemetryData) {
	playerCar := telemetry.Cars[telemetry.PlayerCarIndex]
	for int(playerCar.CurrentLapNum) > len(ui.playerLaps) {
		ui.playerLaps = append(ui.playerLaps, [4]float32{0, 0, 0, 0})
	}
	if playerCar.CurrentLapNum >= 2 {
		ui.playerLaps[playerCar.CurrentLapNum-2][3] = playerCar.LastlapTime
	}
	ui.playerLaps[playerCar.CurrentLapNum-1][3] = playerCar.CurrentlapTime
	switch playerCar.Sector {
	case 0:
		if playerCar.CurrentLapNum >= 2 {
			ui.playerLaps[playerCar.CurrentLapNum-2][0] = playerCar.Sector1Time
			ui.playerLaps[playerCar.CurrentLapNum-2][1] = playerCar.Sector2Time
			ui.playerLaps[playerCar.CurrentLapNum-2][2] = playerCar.LastlapTime - playerCar.Sector1Time - playerCar.Sector2Time

			ui.playerLaps[playerCar.CurrentLapNum-1][0] = playerCar.CurrentlapTime
		}
	case 1:
		if playerCar.CurrentLapNum >= 1 {
			ui.playerLaps[playerCar.CurrentLapNum-1][0] = playerCar.Sector1Time
			ui.playerLaps[playerCar.CurrentLapNum-1][1] = playerCar.CurrentlapTime - playerCar.Sector1Time
		}
	case 2:
		if playerCar.CurrentLapNum >= 1 {
			ui.playerLaps[playerCar.CurrentLapNum-1][0] = playerCar.Sector1Time
			ui.playerLaps[playerCar.CurrentLapNum-1][1] = playerCar.Sector2Time
			ui.playerLaps[playerCar.CurrentLapNum-1][2] = playerCar.CurrentlapTime - playerCar.Sector1Time - playerCar.Sector2Time
		}
	}
}

func (ui *UI) renderCars(sortedCars []f1.CarData, trackSize float32, sessionType byte) {
	ui.driverTable.Rows = make([][]string, len(sortedCars))
	for i, car := range sortedCars {
		if car.CarPosition == 0 {
			continue
		}
		lapColor := "fg-white"
		if car.LastlapTime == car.BestlapTime {
			lapColor = "fg-magenta"
		}
		nameColor := "fg-white"
		switch car.InPits {
		case 1:
			nameColor = "fg-pitting"
		case 2:
			nameColor = "fg-inpits"
		}
		driverName := f1.Drivers[car.DriverID]
		if driverName == "" {
			driverName = "LUA"
		}
		lapToShow := car.BestlapTime
		if sessionType == 3 {
			lapToShow = car.LastlapTime
		}
		ui.driverTable.Rows[i] = []string{
			fmt.Sprintf("%2d [▶](%s) [%s](%s) (%d) [o](%s)  | [%s](%s) | [%s (%.1f%%)](%s)",
				car.CarPosition,
				TeamColors[car.TeamID],
				driverName,
				nameColor,
				car.CurrentLapNum,
				TyreColors[car.TyreCompound],
				floatToTime(lapToShow),
				lapColor,
				floatToTime(car.CurrentlapTime),
				100*(car.LapDistance/trackSize),
				SectorColors[car.Sector],
			),
		}
	}
}

func sortCars(cars [20]f1.CarData) []f1.CarData {
	sortedCars := make([]f1.CarData, 20)
	for _, car := range cars {
		if car.CarPosition == 0 {
			continue
		}
		sortedCars[car.CarPosition-1] = car
	}
	return sortedCars
}

func floatToTime(t float32) string {
	m := int(t) / 60
	t -= float32(m * 60)
	s := int(t)
	t -= float32(s)
	r := int(t * 10000)
	return fmt.Sprintf("%d:%02d.%04d", m, s, r)
}

func (ui *UI) floatToDistance(d float32) string {
	unit := "mi"
	conversion := float32(1609.34)
	if ui.speedUnit.Load().(SpeedUnit) == KPH {
		unit = "km"
		conversion = 1000
	}

	return fmt.Sprintf("%.1f%s", d/conversion, unit)
}

func (ui *UI) renderSpeed(mps float32) string {
	unit := "mph"
	conversion := float32(2.23694)
	if ui.speedUnit.Load().(SpeedUnit) == KPH {
		unit = "km/h"
		conversion = 3.6
	}
	speed := mps * conversion
	str := strconv.Itoa(int(speed))

	return renderASCII(str + " " + unit)
}

func renderASCII(text string) string {
	result := ""
	for i := 0; i < 5; i++ {
		result += "  "
		for _, char := range text {
			result += characters[char][i] + "  "
		}
		result += "\n"
	}
	return result
}
