package network

import (
	"encoding/json"
	"io"
	"os"

	"google.golang.org/protobuf/proto"
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
		_, err := p.out.Write([]byte(v))
		if err != nil {
			return
		}
	case []byte:
		_, err := p.out.Write(v)
		if err != nil {
			return
		}
	case proto.Message:
		data, err := proto.Marshal(v)
		if err != nil {
			return
		}
		_, err = p.out.Write(data)
		if err != nil {
			return
		}
	case interface{}:
		data, err := json.Marshal(v)
		if err != nil {
			return
		}

		_, err = p.out.Write(data)
		if err != nil {
			return
		}
	}
}

func (p *Printer) Clear() {
	_, err := p.out.Write([]byte("\033[2J\033[H"))
	if err != nil {
		return
	}
}
