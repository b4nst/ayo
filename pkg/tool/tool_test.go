package tool

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ollama/ollama/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadTool(t *testing.T) {
	t.Parallel()

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		_, err := LoadTool("missing.json")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "no such file or directory")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		t.Parallel()
		f, err := os.CreateTemp(t.TempDir(), "bad.json")
		require.NoError(t, err)
		defer f.Close()
		_, err = f.WriteString("bad json")
		require.NoError(t, err)

		_, err = LoadTool(f.Name())
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid character")
		}
	})

	t.Run("nominal", func(t *testing.T) {
		t.Parallel()
		tool := &Tool{
			Cmd:  "echo",
			Args: []string{"Hello, World!"},
			Tool: api.Tool{
				Type: "function",
				Function: api.ToolFunction{
					Name:        "echo",
					Description: "Echoes a message",
				},
			},
		}
		f, err := os.CreateTemp(t.TempDir(), "nominal.json")
		require.NoError(t, err)
		defer f.Close()
		err = json.NewEncoder(f).Encode(tool)
		require.NoError(t, err)

		actual, err := LoadTool(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, tool, actual)
	})
}

func TestLoadTools(t *testing.T) {
	t.Parallel()

	t.Run("missing directory", func(t *testing.T) {
		t.Parallel()

		_, err := LoadTools("missing")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "no such file or directory")
		}
	})

	t.Run("no tools", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()

		tools, err := LoadTools(dir)
		assert.NoError(t, err)
		assert.Empty(t, tools)
	})

	t.Run("bad json", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		f, err := os.CreateTemp(dir, "*_bad.json")
		require.NoError(t, err)
		defer f.Close()
		_, err = f.WriteString("bad json")
		require.NoError(t, err)

		_, err = LoadTools(dir)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid character")
		}
	})

	t.Run("nominal", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		tools := map[string]*Tool{
			"echo": {
				Cmd:  "echo",
				Args: []string{"Hello, World!"},
				Tool: api.Tool{
					Type: "function",
					Function: api.ToolFunction{
						Name:        "echo",
						Description: "Echoes a message",
					},
				},
			},
			"cat": {
				Cmd:  "cat",
				Args: []string{"file.txt"},
				Tool: api.Tool{
					Type: "function",
					Function: api.ToolFunction{
						Name:        "cat",
						Description: "Cats a file",
					},
				},
			},
		}
		for _, tool := range tools {
			f, err := os.CreateTemp(dir, "*_tool.json")
			require.NoError(t, err)
			defer f.Close()
			err = json.NewEncoder(f).Encode(tool)
			require.NoError(t, err)
		}

		actual, err := LoadTools(dir)
		assert.NoError(t, err)
		assert.Equal(t, tools, actual)
	})
}
