package main

import (
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/cli/client"
	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	args := []string{"--secret", "mysecret"}
	assert.Equal(t, "mysecret", getValue(args, "--secret"))
	assert.Equal(t, "NO_VALUE", getValue(args, "--missing"))
}

func TestParseParametersValid(t *testing.T) {
	params := "a=1,b=2"
	parsed := parseParameters(params)
	assert.Equal(t, "1", parsed["a"])
	assert.Equal(t, "2", parsed["b"])
}

func TestCreateJobCommand_MissingImageName(t *testing.T) {
	cmds := registerCommands(&client.GatewayClient{})
	var createJobCmd *Command
	for _, cmd := range cmds {
		if cmd.Name == "create-job" {
			createJobCmd = &cmd
			break
		}
	}
	assert.NotNil(t, createJobCmd)

	args := []string{
		"--job-name", "TestJob",
		"--image-version", "1.0",
		"--parameters", "a=1",
	}
	err := createJobCmd.Execute(args)
	assert.NoError(t, err)
}

func TestParseParametersInvalid(t *testing.T) {
	params := "a=1,b,c=3"
	parsed := parseParameters(params)
	assert.Nil(t, parsed)
}

func TestAllRequiredArgumentsProvided(t *testing.T) {
	cmd := Command{
		Parameters: map[string]bool{
			"--a": true,
			"--b": false,
		},
	}
	args := []string{"--a", "1"}
	assert.True(t, cmd.allRequiredArgumentsProvided(args))

	args = []string{"--b", "2"}
	assert.False(t, cmd.allRequiredArgumentsProvided(args))
}

func TestPrintHelpForKnownAndUnknown(t *testing.T) {
	printHelp("create-job")
	printHelp("help")
	printHelp("invalid")
}

func TestIsMissingArgumentsTriggersHelp(t *testing.T) {
	cmd := Command{
		Name: "test-cmd",
		Parameters: map[string]bool{
			"--x": true,
		},
	}
	args := []string{}
	assert.True(t, cmd.isMissingArguments(args))
}

func TestHelpCommandExecutes(t *testing.T) {
	cmds := registerCommands(&client.GatewayClient{})
	var helpCmd *Command
	for _, cmd := range cmds {
		if cmd.Name == "help" {
			helpCmd = &cmd
			break
		}
	}
	assert.NotNil(t, helpCmd)
	err := helpCmd.Execute([]string{})
	assert.NoError(t, err)
}

func TestLoginCommandInputValidation(t *testing.T) {
	cmds := registerCommands(&client.GatewayClient{})
	var loginCmd *Command
	for _, cmd := range cmds {
		if cmd.Name == "login" {
			loginCmd = &cmd
			break
		}
	}
	assert.NotNil(t, loginCmd)
	err := loginCmd.Execute([]string{})
	assert.NoError(t, err)
}

// Table Driven Tests -------------------------------------------------------------------
func TestCreateJobCommand_TableDrivenInvalidInputs(t *testing.T) {
	cmds := registerCommands(&client.GatewayClient{})
	var cmd *Command
	for _, c := range cmds {
		if c.Name == "create-job" {
			cmd = &c
		}
	}
	assert.NotNil(t, cmd)

	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "missing --image-version",
			args:      []string{"--job-name", "J1", "--image-name", "img", "--parameters", "a=1"},
			wantError: false,
		},
		{
			name:      "missing --parameters",
			args:      []string{"--job-name", "J1", "--image-name", "img", "--image-version", "1.0"},
			wantError: false,
		},
		{
			name:      "invalid parameters format",
			args:      []string{"--job-name", "J1", "--image-name", "img", "--image-version", "1.0", "--parameters", "a=1,b,c=3"},
			wantError: true,
		},
		{
			name:      "missing --job-name",
			args:      []string{"--image-name", "img", "--image-version", "1.0", "--parameters", "a=1"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Execute(tt.args)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
