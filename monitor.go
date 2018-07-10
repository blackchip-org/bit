package bit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	AnsiCursorOff = "\033[?25l"
	AnsiCursorOn  = "\033[?25h"
	AnsiClearLine = "\033[2K"
	AnsiReset     = "\033[0m"
	AnsiWhite     = "\033[37;1m"
	AnsiYellow    = "\033[33m"
	AnsiBlue      = "\033[36m"
	AnsiGreen     = "\033[32m"
	AnsiRed       = "\033[31m"
)

const refreshInterval = time.Millisecond * 100

type Monitor struct {
	taskName  string
	planStart time.Time
	taskStart time.Time
	buf       bytes.Buffer
	out       io.Writer
	updates   chan string
	stop      chan bool
	ticker    *time.Ticker
}

func NewMonitor() *Monitor {
	return &Monitor{
		updates: make(chan string),
		stop:    make(chan bool),
		out:     os.Stdout,
	}
}

func (m *Monitor) Update(line string) {
	m.updates <- line
}

func (m *Monitor) Run() {
	m.ticker = time.NewTicker(refreshInterval)
	defer m.ticker.Stop()

	fmt.Fprint(m.out, AnsiCursorOff)
	defer m.resetTerminal()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGQUIT)
	defer signal.Stop(sig)
	go m.signalHandler(sig)

	for {
		select {
		case line := <-m.updates:
			m.printOutputLine(line)
		case <-m.ticker.C:
			m.flushOutput()
		case <-m.stop:
			continue
		}
	}
}

func (m *Monitor) Stop() {
	m.stop <- true
	m.flushOutput()
	m.printStatusLine(AnsiGreen)
	fmt.Fprintln(m.out)
	m.resetTerminal()
}

func (m *Monitor) Error() {
	m.stop <- true
	m.flushOutput()
	m.printStatusLine(AnsiRed)
	fmt.Fprintln(m.out)
	m.resetTerminal()
}

func (m *Monitor) SetTaskName(name string) {
	now := time.Now().Round(time.Second)
	if m.taskName == "" {
		m.planStart = now
	} else {
		m.printStatusLine(AnsiGreen)
		fmt.Fprintln(m.out)
	}
	m.taskName = name
	m.taskStart = now
}

func (m *Monitor) printOutputLine(line string) {
	lineColor := AnsiBlue
	if line[0] == '+' {
		lineColor = AnsiYellow
	}
	m.buf.WriteString(AnsiReset)
	m.buf.WriteString(lineColor)
	m.buf.WriteString(line)
}

func (m *Monitor) printStatusLine(ansiColor string) {
	now := time.Now().Round(time.Second)
	elapsedPlan := now.Sub(m.planStart)
	elapsedTask := now.Sub(m.taskStart)

	planTime := formatElapsed(elapsedPlan)
	taskTime := formatElapsed(elapsedTask)
	times := ""
	if taskTime == planTime {
		times = planTime
	} else {
		times = taskTime + " " + planTime
	}
	width, err := terminalWidth()
	if err != nil {
		width = 80
	}
	nspaces := width - len(m.taskName) - len(times) - 1 // leave space on end
	pad := strings.Repeat(" ", nspaces)

	fmt.Fprint(m.out, AnsiReset)
	fmt.Fprint(m.out, AnsiClearLine)
	fmt.Fprint(m.out, "\r")
	fmt.Fprint(m.out, ansiColor)
	fmt.Fprintf(m.out, m.taskName+pad+times)
}

func (m *Monitor) flushOutput() {
	fmt.Fprint(m.out, AnsiReset)
	fmt.Fprint(m.out, AnsiClearLine)
	fmt.Fprint(m.out, "\r")
	fmt.Fprint(m.out, m.buf.String())
	m.printStatusLine(AnsiWhite)
	m.buf.Reset()
}

func (m *Monitor) signalHandler(sig chan os.Signal) {
	<-sig
	m.resetTerminal()
	os.Exit(0)
}

func (m *Monitor) resetTerminal() {
	fmt.Fprint(m.out, AnsiCursorOn)
	fmt.Fprintf(m.out, AnsiReset)
}

// https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
func terminalWidth() (int, error) {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	ws := &winsize{}
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		return 0, errors.New("unable to get terminal size")
	}
	return int(ws.Col), nil
}

func formatElapsed(elapsed time.Duration) string {
	total := int(elapsed.Seconds())
	seconds := total % 60
	minutes := total % 3600 / 60
	hours := total / 3600
	return fmt.Sprintf("%3d:%02d:%02d", hours, minutes, seconds)
}
