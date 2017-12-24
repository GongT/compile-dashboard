package process

import (
	"os/exec"
	"bytes"
	"github.com/gongt/compile-dashboard/lib"
)

type ChildProcess struct {
	cmd        *exec.Cmd
	OutputPipe *bufferWriter
	Stop       chan error
	isRunning  bool
}

type bufferWriter struct {
	Output      chan []byte
	Clear       chan bool
	lastFewChar []byte
	needClear   bool
}

const maxSaveLen = 10

func newBufferWriter() *bufferWriter {
	ret := bufferWriter{
		make(chan []byte, 256),
		make(chan bool),
		make([]byte, 0, maxSaveLen),
		false,
	}
	return &ret
}

func (pipe *bufferWriter) Write(p []byte) (n int, err error) {
	outputLength := len(p)
	outputStart := bytes.LastIndex(p, []byte("\x1Bc"))
	if outputStart >= 0 {
		// found \ec in this output
		pipe.needClear = true
		outputStart += 2
	} else {
		lastCharIndex := len(pipe.lastFewChar) - 1
		if lastCharIndex >= 0 && pipe.lastFewChar[lastCharIndex] == 0x1B && p[0] == 'c' {
			// if last end with \e, and this start with c
			pipe.needClear = true
			outputStart = 1
		} else {
			// nothing likes \ec
			outputStart = 0
		}
	}

	outputEnd := outputLength
	if p[outputEnd-1] == 0x1B {
		outputEnd--
	}

	if pipe.needClear {
		pipe.needClear = false
		pipe.Clear <- true
		lib.MainLogger.Println("clear the panel")
	}

	pipe.Output <- p[outputStart:outputEnd]

	if outputLength >= maxSaveLen {
		pipe.lastFewChar = pipe.lastFewChar[:maxSaveLen]
		copy(pipe.lastFewChar, p[outputLength-maxSaveLen:])
	} else {
		pipe.lastFewChar = pipe.lastFewChar[:outputLength]
		copy(pipe.lastFewChar, p)
	}

	return outputLength, nil
}

func NewChildProcess(script string) *ChildProcess {
	cmd := exec.Command("bash", "-c", script)
	pipe := newBufferWriter()

	cmd.Stdout = pipe
	cmd.Stderr = pipe

	cp := ChildProcess{
		cmd,
		pipe,
		make(chan error),
		false,
	}

	go func() {
		cp.isRunning = true
		lib.MainLogger.Println("subprocess start to run: ", cmd)
		err := cmd.Run() // this will block until process exit
		cp.isRunning = false

		cp.Stop <- err

		close(pipe.Output)
		close(pipe.Clear)
		lib.MainLogger.Println("subprocess end: ", script)
	}()

	return &cp
}

func (cp *ChildProcess) Close() error {
	err := cp.cmd.Process.Kill()
	if err != nil {
		return err
	}
	return nil
}
