package labeler

import (
	"fmt"
	"strings"

	"github.com/tanmancan/label-it/v1/internal/config"
	"github.com/tanmancan/label-it/v1/internal/gitapi"
)

func printPrLabel(prLabels []gitapi.PrLabel) {
	fmt.Println("PR\tLabels")
	fmt.Println("--\t------")
	for _, prLabel := range prLabels {
		fmt.Printf("%[1]d\t%[2]s\n", prLabel.Issue, strings.Join(prLabel.Labels, ", "))
	}
	fmt.Print("\n")
}

// LabelPr adds labels to a given list of pull requests via the Github API
func LabelPr(prLabels []gitapi.PrLabel) {
	updateCount := len(prLabels)
	fmt.Printf("Found %[1]d matching pull request.\n", updateCount)
	printPrLabel(prLabels)

	if config.DryRun == true {
		fmt.Println("Perform dry run. Pull request were not updated.")
	}

	if config.DryRun == false {
		c := make(chan string, updateCount)
		for _, prLabel := range prLabels {
			go gitapi.AddLabels(prLabel, c)
		}

		for i := 0; i < updateCount; i++ {
			fmt.Println(<-c)
		}
	}
}
