package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/luan/f1-telemetry/f1"
)

const (
	DBName = "f1telemetry"
)

func main() {
	dataChan := make(chan f1.TelemetryData, 1000)

	go serveTelemetry(dataChan)

	uiDataChan := make(chan f1.TelemetryData, 1000)

	influxDataChan := make(chan f1.TelemetryData, 1000)
	go influx(influxDataChan)

	go func() {
		for data := range dataChan {
			uiDataChan <- data
			influxDataChan <- data
		}
	}()

	ui := NewUI(uiDataChan)
	ui.Start()
}

func serveTelemetry(dataChan chan<- f1.TelemetryData) {
	serverAddr, err := net.ResolveUDPAddr("udp", ":20777")
	if err != nil {
		log.Fatal(err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer serverConn.Close()

	for {
		var telemetry f1.TelemetryData
		err := binary.Read(serverConn, binary.LittleEndian, &telemetry)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		dataChan <- telemetry
	}
}

func influx(dataChan <-chan f1.TelemetryData) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		log.Fatal(err)
	}

	for data := range dataChan {
		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  DBName,
			Precision: "us",
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create a point and add to batch
		tags := map[string]string{"driver": "self"}
		fields := map[string]interface{}{
			"time":                    data.Speed,
			"laptime":                 data.Laptime,
			"lapdistance":             data.Lapdistance,
			"totaldistance":           data.Totaldistance,
			"x":                       data.X,
			"y":                       data.Y,
			"z":                       data.Z,
			"speed":                   data.Speed,
			"xv":                      data.Xv,
			"yv":                      data.Yv,
			"zv":                      data.Zv,
			"xr":                      data.Xr,
			"yr":                      data.Yr,
			"zr":                      data.Zr,
			"xd":                      data.Xd,
			"yd":                      data.Yd,
			"zd":                      data.Zd,
			"susp-pos":                data.SuspPos,
			"susp-vel":                data.SuspVel,
			"wheel-speed":             data.WheelSpeed,
			"throttle":                data.Throttle,
			"steer":                   data.Steer,
			"brake":                   data.Brake,
			"clutch":                  data.Clutch,
			"gear":                    data.Gear,
			"gforce-lat":              data.GforceLat,
			"gforce-lon":              data.GforceLon,
			"lap":                     data.Lap,
			"enginerate":              data.Enginerate,
			"sli-pro-native-support":  data.SliProNativeSupport,
			"car-position":            data.CarPosition,
			"kers-level":              data.KersLevel,
			"kers-max-level":          data.KersMaxLevel,
			"drs":                     data.DRS,
			"traction-control":        data.TractionControl,
			"anti-lock-brakes":        data.AntiLockBrakes,
			"fuel-in-tank":            data.FuelInTank,
			"fuel-capacity":           data.FuelCapacity,
			"in-pits":                 data.InPits,
			"sector":                  data.Sector,
			"sector1-time":            data.Sector1Time,
			"sector2-time":            data.Sector2Time,
			"brakes-temp":             data.BrakesTemp,
			"tyres-pressure":          data.TyresPressure,
			"team-info":               data.TeamInfo,
			"total-laps":              data.TotalLaps,
			"track-size":              data.TrackSize,
			"last-lap-time":           data.LastLapTime,
			"max-rpm":                 data.MaxRpm,
			"idle-rpm":                data.IdleRpm,
			"max-gears":               data.MaxGears,
			"session-type":            data.SessionType,
			"drsallowed":              data.Drsallowed,
			"track-number":            data.TrackNumber,
			"vehiclefiaflags":         data.Vehiclefiaflags,
			"era":                     data.Era,
			"engine-temperature":      data.EngineTemperature,
			"gforce-vert":             data.GforceVert,
			"ang-vel-x":               data.AngVelX,
			"ang-vel-y":               data.AngVelY,
			"ang-vel-z":               data.AngVelZ,
			"tyres-temperature":       data.TyresTemperature,
			"tyres-wear":              data.TyresWear,
			"tyre-compound":           data.TyreCompound,
			"front-brake-bias":        data.FrontBrakeBias,
			"fuel-mix":                data.FuelMix,
			"currentlapinvalid":       data.Currentlapinvalid,
			"tyres-damage":            data.TyresDamage,
			"front-left-wing-damage":  data.FrontLeftWingDamage,
			"front-right-wing-damage": data.FrontRightWingDamage,
			"rear-wing-damage":        data.RearWingDamage,
			"engine-damage":           data.EngineDamage,
			"gear-box-damage":         data.GearBoxDamage,
			"exhaust-damage":          data.ExhaustDamage,
			"pit-limiter-status":      data.PitLimiterStatus,
			"pit-speed-limit":         data.PitSpeedLimit,
			"session-time-left":       data.SessionTimeLeft,
			"rev-lights-percent":      data.RevLightsPercent,
			"is-spectating":           data.IsSpectating,
			"spectator-car-index":     data.SpectatorCarIndex,
			"num-cars":                data.NumCars,
			"player-car-index":        data.PlayerCarIndex,
		}

		pt, err := client.NewPoint("telemetry", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		for _, car := range data.Cars {
			tags := map[string]string{"driver": strconv.Itoa(int(car.DriverID))}
			fields := map[string]interface{}{
				"lastlap-time":      car.LastlapTime,
				"currentlap-time":   car.CurrentlapTime,
				"bestlap-time":      car.BestlapTime,
				"sector1-time":      car.Sector1Time,
				"sector2-time":      car.Sector2Time,
				"lap-distance":      car.LapDistance,
				"driver-id":         car.DriverID,
				"team-id":           car.TeamID,
				"car-position":      car.CarPosition,
				"current-lap-num":   car.CurrentLapNum,
				"tyre-compound":     car.TyreCompound,
				"in-pits":           car.InPits,
				"sector":            car.Sector,
				"currentlapinvalid": car.Currentlapinvalid,
				"penalties":         car.Penalties,
			}

			pt, err := client.NewPoint("car", tags, fields, time.Now())
			if err != nil {
				log.Fatal(err)
			}
			bp.AddPoint(pt)
		}

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}
	}
}
