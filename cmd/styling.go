package cmd

import (
	"github.com/fatih/color"
	"github.com/kyokomi/emoji/v2"
)

var boldMessage = color.New(color.Bold)
var errorMessage = color.New(color.FgRed, color.Bold)
var goodbyeMessage = color.New(color.FgGreen, color.Bold)
var informationHeading = color.New(color.FgRed, color.Bold).Add(color.Underline)
var informationDetails = color.New(color.FgBlue, color.Bold)
var questionLabel = color.New(color.FgBlack, color.Bold)
var answerColour = color.New(color.FgBlue, color.Bold)
var selectionColour = color.New(color.FgBlue, color.Bold).SprintFunc()

var goodbyeEmoji = emoji.Sprint(":wave: ")
var resultEmoji = emoji.Sprint(":pencil2: ")
var requiredEmoji = emoji.Sprint(":x: ")
var defaultEmoji = emoji.Sprint(":information_desk_person: ")
