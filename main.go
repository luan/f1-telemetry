package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/luan/f1-telemetry/f1"
)

func main() {
	dataChan := make(chan f1.TelemetryData, 1000)
	ui := NewUI(dataChan)

	go serveTelemetry(dataChan)
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
