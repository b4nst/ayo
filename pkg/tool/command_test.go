package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderArgs(t *testing.T) {
	t.Parallel()

	t.Run("no args", func(t *testing.T) {
		t.Parallel()
		args := []string{}
		values := map[string]any{"message": "Hello, World!"}
		expected := []string{}

		actual, err := RenderArgs(args, values)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("no values", func(t *testing.T) {
		t.Parallel()
		args := []string{"echo", "Hello, World!"}
		values := map[string]any{}
		expected := []string{"echo", "Hello, World!"}

		actual, err := RenderArgs(args, values)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("template error", func(t *testing.T) {
		t.Parallel()
		args := []string{"{{.message}"}
		values := map[string]any{"message": "Hello, World!"}

		_, err := RenderArgs(args, values)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "bad character")
		}
	})

	t.Run("missing value", func(t *testing.T) {
		t.Parallel()
		args := []string{"echo", "{{.message}}"}
		values := map[string]any{}

		_, err := RenderArgs(args, values)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "map has no entry for key \"message\"")
		}
	})

	t.Run("nominal", func(t *testing.T) {
		t.Parallel()
		args := []string{"echo", "{{.message}}"}
		values := map[string]any{"message": "Hello, World!"}
		expected := []string{"echo", "Hello, World!"}

		actual, err := RenderArgs(args, values)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
