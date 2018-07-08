package bit

import (
	"io/ioutil"
	"os"
)

type FileSystem interface {
	LoadText(path string) (string, error)
	SaveText(path string, text string, mode os.FileMode) error
	Mkdir(path string, mode os.FileMode) error
	Mode(path string) (os.FileMode, error)
	Getwd() (string, error)
	Chdir(path string) error
}

type actual struct{}

var fs FileSystem = actual{}

func (f actual) LoadText(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f actual) SaveText(path string, text string, mode os.FileMode) error {
	return ioutil.WriteFile(path, []byte(text), mode)
}

func (f actual) Mkdir(path string, mode os.FileMode) error {
	return os.MkdirAll(path, mode)
}

func (f actual) IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func (f actual) Mode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode(), nil
}

func (f actual) Getwd() (string, error) {
	return os.Getwd()
}

func (f actual) Chdir(path string) error {
	return os.Chdir(path)
}

type mockFile struct {
	text string
	mode os.FileMode
	err  error
}

type mock struct {
	cwd   string
	files map[string]mockFile
}

func NewMockFileSystem() FileSystem {
	return &mock{
		cwd:   "/",
		files: map[string]mockFile{},
	}
}

func (f mock) LoadText(path string) (string, error) {
	file, ok := f.files[path]
	if !ok {
		return "", os.ErrNotExist
	}
	return file.text, nil
}

func (f *mock) SaveText(path string, text string, mode os.FileMode) error {
	file := mockFile{text: text, mode: mode}
	f.files[path] = file
	return nil
}

func (f *mock) Mkdir(path string, mode os.FileMode) error {
	prev, ok := f.files[path]
	if ok && !prev.mode.IsDir() {
		return os.ErrExist
	}
	file := mockFile{mode: mode | os.ModeDir}
	f.files[path] = file
	return nil
}

func (f mock) IsDir(path string) (bool, error) {
	file, ok := f.files[path]
	if !ok {
		return false, os.ErrNotExist
	}
	return file.mode.IsDir(), nil
}

func (f mock) Mode(path string) (os.FileMode, error) {
	file, ok := f.files[path]
	if !ok {
		return 0, os.ErrNotExist
	}
	return file.mode, nil
}

func (f mock) Getwd() (string, error) {
	return f.cwd, nil
}

func (f *mock) Chdir(path string) error {
	f.cwd = path
	return nil
}
