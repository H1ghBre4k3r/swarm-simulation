package process

import (
	"bufio"
	"io"
	"os/exec"
	"sync"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

type Process struct {
	c  *exec.Cmd
	si io.WriteCloser
	so io.ReadCloser

	inputChannel  chan string
	outputChannel chan string

	Out <-chan string
	In  chan<- string

	running bool
	lock    sync.Mutex
}

func Spawn(command string, args ...string) (*Process, error) {
	p := &Process{}
	p.c = exec.Command(command, args...)
	si, err := p.c.StdinPipe()
	if err != nil {
		return nil, err
	}
	p.si = si

	so, err := p.c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	p.so = so

	p.inputChannel = make(chan string)
	p.outputChannel = make(chan string)

	p.In = p.inputChannel
	p.Out = p.outputChannel

	return p, nil
}

func (p *Process) Start() error {
	pipe, err := p.c.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = p.c.Start()
	if err != nil {
		return err
	}
	p.running = true
	go func() {
		p.c.Wait()
		p.running = false
	}()
	reader := bufio.NewReader(pipe)
	go func() {
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				// panic(err)
				p.Stop()
				break
			}
			println(string(line))
		}
	}()

	go p.writeStdIn()
	go p.readStdOut()
	return nil
}

func (p *Process) writeStdIn() {
	defer p.si.Close()
	for {
		message, ok := <-p.inputChannel
		if !ok {
			p.Stop()
			break
		}

		if message[len(message)-1] != '\n' {
			message += "\n"
		}

		_, err := p.si.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}
}

func (p *Process) readStdOut() {
	defer p.so.Close()
	reader := bufio.NewReader(p.so)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			p.Stop()
			break
		}
		p.outputChannel <- string(line)
	}
}

func (p *Process) Stop() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !util.IsChannelClosed(p.inputChannel) {
		close(p.inputChannel)
	}
	if !util.IsChannelClosed(p.outputChannel) {
		close(p.outputChannel)
	}
	return p.c.Process.Kill()
}

func (p *Process) IsRunning() bool {
	return p.running
}
