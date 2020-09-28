package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tanmancan/label-it/v1/internal/common"
	"github.com/tanmancan/label-it/v1/internal/config"
	"github.com/tanmancan/label-it/v1/internal/gitapi"
	"github.com/tanmancan/label-it/v1/internal/labeler"
)

// Display list of labels to be added, if found
func printLabelSummary(prLabels []gitapi.PrLabel) {
	updateCount := len(prLabels)
	fmt.Printf("Found %[1]d matching pull request.\n", updateCount)
	fmt.Println("PR\tLabels")
	fmt.Println("--\t------")
	for _, prLabel := range prLabels {
		fmt.Printf("%[1]d\t%[2]s\n", prLabel.Issue, strings.Join(prLabel.Labels, ", "))
	}
	fmt.Print("\n")
}

// Ask users for confirmation before applying labels
func userConfirm() bool {
	fmt.Println("Do you want to continue? (y/n)")

	if config.AutoConfirm == true {
		fmt.Println("y")
		return true
	}

	var userInput string

	_, err := fmt.Scanln(&userInput)
	common.CheckErr(err)

	switch strings.ToLower(userInput) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return userConfirm()
	}
}

func main() {
	err := config.SetupArgs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.LoadYaml()

	prList := gitapi.ListPulls()
	prLabels := labeler.RuleParser(prList)

	printLabelSummary(prLabels)

	if config.DryRun == true {
		fmt.Println("Perform dry run. Pull requests were not updated.")
		return
	}

	if len(prLabels) == 0 {
		return
	}

	confirm := userConfirm()

	if confirm == true {
		labeler.LabelPr(prLabels)
	}
}
