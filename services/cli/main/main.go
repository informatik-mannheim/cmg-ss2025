package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/informatik-mannheim/cmg-ss2025/services/cli"
	"github.com/informatik-mannheim/cmg-ss2025/services/cli/client"
	"log"
	"os"
	"strings"
)

func getValue(args []string, arg string) string {
	// search args[]: does it contain arg?
	for i, v := range args {
		// do this if we find the arg:
		if v == arg {
			// check if there is anything in the array AFTER arg
			if i+1 < len(args) {
				// the thing after arg is the value we are looking for
				value := args[i+1]
				// unless it starts with "--", in which case the user did not provide a value, but there is another flag instead
				if !strings.HasPrefix(value, "--") {
					return value
				}
			}
		}
	}
	// return NO_VALUE in case no value was provided or found
	return "NO_VALUE"
}

func parseParameters(paramsCsv string) map[string]string {
	result := make(map[string]string)

	pairs := strings.Split(paramsCsv, ",")
	for _, pair := range pairs {
		// überspringe leere Teile (z.B. durch Doppelkomma)
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			// Fehlerhafte Eingabe → gib nil zurück
			return nil
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}

	return result
}

func printGeneralHelp() {
	fmt.Println("COMMAND DESCRIPTION")
	for _, command := range allCommands {
		fmt.Println(command.Name, " - ", command.Description)
	}
	fmt.Println("For help on a specific command, simply enter the command without any arguments.")
}

type Command struct {
	Name        string
	Description string
	Parameters  map[string]bool
	ParamOrder  []string
	Execute     func(args []string) error
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func (c *Command) allRequiredArgumentsProvided(providedArgs []string) bool {
	allArgumentsProvided := true
	for param, required := range c.Parameters {
		if required {
			// handle missing required argument
			if !contains(providedArgs, param) {
				fmt.Println("Missing required argument:", param)
				allArgumentsProvided = false
				continue
			}
			// handle missing value for argument
			if getValue(providedArgs, param) == "NO_VALUE" {
				fmt.Println("You need to provide a value for ", param)
				allArgumentsProvided = false
				continue
			}
		}
	}
	return allArgumentsProvided
}

var allCommands []Command

func (c *Command) isMissingArguments(args []string) bool {
	if !c.allRequiredArgumentsProvided(args) {
		printHelp(c.Name)
		return true
	}
	return false
}

func registerCommands(gatewayClient *client.GatewayClient) []Command {

	// Help command –––––––––––––––––––––––––––––––––––––––––
	helpCommand := Command{
		Name:        "help",
		Description: "Show this help",
		Parameters:  map[string]bool{},
	}
	helpCommand.Execute = func(args []string) error {
		printGeneralHelp()
		return nil
	}
	allCommands = append(allCommands, helpCommand)

	// Create Job Command –––––––––––––––––––––––––––––––––––––––––
	createJobCommand := Command{
		Name:        "create-job",
		Description: "Create a new job",
		Parameters: map[string]bool{
			"--image-name":    true,
			"--image-version": true,
			"--job-name":      true,
			"--creation-zone": false,
			"--parameters":    true,
		},
		ParamOrder: []string{"--job-name", "--creation-zone", "--image-name", "--image-version", "--parameters"},
	}
	createJobCommand.Execute = func(args []string) error {
		// handle image_id
		if createJobCommand.isMissingArguments(args) {
			return nil
		}
		imageName := getValue(args, "--image-name")
		imageVersion := getValue(args, "--image-version")
		containerImage := cli.ContainerImage{
			Name:    imageName,
			Version: imageVersion,
		}
		jobName := getValue(args, "--job-name")
		creationZone := getValue(args, "--creation-zone")
		if creationZone == "NO_VALUE" {
			creationZone = ""
		}
		parametersValue := getValue(args, "--parameters")
		parameters := parseParameters(parametersValue)
		if parameters == nil {
			return errors.New("there was an error parsing the parameters")
		}

		// create the job once all checks have passed
		gatewayClient.CreateJob(jobName, creationZone, containerImage, parameters)
		return nil
	}
	allCommands = append(allCommands, createJobCommand)

	// Get job by its ID Command –––––––––––––––––––––––––––––––––––––––––
	getJobByIdCommand := Command{
		Name:        "get-job",
		Description: "Get job by its id",
		Parameters: map[string]bool{
			"--id": true,
		},
		ParamOrder: []string{"--id"},
	}
	getJobByIdCommand.Execute = func(args []string) error {
		if getJobByIdCommand.isMissingArguments(args) {
			return nil
		}
		Id := getValue(args, "--id")
		gatewayClient.GetJobById(Id)
		fmt.Printf("Getting job by id %s\n", getValue(args, "--id"))
		return nil
	}
	allCommands = append(allCommands, getJobByIdCommand)

	// Get the outcome of job Command –––––––––––––––––––––––––––––––––––––––––
	getJobOutcomeCommand := Command{
		Name:        "get-job-outcome",
		Description: "Get job outcome",
		Parameters: map[string]bool{
			"--id": true,
		},
		ParamOrder: []string{"--id"},
	}
	getJobOutcomeCommand.Execute = func(args []string) error {
		if getJobOutcomeCommand.isMissingArguments(args) {
			return nil
		}
		Id := getValue(args, "--id")
		gatewayClient.GetJobOutcome(Id)
		fmt.Printf("Getting job outcome for job  %s\n", getValue(args, "--id"))
		return nil
	}
	allCommands = append(allCommands, getJobOutcomeCommand)

	loginCommand := Command{
		Name:        "login",
		Description: "Log in by providing a secret",
		Parameters: map[string]bool{
			"--secret": true,
		},
		ParamOrder: []string{"--secret"},
	}
	loginCommand.Execute = func(args []string) error {
		if loginCommand.isMissingArguments(args) {
			return nil
		}
		secret := getValue(args, "--secret")
		gatewayClient.Login(secret)
		return nil
	}

	allCommands = append(allCommands, loginCommand)
	return allCommands
}

func printHelp(arg string) {
	if arg == "help" {
		printGeneralHelp()
		return
	}
	for _, command := range allCommands {
		if command.Name == arg {
			usage := "Usage: \n\t" + arg
			for _, param := range command.ParamOrder {
				required := command.Parameters[param]
				usage += " " + param + " <value>"
				if required {
					usage += " [required]"
				}
			}
			fmt.Println(usage)
			return
		}
	}
	fmt.Println("Unknown command:", arg)

}

func buildCLI(gatewayClient *client.GatewayClient) {
	commands := registerCommands(gatewayClient)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) > 0 {
			if len(args) == 1 {
				printHelp(args[0])
				fmt.Print("> ")
				continue
			}
			// check if the provided command exists
			for _, command := range commands {
				if command.Name == args[0] {
					// execute the function associated with the provided command
					err := command.Execute(args[1:])
					if err != nil {
						return
					}
				}
			}
		}
		fmt.Print("> ")
	}
}

func main() {
	gatewayUrl := os.Getenv("GATEWAY_URL")
	if gatewayUrl == "" {
		log.Fatal("Fehlende Umgebungsvariable GATEWAY_URL")
	}
	gateway := client.NewGatewayClient(gatewayUrl)

	buildCLI(gateway)
}
