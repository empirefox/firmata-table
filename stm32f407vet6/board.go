//go:generate stringer -type=PinName
package stm32f407vet6

import (
	"fmt"

	"github.com/empirefox/firmata-table/pintable"
)

// Toggle this to generate PinName
type PinName int

const (
	//J2 connector Left, Right side - USB side,
	PE2 PinName = iota //D0

	PE3  //D1 - BUTTON K1 (USE INPUT_PULLUP)
	PE4  //D2 - BUTTON K0 (USE INPUT_PULLUP)
	PE5  //D3
	PE6  //D4
	PC13 //D5
	PC0  //D6
	PC1  //D7
	PC2  //D8 - SPI2
	PC3  //D9 - SPI2
	PA0  //D10 - BUTTON WK_UP (USE INPUT_PULLDOWN for GP button)
	PA1  //D11
	PA2  //D12
	PA3  //D13
	PA4  //D14
	PA5  //D15 - SPI1 (or SPI3)
	PA6  //D16 - BUILTIN LED D2 - SPI1 (or SPI3)
	PA7  //D17 - BUILTIN LED D3 - SPI1 (or SPI3)
	PC4  //D18
	PC5  //D19
	PB0  //D20 - Flash_CS
	PB1  //D21 - LCD_BL driver
	PE7  //D22
	PE8  //D23
	PE9  //D24 - SPI2
	PE10 //D25 - SPI2
	PE11 //D26 - SPI2
	PE12 //D27
	PE13 //D28
	PE14 //D29
	PE15 //D30
	PB10 //D31 - I2C2
	PB11 //D32 - I2C2
	PB12 //D33 - SPI2
	PB13 //D34 - SPI2
	PB14 //D35 - SPI3 (or SPI1) - FLASH/NRF

	//J3 connector Left, Right side - SDIO side
	PE1  //D36
	PE0  //D37
	PB9  //D38
	PB8  //D39 - I2C1
	PB7  //D40 - I2C1
	PB6  //D41
	PB5  //D42 - SPI3 (or SPI1) - FLASH/NRF
	PB3  //D43 - SPI3 (or SPI1) - FLASH/NRF
	PD7  //D44
	PD6  //D45 - USART2
	PD5  //D46 - USART2
	PD4  //D47
	PD3  //D48
	PD2  //D49 - SDIO
	PD1  //D50
	PD0  //D51
	PC12 //D52 - SDIO
	PC11 //D53 - SDIO
	PC10 //D54 - SDIO
	PA15 //D55
	PA12 //D56 - USB DM
	PA11 //D57 - USB DP
	PA10 //D58 - RX1
	PA9  //D59 - TX1
	PA8  //D60
	PC9  //D61 - SDIO
	PC8  //D62 - SDIO
	PC7  //D63
	PC6  //D64
	PD15 //D65
	PD14 //D66
	PD13 //D67
	PD12 //D68
	PD11 //D69
	PD10 //D70
	PD9  //D71 - USART3
	PD8  //D72 - USART3
	PB15 //D73
	PA13 //D74 SWDIO (JTAG)
	PA14 //D75 SWCLK (JTAG)
	PB4  //D76 CONNECTED TO NRF HEADER Pin 7, Winbond Flash Pin 2, NRST JTAG Pin 3
	PB2  //D77 BOOT1 - J3

	//  PC14,	// XTAL
	//  PC15, // XTAL

	// Duplicated to have A0-A5 as F407 do not have Uno like connector
	// and to be aligned with PinMap_ADC
	PC0_2 //D78/A0
	PC1_2 //D79/A1
	PC2_2 //D80/A2 - SPI2
	PC3_2 //D81/A3 - SPI2
	PA0_2 //D82/A4 - BUTTON WK_UP (USE INPUT_PULLDOWN for GP button)
	PA1_2 //D83/A5
	PA2_2 //D84/A6
	PA3_2 //D85/A7
	PC4_2 //D86/A8
	PC5_2 //D87/A9
	PB0_2 //D88/A10 - Flash_CS
	PB1_2 //D89/A11 - LCD_BL driver
	PEND
)

var headers = []pintable.Header{
	{
		Name: "J2",
		Padding: map[int]string{
			1:  "5V",
			2:  "5V",
			3:  "5V",
			4:  "5V",
			5:  "3V3",
			6:  "3V3",
			7:  "3V3",
			8:  "3V3",
			9:  "GND",
			10: "GND",
			21: "VREF-",
			22: "VREF+",
		},
		Total: 48,
		Start: int(PE2),
	},
	{
		Name: "J3",
		Padding: map[int]string{
			1:  "3V3",
			2:  "3V3",
			3:  "3V3",
			4:  "3V3",
			5:  "BOOT0",
			6:  "BOOT1",
			7:  "GND",
			8:  "GND",
			9:  "GND",
			10: "GND",
		},
		Total: 48,
		Start: int(PE1),
	},
}

var analogPins = []int{
	int(PC0), //D78/A0
	int(PC1), //D79/A1
	int(PC2), //D80/A2 - SPI2
	int(PC3), //D81/A3 - SPI2
	int(PA0), //D82/A4 - BUTTON WK_UP (USE INPUT_PULLDOWN for GP button)
	int(PA1), //D83/A5
	int(PA2), //D84/A6
	int(PA3), //D85/A7
	int(PC4), //D86/A8
	int(PC5), //D87/A9
	int(PB0), //D88/A10 - Flash_CS
	int(PB1), //D89/A11 - LCD_BL driver
}

var stringer = func(pid int) fmt.Stringer {
	var pn interface{} = PinName(pid)
	return pn.(fmt.Stringer)
}

var Board = &pintable.Board{
	Headers:    headers,
	AnalogPins: analogPins,
	Stringer:   stringer,
	PinEnd:     int(PEND),
}

// TODO delete
var Tags = map[PinName]string{
	PA15: "SPL1", // SS1 SS3
	PB9:  "SPL2", // SS2
	PB12: "SPL2", // SS2

	PA4: "SPL1", // SS1 Default, SS3
	PA5: "SPL1", // PA5
	PA6: "SPL1", // MISO
	PA7: "SPL1", // MOSI

	PB5:  "SPL3",
	PB4:  "SPL3",
	PB3:  "SPL3",
	PB14: "SPL2",
	PB13: "SPL2",
	PB10: "SPL2",
	PC3:  "SPL2",
	PC2:  "SPL2",

	PA9:  "TX",
	PA10: "RX",
}
