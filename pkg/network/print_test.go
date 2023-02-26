package network

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/AlexEkdahl/snakes/pkg/network/protobuf"
)

func TestPrinter_PrintMessage(t *testing.T) {
	type args struct {
		s interface{}
	}

	msg := &protobuf.Message{
		Type: &protobuf.Message_Move{
			Move: &protobuf.MoveMessage{
				Direction: protobuf.Direction(0),
			},
		},
	}
	tests := []struct {
		name    string
		p       *Printer
		args    args
		wantErr bool
	}{
		{
			name: "Prints string successfully",
			p:    NewPrinter(),
			args: args{s: "test message"},
		},
		{
			name: "Prints []byte successfully",
			p:    NewPrinter(),
			args: args{s: []byte("test message")},
		},
		{
			name: "Prints proto.Message successfully",
			p:    NewPrinter(),
			args: args{s: msg},
		},
		{
			name: "Prints struct successfully",
			p:    NewPrinter(),
			args: args{s: struct{ Field string }{Field: "test message"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.p.out = &buf

			tt.p.PrintMessage(tt.args.s)

			if got := buf.String(); len(got) == 0 {
				t.Errorf("Printer.PrintMessage() = %v, want non-empty string", got)
			}
		})
	}
}

func TestPrinter_Clear(t *testing.T) {
	tests := []struct {
		name    string
		p       *Printer
		wantErr bool
	}{
		{
			name: "Clears the screen successfully",
			p:    NewPrinter(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.p.out = &buf

			tt.p.Clear()

			if got := buf.String(); !reflect.DeepEqual(got, "\033[2J\033[H") {
				t.Errorf("Printer.Clear() = %v, want \033[2J\033[H", got)
			}
		})
	}
}
