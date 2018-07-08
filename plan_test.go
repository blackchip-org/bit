package bit

import (
	"reflect"
	"strings"
	"testing"
)

func TestSingleFile(t *testing.T) {
	fs = NewMockFileSystem()
	fs.SaveText("task/single", "", 0755)
	p := LoadPlan("task/single")
	if len(p.Errors) > 0 {
		t.Errorf("unexpected errors: \n%v", strings.Join(p.Errors, "\n"))
	}
	want := []string{"task/single"}
	have := p.Tasks
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestDir(t *testing.T) {
	fs = NewMockFileSystem()
	taskList := "task1 \n task2 \n task3"
	fs.Mkdir("plan", 0755)
	fs.SaveText("plan/.tasks", taskList, 0644)
	fs.SaveText("plan/task1", "", 0755)
	fs.SaveText("plan/task2", "", 0755)
	fs.SaveText("plan/task3", "", 0755)
	p := LoadPlan("plan")

	if len(p.Errors) > 0 {
		t.Errorf("unexpected errors: \n%v", strings.Join(p.Errors, "\n"))
	}
	want := []string{"plan/task1", "plan/task2", "plan/task3"}
	have := p.Tasks
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestComment(t *testing.T) {
	fs = NewMockFileSystem()
	taskList := "task1 \n          #task2 \n task3"
	fs.Mkdir("plan", 0755)
	fs.SaveText("plan/.tasks", taskList, 0644)
	fs.SaveText("plan/task1", "", 0755)
	fs.SaveText("plan/task2", "", 0755)
	fs.SaveText("plan/task3", "", 0755)
	p := LoadPlan("plan")

	if len(p.Errors) > 0 {
		t.Errorf("unexpected errors: \n%v", strings.Join(p.Errors, "\n"))
	}
	want := []string{"plan/task1", "plan/task3"}
	have := p.Tasks
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestBlankLines(t *testing.T) {
	fs = NewMockFileSystem()
	taskList := "\n\n\n task1 \n      \n task2 \n task3"
	fs.Mkdir("plan", 0755)
	fs.SaveText("plan/.tasks", taskList, 0644)
	fs.SaveText("plan/task1", "", 0755)
	fs.SaveText("plan/task2", "", 0755)
	fs.SaveText("plan/task3", "", 0755)
	p := LoadPlan("plan")

	if len(p.Errors) > 0 {
		t.Errorf("unexpected errors: \n%v", strings.Join(p.Errors, "\n"))
	}
	want := []string{"plan/task1", "plan/task2", "plan/task3"}
	have := p.Tasks
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestMissingTask(t *testing.T) {
	fs = NewMockFileSystem()
	taskList := "task1 \n task2 \n task3"
	fs.Mkdir("plan", 0755)
	fs.SaveText("plan/.tasks", taskList, 0644)
	fs.SaveText("plan/task1", "", 0755)
	fs.SaveText("plan/task3", "", 0755)
	p := LoadPlan("plan")

	want := []string{"plan/task2: file does not exist"}
	have := p.Errors
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestMissingTaskList(t *testing.T) {
	fs = NewMockFileSystem()
	fs.Mkdir("plan", 0755)
	p := LoadPlan("plan")

	want := []string{"plan/.tasks: file does not exist"}
	have := p.Errors
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestTaskNotExecutable(t *testing.T) {
	fs = NewMockFileSystem()
	taskList := "task1 \n task2 \n task3"
	fs.Mkdir("plan", 0755)
	fs.SaveText("plan/.tasks", taskList, 0644)
	fs.SaveText("plan/task1", "", 0755)
	fs.SaveText("plan/task2", "", 0644)
	fs.SaveText("plan/task3", "", 0755)
	p := LoadPlan("plan")

	want := []string{"plan/task2: not executable"}
	have := p.Errors
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestMultipleMissingTasks(t *testing.T) {
	fs = NewMockFileSystem()
	taskList := "task1 \n task2 \n task3"
	fs.Mkdir("plan", 0755)
	fs.SaveText("plan/.tasks", taskList, 0644)
	p := LoadPlan("plan")

	want := []string{
		"plan/task1: file does not exist",
		"plan/task2: file does not exist",
		"plan/task3: file does not exist",
	}
	have := p.Errors
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestNestedDir(t *testing.T) {
	fs = NewMockFileSystem()
	fs.Mkdir("plan", 0755)
	fs.Mkdir("plan/alpha", 0755)
	fs.Mkdir("plan/alpha/bravo", 0755)
	fs.SaveText("plan/.tasks", "alpha \n bar", 0644)
	fs.SaveText("plan/alpha/.tasks", "bravo", 0644)
	fs.SaveText("plan/alpha/bravo/.tasks", "foo", 0644)
	fs.SaveText("plan/alpha/bravo/foo", "", 0755)
	fs.SaveText("plan/bar", "", 0755)
	p := LoadPlan("plan")

	if len(p.Errors) > 0 {
		t.Errorf("unexpected errors: \n%v", strings.Join(p.Errors, "\n"))
	}
	want := []string{"plan/alpha/bravo/foo", "plan/bar"}
	have := p.Tasks
	if !reflect.DeepEqual(want, have) {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}
