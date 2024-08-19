package tool

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/ollama/ollama/api"
	"golang.org/x/sync/errgroup"
)

const Ext = ".json"

type Tool struct {
	api.Tool
	Cmd         string   `json:"cmd"`
	Args        []string `json:"args"`
	ConfirmExec bool     `json:"confirm_exec"`
}

func (t *Tool) Run(values map[string]any) ([]byte, error) {
	args, err := RenderArgs(t.Args, values)
	if err != nil {
		return nil, err
	}
	return exec.Command(t.Cmd, args...).CombinedOutput() //nolint:gosec // This is a tool runner.
}

// Load loads a tool from a json file.
func Load(file string) (*Tool, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tool Tool
	err = json.NewDecoder(f).Decode(&tool)
	if err != nil {
		return nil, err
	}

	return &tool, nil
}

// LoadAll loads a list of tools from a directory.
func LoadAll(dir string) (map[string]*Tool, error) {
	tools := make(map[string]*Tool)

	// Get the list of json files in the directory
	files, err := ListFiles(dir)
	if err != nil {
		return tools, err
	}

	// Load the tools concurrently, pool of 10 goroutines max
	lock := sync.Mutex{}
	eg := new(errgroup.Group)
	eg.SetLimit(10)
	for _, f := range files {
		file := f // capture the loop variable
		eg.Go(func() error {
			tool, err := Load(file)
			if err != nil {
				return err
			}

			lock.Lock()
			defer lock.Unlock()
			tools[tool.Function.Name] = tool
			return nil
		})
	}

	return tools, eg.Wait()
}

// ListToolFiles lists the tool files in a directory.
func ListFiles(dir string) ([]string, error) {
	// Get the list of json files in the directory
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == Ext {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
