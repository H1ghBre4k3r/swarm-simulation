package process

import (
	"bufio"
	"io"
	"os/exec"
)

type Process struct {
	c *exec.Cmd
	// si io.WriteCloser
	// so io.ReadCloser

	inputChannel  chan string
	outputChannel chan string

	Out <-chan string
	In  chan<- string
}

func Spawn(command string, args ...string) (*Process, error) {
	p := &Process{}
	p.c = exec.Command(command, args...)

	si, err := p.c.StdinPipe()
	if err != nil {
		return nil, err
	}

	so, err := p.c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	p.inputChannel = make(chan string)
	p.outputChannel = make(chan string)

	p.In = p.inputChannel
	p.Out = p.outputChannel

	err = p.c.Start()
	if err != nil {
		return nil, err
	}

	go p.writeStdIn(si)
	go p.readStdOut(so)

	return p, nil
}

func (p *Process) writeStdIn(si io.WriteCloser) {
	defer si.Close()
	for {
		message, ok := <-p.inputChannel
		if !ok {
			p.Stop()
			break
		}

		if message[len(message)-1] != '\n' {
			message += "\n"
		}

		_, err := si.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}
}

func (p *Process) readStdOut(so io.ReadCloser) {
	defer so.Close()
	reader := bufio.NewReader(so)
	for {
		answer, err := reader.ReadString('\n')
		if err != nil {
			p.Stop()
			break
		}
		p.outputChannel <- answer
	}
}

func (p *Process) Stop() error {
	if !IsClosed(p.inputChannel) {
		close(p.inputChannel)
	}
	if !IsClosed(p.outputChannel) {
		close(p.outputChannel)
	}
	return p.c.Process.Kill()
}

func IsClosed(ch chan string) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func (p *Process) IsRunning() bool {
	return !p.c.ProcessState.Exited()
}
