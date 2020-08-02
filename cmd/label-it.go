package main

import (
	"label-it/internal/config"
	"label-it/internal/gitapi"
	"label-it/internal/labeler"
)

func main() {
	config.SetupArgs()
	config.LoadYaml()

	prList := gitapi.ListPulls()
	prLabels := labeler.RuleParser(prList)
	labeler.LabelPr(prLabels)
}
