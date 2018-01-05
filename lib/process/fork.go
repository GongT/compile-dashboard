package process

import (
	"os/exec"
	"bytes"
	"github.com/gongt/compile-dashboard/lib"
	"fmt"
	"os"
	"time"
)

type ChildProcess struct {
	cmd        *exec.Cmd
	title      string
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
		script,
		pipe,
		make(chan error, 1),
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
		close(cp.Stop)
		lib.MainLogger.Println("subprocess end: ", script)
	}()

	return &cp
}

func (cp *ChildProcess) Inspect() string {
	return fmt.Sprintf("child process: %d", cp.cmd.Process.Pid)
}
func (cp *ChildProcess) Close() (err error) {
	pid := fmt.Sprint(cp.cmd.Process.Pid)
	lib.MainLogger.Printf("stop child process: %s: %s", pid, cp.title)
	if cp.cmd.ProcessState != nil && cp.cmd.ProcessState.Exited() {
		return
	}

	cp.cmd.Process.Signal(os.Interrupt)
	select {
	case <-time.After(3 * time.Second):
		err = cp.cmd.Process.Kill()
		lib.MainLogger.Println("process[" + pid + "] killed as timeout reached")
		if err != nil {
			lib.MainLogger.Println("process["+pid+"] can not kill:", err)
			return
		}
	case err = <-cp.Stop:
		if err != nil {
			lib.MainLogger.Println("process["+pid+"] done with error = %v", err)
		} else {
			lib.MainLogger.Println("process[" + pid + "] done gracefully without error")
		}
	}
	return
}
