package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji/v2"
)

func language() string {
	templates := promptui.SelectTemplates{
		Selected: `{{ "✔" | green | bold }} {{ "Language Choosen" | bold }}: {{ . | blue | bold }}`,
	}

	prompt := promptui.Select{
		Label: "Select your preferred language",
		Items: []string{"Node", "Typescript", "Python", "Golang", "Java", "Ruby"},
		Templates: &templates,
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

type QuestionDetails struct {
	QuestionLabel   string
	DefaultResponse string
	ResponseLabel   string
}

func prompt(reference string) string {
	questionDetails := GetQuestionDetails(reference)
	questionLabel := color.New(color.FgBlack, color.Bold)
	answer := color.New(color.FgBlue, color.Bold)
	resultEmoji := emoji.Sprint(":pencil2: ")

	prompt := promptui.Prompt{
		Label: questionDetails.QuestionLabel,
	}

	result, error := prompt.Run()

	if result == "quit" || result == "q" {
		goodbyeEmoji := emoji.Sprint(":wave: ")
		goodbyeMessage := color.New(color.FgGreen, color.Bold)
		goodbyeMessage.Println("Exiting WebSocket Generator, Bye " + goodbyeEmoji)
		os.Exit(-1)
	}

	if error != nil {
		fmt.Printf("Prompt failed %v\n", error)
		return ""
	}

	questionLabel.Println(questionDetails.ResponseLabel)

	if (result == "" && questionDetails.DefaultResponse != "") {
		answer.Println(resultEmoji + questionDetails.DefaultResponse)
		defaultMessage := "You have made no entry so the default value has been used"
		defaultEmoji := emoji.Sprint(":information_desk_person: ")
		fmt.Println(defaultEmoji + defaultMessage)
		return questionDetails.DefaultResponse
	}

	if (reference == "InfrastructureFilePath" || reference == "WebsocketFilePath") {
		currentDirectory, error := os.Getwd()
		if error != nil {
			fmt.Println(error)
		}
		destinationFilePath := filepath.Join(currentDirectory, result)
		answer.Println(resultEmoji + destinationFilePath + result)
		return destinationFilePath + result
	}

	answer.Println(resultEmoji + result)
	return result
}

func requiredPrompt(label string) string {
	result := prompt("AWSAccountID")

	if result == "" {
		requiredEmoji := emoji.Sprint(":x: ")
		fmt.Println(requiredEmoji + "This value is required. Please enter it now")
		result = requiredPrompt(label)
	}

	return result
}

func yesNo(label string) bool {
	templates := promptui.SelectTemplates{
		Selected: `{{ "✔" | green | bold }} {{ "Answer" | bold }}: {{ . | blue | bold }}`,
	}

	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
		Templates: &templates,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Yes No prompt failed %v\n. False is returned as the default", err)
		return false
	}

	selectionColour := color.New(color.FgBlue, color.Bold).SprintFunc()
	fmt.Println("You chose " + selectionColour(result) + " to adding another environment")

	if result == "Yes" {
		return true
	} else {
		return false
	}
}
