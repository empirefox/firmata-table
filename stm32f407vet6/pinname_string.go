// Code generated by "stringer -type=PinName"; DO NOT EDIT.

package stm32f407vet6

import "fmt"

const _PinName_name = "PE2PE3PE4PE5PE6PC13PC0PC1PC2PC3PA0PA1PA2PA3PA4PA5PA6PA7PC4PC5PB0PB1PE7PE8PE9PE10PE11PE12PE13PE14PE15PB10PB11PB12PB13PB14PE1PE0PB9PB8PB7PB6PB5PB3PD7PD6PD5PD4PD3PD2PD1PD0PC12PC11PC10PA15PA12PA11PA10PA9PA8PC9PC8PC7PC6PD15PD14PD13PD12PD11PD10PD9PD8PB15PA13PA14PB4PB2PC0_2PC1_2PC2_2PC3_2PA0_2PA1_2PA2_2PA3_2PC4_2PC5_2PB0_2PB1_2PEND"

var _PinName_index = [...]uint16{0, 3, 6, 9, 12, 15, 19, 22, 25, 28, 31, 34, 37, 40, 43, 46, 49, 52, 55, 58, 61, 64, 67, 70, 73, 76, 80, 84, 88, 92, 96, 100, 104, 108, 112, 116, 120, 123, 126, 129, 132, 135, 138, 141, 144, 147, 150, 153, 156, 159, 162, 165, 168, 172, 176, 180, 184, 188, 192, 196, 199, 202, 205, 208, 211, 214, 218, 222, 226, 230, 234, 238, 241, 244, 248, 252, 256, 259, 262, 267, 272, 277, 282, 287, 292, 297, 302, 307, 312, 317, 322, 326}

func (i PinName) String() string {
	if i < 0 || i >= PinName(len(_PinName_index)-1) {
		return fmt.Sprintf("PinName(%d)", i)
	}
	return _PinName_name[_PinName_index[i]:_PinName_index[i+1]]
}
