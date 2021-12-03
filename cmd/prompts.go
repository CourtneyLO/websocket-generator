package cmd

import (
	"fmt"
	"os"
	"github.com/manifoldco/promptui"
)

func language() string {
	fmt.Printf("Language")
	prompt := promptui.Select{
		Label: "Select your preferred language",
		Items: []string{"Node", "Typescript", "Python", "Golang", "Java", "Ruby"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Language prompt failed %v\n, Node will be used as a default", err)
		return "Node"
	}

	switch result {
	case "Node":
		return result
	case "Typescript":
		fmt.Println(result + " has not yet be included. Node will be used as a default.")
		return "Node"
	case "Python":
		fmt.Println(result + " has not yet be included. Node will be used as a default.")
		return "Node"
	case "Golang":
		fmt.Println(result + " has not yet be included. Node will be used as a default.")
		return "Node"
	case "Java":
		fmt.Println(result + " has not yet be included. Node will be used as a default.")
		return "Node"
	}

	// Default
	return "Node"
}

func prompt(label string, defaultValue string) string {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()

	if result == "quit" {
		fmt.Println("Exiting WebSocket Generator Int, Bye!")
		os.Exit(-1)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	if (result == "" && defaultValue != "") {
		defaultMessage := "You have made no entry so the default value will be used"
		fmt.Println(defaultMessage)
		return defaultValue
	}

	return result
}

func requiredPrompt(label string) string {
	result := prompt(label, "")

	if result == "" {
		fmt.Println("This value is required. Please enter it now")
		result = requiredPrompt(label)
	}

	return result
}

func yesNo(label string) (bool, string) {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Yes No prompt failed %v\n. False is returned as the default", err)
		return false, ""
	}

	if result == "Yes" {
		return true, result
	} else {
		return false, result
	}
}
