package file

import (
	"encoding/json"
	"os"
)

// const dataPath = "/app/data/secret.json"

type File struct {
	path string
}

func New(path string) *File {
	return &File{
		path: path,
	}
}

func (f *File) Read() []byte {
	file, err := os.ReadFile(f.path)
	if err != nil {
		return nil // file doesn't exist yet, start fresh
	}
	return file
}

func (f *File) ReadAsJson() map[string]any {
	file, err := os.ReadFile(f.path)
	if err != nil {
		return nil // file doesn't exist yet, start fresh
	}
	var data map[string]any
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil
	}
	return data
}

func (f *File) Save(content []byte) error {
	return os.WriteFile(f.path, content, 0644)
}
