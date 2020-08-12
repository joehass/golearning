package command

import (
	"fmt"
	"testing"
)

//命令模式

type TV struct {
}

func (p TV) Open() {
	fmt.Println("play...")
}

func (p TV) Close() {
	fmt.Println("stop...")
}

type Command interface {
	Press()
}

type OpenCommand struct {
	tv TV
}

func (p OpenCommand) Press() {
	p.tv.Open()
}

type CloseCommand struct {
	tv TV
}

func (p CloseCommand) Press() {
	p.tv.Close()
}

type Invoker struct {
	cmd Command
}

func (p *Invoker) SetCommand(cmd Command) {
	p.cmd = cmd
}

func (p *Invoker) Do() {
	p.cmd.Press()
}

func TestComm(t *testing.T) {
	var tv TV
	OpenCommand := OpenCommand{tv}
	invoker := Invoker{OpenCommand}
	invoker.Do()

	invoker.SetCommand(CloseCommand{tv})
	invoker.Do()
}
