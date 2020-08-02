package labeler

import (
	"fmt"
	"label-it/internal/config"
	"label-it/internal/gitapi"
	"strings"
)

func printPrLabel(prLabels []gitapi.PrLabel) {
	fmt.Println("Dry run - List of pull request #s and labels to be applied")
	fmt.Println("PR\tLabels")
	for _, prLabel := range prLabels {
		fmt.Printf("%[1]d\t%[2]s\n", prLabel.Issue, strings.Join(prLabel.Labels, ", "))
	}
}

// LabelPr adds labels to a given list of pull requests via the Github API
func LabelPr(prLabels []gitapi.PrLabel) {
	if config.DryRun == true {
		printPrLabel(prLabels)
	}
	if config.DryRun == false {
		for _, prLabel := range prLabels {
			gitapi.AddLabels(prLabel)
		}
	}
}
