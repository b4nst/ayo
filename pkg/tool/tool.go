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
	return exec.Command(t.Cmd, args...).CombinedOutput()
}

// LoadTool loads a tool from a json file
func LoadTool(file string) (*Tool, error) {
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

// LoadTools loads a list of tools from a directory
func LoadTools(dir string) ([]*Tool, error) {
	// Get the list of json files in the directory
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Load the tools concurrently, pool of 10 goroutines max
	var tools []*Tool
	lock := sync.Mutex{}
	eg := new(errgroup.Group)
	eg.SetLimit(10)
	for _, f := range files {
		file := f // capture the loop variable
		eg.Go(func() error {
			tool, err := LoadTool(file)
			if err != nil {
				return err
			}

			lock.Lock()
			defer lock.Unlock()
			tools = append(tools, tool)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return tools, nil
}
