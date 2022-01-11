package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
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

	prompt := promptui.Prompt{
		Label: questionDetails.QuestionLabel,
	}

	answer := questionDetails.DefaultResponse

	result, error := prompt.Run()

	if error != nil {
		fmt.Printf("Prompt failed %v\n", error)
		return ""
	}

	if result != "" {
		answer = strings.ReplaceAll(result, " ", "")
	} else {
		defaultMessage := "You have made no entry so the default value has been used"
		fmt.Println(defaultEmoji + defaultMessage)
	}

	if answer == "quit" || answer == "q" {
		goodbyeMessage.Println("Exiting WebSocket Generator, no configuration file has been created. Bye " + goodbyeEmoji)
		os.Exit(-1)
	}

	questionLabel.Println(questionDetails.ResponseLabel)

	if reference == "Environment" {
		answer = strings.ToLower(answer)
	}

	if (reference == "InfrastructureFilePath" || reference == "WebsocketFilePath") {
		currentDirectory, error := os.Getwd()
		if error != nil {
			fmt.Println(error)
		}

		destinationFilePath := filepath.Join(currentDirectory, answer)
		answerColour.Println(resultEmoji + destinationFilePath)

		if !strings.HasPrefix(answer, "/") {
			answer = "/" + answer
		}

		if strings.HasSuffix(answer, "/") {
			answer = strings.TrimRight(answer, "/")
		}
	} else {
		answerColour.Println(resultEmoji + answer)
	}


	fmt.Println(boldMessage.Sprintf("Answer:"), answer)
	return answer
}

func requiredPrompt(label string) string {
	result := prompt("AWSAccountID")

	if result == "" {
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

	fmt.Println("You chose " + selectionColour(result) + " to adding another environment")

	if result == "Yes" {
		return true
	} else {
		return false
	}
}
