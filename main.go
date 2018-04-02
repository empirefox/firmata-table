package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/empirefox/firmata"
	"github.com/empirefox/firmata-table/stm32f407vet6"
)

var (
	addr = flag.String("addr", "", "firmata server address")

	viewBoard = stm32f407vet6.Board
)

// use this to print Markdown version PIN table.
func main() {
	flag.Parse()
	if *addr == "" {
		flag.PrintDefaults()
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	var d net.Dialer
	sp, err := d.DialContext(ctx, "tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer sp.Close()

	board := firmata.NewFirmata(sp)

	board.OnConnected = func() {
		fmt.Println("Connected to board")
		fmt.Println("Protocol", board.ProtocolVersion.Server.Name)
		fmt.Println("Firmware", string(board.FirmwareName), board.FirmwareVersion.Server.Name)
		fmt.Println("Pins total:", len(board.Pins_l))
		fmt.Println("AnalogPins total:", len(board.AnalogPins_l))
	}

	pinStateId := 0
	board.OnPinState = func(pin *firmata.Pin) {
		fmt.Print(pinStateId, " ")
		pinStateId++
		if pinStateId > len(board.Pins_l) {
			fmt.Printf("OnPinState: %#v\n", pin)
		} else if pinStateId == len(board.Pins_l) {
			fmt.Println("")
			fmt.Println("OnPinState finished")
		}
	}
	board.OnAnalogMessage = func(pin *firmata.Pin) {
		fmt.Printf("OnAnalogMessage: %#v\n", pin)
	}

	port := 0
	board.OnDigitalMessage = func(pins []*firmata.Pin) {
		fmt.Print(port, " ")
		port++
		if port > 9 {
		} else if port == 9 {
			fmt.Printf("OnDigitalMessage: %#v\n", pins)
			fmt.Println("OnDigitalMessage finished:")
			tables, err := viewBoard.HeadersToMarkdownTables(board.Pins_l)
			if err != nil {
				log.Fatal(err)
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

	fmt.Println("connecting.....")
	ctx, _ = context.WithTimeout(context.Background(), time.Second*20)
	err = board.Handshake(ctx)
	if err != nil {
		log.Fatal(err)
	}

	board.Loop(func() {
		fmt.Println("PinStateQuery_l:")
		var i byte
		total := byte(len(board.Pins_l))
		for i = 0; i < total; i++ {
			if err = board.PinStateQuery_l(i); err != nil {
				log.Fatal(err)
			}
		}
	})

	board.Loop(func() {
		fmt.Println("ReportDigital_l:")
		var i byte
		for i = 0; i < 9; i++ {
			if err = board.ReportDigital_l(i, 1); err != nil {
				log.Fatal(err)
			}
			if err = board.ReportDigital_l(i, 0); err != nil {
				log.Fatal(err)
			}
		}
	})

	var pin1 byte = 66
	var pin2 byte = 67
	board.Loop(func() {
		if err = board.SetPinMode_l(pin1, firmata.PIN_MODE_OUTPUT); err != nil {
			log.Fatal(err)
		}
		if err = board.SetPinMode_l(pin2, firmata.PIN_MODE_OUTPUT); err != nil {
			log.Fatal(err)
		}
	})

	//	if err := board.ReportAnalog(1, 3); err != nil {
	//		panic(err)
	//	}

	go func() {
		var level byte
		for {
			board.Loop(func() {
				if err := board.SetDigitalPinValue_l(pin1, level); err != nil {
					log.Fatal(err)
				}
				level ^= 1
				if err := board.SetDigitalPinValue_l(pin2, level); err != nil {
					log.Fatal(err)
				}
			})

			//			fmt.Println("PD15:", level)
			time.Sleep(3000 * time.Millisecond)
		}
	}()

	<-board.CloseNotifier()
	err = board.ClosedError
	if err != nil {
		log.Fatal(err)
	}
}
