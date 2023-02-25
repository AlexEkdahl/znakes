package network

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/proto"
)

type Printer struct {
	out io.Writer
}

func NewPrinter() *Printer {
	return &Printer{os.Stdout}
}

func (p *Printer) PrintMessage(s interface{}) {
	p.Clear()
	switch v := s.(type) {
	case string:
		p.out.Write([]byte(v))
	case []byte:
		p.out.Write(v)
	case proto.Message:
		data, err := proto.Marshal(v)
		if err == nil {
			p.out.Write(data)
		}
	case interface{}:
		data, err := json.Marshal(v)
		if err == nil {
			p.out.Write(data)
		}
	default:
		data, err := json.Marshal(v)
		if err == nil {
			p.out.Write(data)
		} else {
			fmt.Printf("Cannot print %T: %v\n", v, err)
		}
	}
}

func (p *Printer) Clear() {
	p.out.Write([]byte("\033[2J\033[H"))
}
