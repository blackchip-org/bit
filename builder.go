package bit

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Builder struct {
	Plan    *Plan
	Monitor *Monitor
}

func NewBuilder(base string) (*Builder, error) {
	b := &Builder{
		Plan:    LoadPlan(base),
		Monitor: NewMonitor(),
	}
	if len(b.Plan.Errors) > 0 {
		errorText := "unable to execute plan: " + base +
			"\n\nThe plan has the following errors:\n - " +
			strings.Join(b.Plan.Errors, "\n - ")
		return nil, errors.New(errorText)
	}
	return b, nil
}

func (b *Builder) Run() error {
	go b.Monitor.Run()
	for _, task := range b.Plan.Tasks {
		if err := b.runTask(task); err != nil {
			b.Monitor.Error()
			return err
		}
	}
	b.Monitor.Stop()
	return nil
}

func (b *Builder) runTask(task string) error {
	b.Monitor.SetTaskName(task)
	wd, err := fs.Getwd()
	if err != nil {
		return err
	}
	if err := fs.Chdir(filepath.Dir(task)); err != nil {
		return err
	}

	var cmdexec string
	var cmdargs []string
	_, err = os.Stat(".exec")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		cmdexec = "./.exec"
		cmdargs = []string{task}
	} else {
		cmdexec = "./" + filepath.Base(task)
		cmdargs = []string{}
	}

	cmd := exec.Command(cmdexec, cmdargs...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	b.attachOuputHandler(stdout)
	b.attachOuputHandler(stderr)

	err = cmd.Run()
	if err != nil {
		return err
	}

	if err := os.Chdir(wd); err != nil {
		return err
	}

	return nil
}

func (b *Builder) attachOuputHandler(output io.Reader) {
	go func() {
		r := bufio.NewReader(output)
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			b.Monitor.Update(line)
		}
	}()
}
