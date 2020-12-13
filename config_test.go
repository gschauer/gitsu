package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/stretchr/testify/require"
)

func Test_writeConfig(t *testing.T) {
	dir := os.TempDir()
	NoError(t, os.MkdirAll(dir, 0700))
	x, _ := ioutil.TempFile(dir, "")
	defer os.Remove(dir)
	configPath = x.Name()

	c := &config{Users: []user{
		{
			Name:  "John DOE",
			Email: "john.doe@acme.org",
		},
	}}

	NoError(t, writeConfig(c))
	data, err := ioutil.ReadFile(configPath)
	NoError(t, err)

	actual := &config{}
	NoError(t, json.Unmarshal(data, actual))
	Equal(t, c, actual)

	// remove the user and check that it is empty
	c.Users = []user{}
	NoError(t, writeConfig(c))

	data, err = ioutil.ReadFile(configPath)
	NoError(t, err)
	NoError(t, json.Unmarshal(data, &actual))
	Equal(t, c, actual)
}

func Test_readConfig(t *testing.T) {
	dir := os.TempDir()
	NoError(t, os.MkdirAll(dir, 0700))
	x, _ := ioutil.TempFile(dir, "")
	defer os.Remove(dir)
	configPath = x.Name()

	expected := &config{Users: []user{
		{
			Name:  "Jane DOE",
			Email: "jane.doe@acme.org",
		},
	}}

	data, err := json.Marshal(expected)
	NoError(t, err)
	NoError(t, ioutil.WriteFile(configPath, data, 0600))
	c, err := readConfig()
	NoError(t, err)
	Equal(t, expected, c)
}
