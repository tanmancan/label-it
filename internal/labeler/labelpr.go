package labeler

import (
	"fmt"

	"github.com/tanmancan/label-it/v1/internal/gitapi"
)

// LabelPr adds labels to a given list of pull requests via the Github API
func LabelPr(prLabels []gitapi.PrLabel) {
	updateCount := len(prLabels)

	c := make(chan string, updateCount)
	for _, prLabel := range prLabels {
		go gitapi.AddLabels(prLabel, c)
	}

	for i := 0; i < updateCount; i++ {
		fmt.Println(<-c)
	}
}
