package bit

import (
	"errors"
	"path/filepath"
	"strings"
)

type Plan struct {
	Tasks  []string
	Errors []string
}

func LoadPlan(path string) *Plan {
	p := &Plan{}
	return p.addPath(path)
}

func (p *Plan) addPath(path string) *Plan {
	mode, err := fs.Mode(path)
	if err != nil {
		p.addError(path, err)
		return p
	}
	if mode.IsDir() {
		return p.addDir(path)
	}
	return p.addFile(path)
}

func (p *Plan) addDir(dir string) *Plan {
	taskFile := filepath.Join(dir, ".tasks")
	taskText, err := fs.LoadText(taskFile)
	if err != nil {
		p.addError(taskFile, err)
		return p
	}
	tasks := strings.Split(taskText, "\n")
	for _, task := range tasks {
		task := strings.TrimSpace(task)
		if task == "" {
			continue
		}
		if task[0] == '#' {
			continue
		}
		p.addPath(filepath.Join(dir, task))
	}
	return p
}

func (p *Plan) addFile(file string) *Plan {
	mode, err := fs.Mode(file)
	if err != nil {
		p.addError(file, err)
		return p
	}
	if mode&0100 == 0 {
		p.addError(file, errors.New("not executable"))
		return p
	}
	p.Tasks = append(p.Tasks, file)
	return p
}

func (p *Plan) addError(path string, err error) {
	p.Errors = append(p.Errors, path+": "+err.Error())
}
