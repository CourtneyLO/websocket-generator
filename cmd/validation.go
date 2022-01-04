package cmd

import (
	"fmt"
	"reflect"
	"strings"
)

func containsElement(givenInterface interface{}, element string) bool {
	value := reflect.ValueOf(givenInterface)

	if value.Kind() != reflect.Slice {
		fmt.Println("Invalid data-type")
		return false
	}

	for i := 0; i < value.Len(); i++ {
		if value.Index(i).Interface() == element {
			return true
		}
	}

	return false
}

func checkForValidArguments(action string, args []string, configFile map[string]interface{}) (bool, string, []string) {
	argumentsValid := true
	argumentInvalidMessage := ""
	argumentsForMessage := []string{}

	if len(args) != 2 {
		argumentsValid = false
		argumentInvalidMessage = PROJECT_ENVIRONMENT_MISSING_MESSAGE
		argumentsForMessage := append(argumentsForMessage, action)
		return argumentsValid, argumentInvalidMessage, argumentsForMessage
	}

	projectNameInArgs := strings.ToLower(args[0])
	environmentInArgs := strings.ToLower(args[1])
	projectNameInConfig := strings.ToLower(fmt.Sprintf("%v", configFile["projectName"]))
	environmentsInConfig := configFile["environments"]

	acceptedEnvironment := containsElement(environmentsInConfig, environmentInArgs)

	if !acceptedEnvironment {
		argumentsValid = false
		argumentInvalidMessage = UNACCEPTED_ENVIRONMENT_MESSAGE
		argumentsForMessage := append(argumentsForMessage, "init", projectNameInConfig)
		return argumentsValid, argumentInvalidMessage, argumentsForMessage
	}

	if projectNameInArgs != projectNameInConfig {
		argumentsValid = false
		argumentInvalidMessage = UNACCEPTED_PROJECT_NAME_MESSAGE
		argumentsForMessage := append(argumentsForMessage, action, projectNameInConfig, environmentInArgs)
		return argumentsValid, argumentInvalidMessage, argumentsForMessage
	}

	return argumentsValid, argumentInvalidMessage, argumentsForMessage
}
