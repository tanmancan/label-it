package main

import (
	"github.com/tanmancan/label-it/v1/internal/config"
	"github.com/tanmancan/label-it/v1/internal/gitapi"
	"github.com/tanmancan/label-it/v1/internal/labeler"
)

func main() {
	config.SetupArgs()
	config.LoadYaml()

	prList := gitapi.ListPulls()
	prLabels := labeler.RuleParser(prList)
	labeler.LabelPr(prLabels)
}
