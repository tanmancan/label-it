package main

import (
	"fmt"
	"os"

	"github.com/tanmancan/label-it/v1/internal/config"
	"github.com/tanmancan/label-it/v1/internal/gitapi"
	"github.com/tanmancan/label-it/v1/internal/labeler"
)

func main() {
	err := config.SetupArgs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.LoadYaml()

	prList := gitapi.ListPulls()
	prLabels := labeler.RuleParser(prList)
	labeler.LabelPr(prLabels)
}
