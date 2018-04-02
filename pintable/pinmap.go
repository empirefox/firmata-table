//go:generate stringer -type=PinMode
package pintable

import (
	"bytes"
	"fmt"
	"sort"
	"text/template"

	"github.com/empirefox/firmata"
)

type PinMode byte

const (
	I       PinMode = 0x00 // same as INPUT defined in Arduino.h
	O       PinMode = 0x01 // same as OUTPUT defined in Arduino.h
	A       PinMode = 0x02 // analog pin in analogInput mode
	PWM     PinMode = 0x03 // digital pin in PWM output mode
	SERVO   PinMode = 0x04 // digital pin in Servo output mode
	SHIFT   PinMode = 0x05 // shiftIn/shiftOut mode
	I2C     PinMode = 0x06 // pin included in I2C setup
	ONEWIRE PinMode = 0x07 // pin configured for 1-wire
	STEPPER PinMode = 0x08 // pin configured for stepper motor
	ENCODER PinMode = 0x09 // pin configured for rotary encoders
	SERIAL  PinMode = 0x0A // pin configured for serial communication
	PULLUP  PinMode = 0x0B // enable internal pull-up resistor for pin
	X       PinMode = 0x7F // pin configured to be ignored by digitalWrite and capabilityResponse
)

func NewPinMode(mode byte) fmt.Stringer {
	return (interface{})(PinMode(mode)).(fmt.Stringer)
}

type Header struct {
	Name    string
	Padding map[int]string
	Total   int
	Start   int
}

type Board struct {
	Headers    []Header
	AnalogPins []int
	Stringer   func(pid int) fmt.Stringer
	PinEnd     int
}

type headerPin struct {
	ID            int
	Padding       bool
	Name          string
	Modes         []string
	Mode          string
	Value         int
	State         int
	AnalogChannel byte
	Digital       int
}

const (
	oddHeaderPinStr     = `| {{.AnalogChannel}} | {{.Modes}} | {{.Mode}} | {{.State}} | {{.Value}} | {{.Digital}} | {{.Name}} | {{.ID}} `
	evenHeaderPinStr    = `{{.ID}} | {{.Name}} | {{.Digital}} | {{.Value}} | {{.State}} | {{.Mode}} | {{.Modes}} | {{.AnalogChannel}} |`
	oddHeaderPadPinStr  = `| | | | | | | {{.Name}} | {{.ID}} `
	evenHeaderPadPinStr = `{{.ID}} | {{.Name}} | | | | | | |`
)

var (
	oddHeaderPinTpl     = template.Must(template.New("oddHeaderPinTpl").Parse(oddHeaderPinStr))
	evenHeaderPinTpl    = template.Must(template.New("evenHeaderPinTpl").Parse(evenHeaderPinStr))
	oddHeaderPadPinTpl  = template.Must(template.New("oddHeaderPadPinTpl").Parse(oddHeaderPadPinStr))
	evenHeaderPadPinTpl = template.Must(template.New("evenHeaderPadPinTpl").Parse(evenHeaderPadPinStr))
)

func (board *Board) HeadersToMarkdownTables(allpins []*firmata.Pin) ([][]byte, error) {
	var tables [][]byte
	for _, header := range board.Headers {
		b, err := board.HeaderToMarkdownTable(allpins, header)
		if err != nil {
			return nil, err
		}
		tables = append(tables, b)
	}
	return tables, nil
}

func (board *Board) HeaderToMarkdownTable(allpins []*firmata.Pin, header Header) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("|A|Modes|M|S|V|D|N|")    // oddHeaderPinStr
	b.WriteString(header.Name)              // header name
	b.WriteString("_ID|N|D|V|S|M|Modes|A|") // evenHeaderPinStr
	b.WriteByte('\n')
	b.WriteString("|:--:|:--|:--|:--|:--|:--|:--|:--:|:--|:--|:--|:--|:--|:--|:--:|")
	b.WriteByte('\n')

	pid := header.Start
	var hp *headerPin
	for i := 0; i < header.Total; i++ {
		fid := i + 1
		pad, ok := header.Padding[fid]
		if ok {
			hp = &headerPin{
				ID:      fid,
				Padding: true,
				Name:    pad,
				Digital: -1,
			}
		} else {
			pin := allpins[pid]
			ms := make([]int, len(pin.Modes))
			for m := range pin.Modes {
				ms = append(ms, int(m))
			}
			sort.Ints(ms)
			modes := make([]string, len(pin.Modes))
			for m := range ms {
				modes = append(modes, NewPinMode(byte(m)).String())
			}
			hp = &headerPin{
				ID:            fid,
				Name:          board.Stringer(pid).(fmt.Stringer).String(),
				Modes:         modes,
				Mode:          NewPinMode(pin.Mode).String(),
				Value:         pin.Value,
				State:         pin.State,
				AnalogChannel: pin.AnalogChannel,
				Digital:       pid,
			}
			for aid, atop := range board.AnalogPins {
				if atop == pid {
					hp.AnalogChannel = byte(aid)
					break
				}
			}
			pid++
		}

		var tpl *template.Template
		odd := fid&1 == 1
		if odd {
			// odd
			if hp.Padding {
				tpl = oddHeaderPadPinTpl
			} else {
				tpl = oddHeaderPinTpl
			}
		} else {
			// even
			if hp.Padding {
				tpl = evenHeaderPadPinTpl
			} else {
				tpl = evenHeaderPinTpl
			}
		}
		if err := tpl.Execute(&b, hp); err != nil {
			return nil, err
		}
		if !odd {
			b.WriteByte('\n')
		}
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}
