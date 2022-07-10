package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonConfigReader_Read(t *testing.T) {
	reader := JsonConfigReader{}
	config, err := reader.Read("./git-notes.json.example")
	assert.NoError(t, err)

	assert.Equal(t, []string{"/Users/tanin/projects/personal-notes", "/Users/tanin/projects/another-personal-notes"}, config.Repos)
}
