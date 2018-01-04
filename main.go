package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/empirefox/firmata"
	"github.com/empirefox/firmata-table/stm32f407vet6"
)

// use this to print Markdown version PIN table.
func main() {
	viewBoard := stm32f407vet6.Board
	sp, err := net.Dial("tcp", "192.168.20.165:3030")
	if err != nil {
		panic(err)
	}

	board := firmata.NewFirmata()

	pinStateId := 0
	board.OnPinState = func(pin *firmata.Pin) {
		fmt.Print(pinStateId, " ")
		pinStateId++
		if pinStateId > len(board.Pins()) {
			fmt.Printf("OnPinState: %#v\n", pin)
		} else if pinStateId == len(board.Pins()) {
			fmt.Println("")
			fmt.Println("OnPinState finished")
		}
	}
	board.OnConnected = func() {
		fmt.Println("Connected to board")
		fmt.Println("Protocol", board.ProtocolVersion.Server.Name)
		fmt.Println("Firmware", board.FirmwareName, board.FirmwareVersion.Server.Name)
		fmt.Println("Pins total:", len(board.Pins()))
		fmt.Println("AnalogPins total:", len(board.AnalogPins()))
	}
	board.OnAnalogMessage = func(pin *firmata.Pin) {
		fmt.Printf("OnAnalogMessage: %#v\n", pin)
	}

	port := 0
	board.OnDigitalMessage = func(pins []*firmata.Pin) {
		fmt.Print(port, " ")
		port++
		if port > 0 {
		} else if port == 9 {
			fmt.Printf("OnDigitalMessage: %#v\n", pins)
			fmt.Println("OnDigitalMessage finished:")
			tables, err := viewBoard.HeadersToMarkdownTables(board.Pins())
			if err != nil {
				panic(err)
			}
			for _, md := range tables {
				fmt.Print(string(md))
			}
		}
	}
	board.OnI2cReply = nil
	board.OnStringData = func(data []byte) {
		fmt.Printf("OnStringData: %s\n", data)
	}
	board.OnSysexResponse = func(buf []byte) {
		fmt.Printf("OnSysexResponse: %#v\n", buf)
	}
	board.OnError = func(err error) {
		fmt.Println("OnError:", err)
	}

	fmt.Println("connecting.....")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	err = board.Dial(ctx, sp)
	defer board.Close()

	if err != nil {
		panic(err)
	}

	fmt.Println("OnPinState:")
	for i := 0; i < len(board.Pins()); i++ {
		if err = board.PinStateQuery(i); err != nil {
			panic(err)
		}
	}

	fmt.Println("OnDigitalMessage:")
	for i := 0; i < 9; i++ {
		if err = board.ReportDigital(i, 1); err != nil {
			panic(err)
		}
		if err = board.ReportDigital(i, 0); err != nil {
			panic(err)
		}
	}

	pin1 := 66
	pin2 := 67
	if err = board.SetPinMode(pin1, firmata.PIN_MODE_OUTPUT); err != nil {
		panic(err)
	}
	if err = board.SetPinMode(pin2, firmata.PIN_MODE_OUTPUT); err != nil {
		panic(err)
	}

	//	if err := board.ReportAnalog(1, 3); err != nil {
	//		panic(err)
	//	}

	level := 0

	for {
		if err := board.SetDigitalPinValue(pin1, level); err != nil {
			panic(err)
		}
		level ^= 1
		if err := board.SetDigitalPinValue(pin2, level); err != nil {
			panic(err)
		}

		//		fmt.Println("PD15:", level)
		time.Sleep(3000 * time.Millisecond)
	}
}
