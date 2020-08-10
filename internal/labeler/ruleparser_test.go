package labeler

import (
	"testing"

	"github.com/tanmancan/label-it/v1/internal/config"
	"github.com/tanmancan/label-it/v1/internal/gitapi"
)

func TestRuleTypeStringValidator(t *testing.T) {
	type args struct {
		r config.RuleTypeString
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Passes all checks",
			args{
				config.RuleTypeString{
					Exact:   "octopus hello",
					NoExact: "world",
					Match:   "^(octopus)",
					NoMatch: "(lion)$",
				},
				"octopus hello",
			},
			true,
		},
		{
			"passes exact check",
			args{
				config.RuleTypeString{
					Exact: "LionOctopus",
				},
				"LionOctopus",
			},
			true,
		},
		{
			"does not pass exact check",
			args{
				config.RuleTypeString{
					Exact: "TigerHippo",
				},
				"HippoTiger",
			},
			false,
		},
		{
			"passes no-exact check",
			args{
				config.RuleTypeString{
					NoExact: "FishHead",
				},
				"BirdWatch",
			},
			true,
		},
		{
			"does not pass no-exact check",
			args{
				config.RuleTypeString{
					NoExact: "TreeSnake",
				},
				"TreeSnake",
			},
			false,
		},
		{
			"passes match check",
			args{
				config.RuleTypeString{
					Match: "\\w*ee",
				},
				"One two three four five",
			},
			true,
		},
		{
			"does not pass no-match check",
			args{
				config.RuleTypeString{
					NoMatch: "^(Tr)",
				},
				"TreeSnake",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RuleTypeStringValidator(tt.args.r, tt.args.s); got != tt.want {
				t.Errorf("RuleTypeStringValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleTypeIntValidator(t *testing.T) {
	type args struct {
		r config.RuleTypeInt
		i int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Passes all checks",
			args{
				config.RuleTypeInt{
					Exact:   5790,
					NoExact: 4,
					Match:   "^(57)",
					NoMatch: "(91)$",
				},
				5790,
			},
			true,
		},
		{
			"passes exact check",
			args{
				config.RuleTypeInt{
					Exact: 23,
				},
				23,
			},
			true,
		},
		{
			"does not pass exact check",
			args{
				config.RuleTypeInt{
					Exact: 23,
				},
				44,
			},
			false,
		},
		{
			"passes no-exact check",
			args{
				config.RuleTypeInt{
					NoExact: 44,
				},
				32,
			},
			true,
		},
		{
			"does not pass no-exact check",
			args{
				config.RuleTypeInt{
					NoExact: 44,
				},
				44,
			},
			false,
		},
		{
			"passes match check",
			args{
				config.RuleTypeInt{
					Match: "^(35)",
				},
				3566,
			},
			true,
		},
		{
			"does not pass no-match check",
			args{
				config.RuleTypeInt{
					NoMatch: "(56)$",
				},
				3456,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RuleTypeIntValidator(tt.args.r, tt.args.i); got != tt.want {
				t.Errorf("RuleTypeIntValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRule_MatchAllRules(t *testing.T) {
	type fields struct {
		Label       string
		HeadRules   config.RuleTypeString
		BaseRules   config.RuleTypeString
		TitleRules  config.RuleTypeString
		BodyRules   config.RuleTypeString
		UserRules   config.RuleTypeString
		NumberRules config.RuleTypeInt
	}
	type args struct {
		pr gitapi.PullRequest
	}
	tests := []struct {
		name   string
		fields Rule
		args   args
		want   bool
	}{
		{
			"Match All Rules",
			Rule{
				Label: "My Test Rule Label",
				HeadRules: config.RuleTypeString{
					Exact:   "octopus-hello-head",
					NoExact: "world",
					Match:   "^(octopus)",
					NoMatch: "(lion)$",
				},
				BaseRules: config.RuleTypeString{
					Exact:   "octopus-hello-base",
					NoExact: "world",
					Match:   "^(octopus)",
					NoMatch: "(lion)$",
				},
				TitleRules: config.RuleTypeString{
					Exact:   "Test PR Title",
					NoExact: "world",
					Match:   "^(Tes)",
					NoMatch: "(PR)$",
				},
				BodyRules: config.RuleTypeString{
					Exact:   "Test PR Body Text",
					NoExact: "world",
					Match:   "^(Tes)",
					NoMatch: "(Body)$",
				},
				UserRules: config.RuleTypeString{
					Exact:   "tanmancan",
					NoExact: "world",
					Match:   "^(tan)",
					NoMatch: "(man)$",
				},
				NumberRules: config.RuleTypeInt{
					Exact:   5790,
					NoExact: 4,
					Match:   "^(57)",
					NoMatch: "(91)$",
				},
			},
			args{
				pr: gitapi.PullRequest{
					Number: 5790,
					Head: gitapi.PrBranch{
						Ref: "octopus-hello-head",
					},
					Base: gitapi.PrBranch{
						Ref: "octopus-hello-base",
					},
					Title: "Test PR Title",
					Body:  "Test PR Body Text",
					User: gitapi.PrUser{
						Login: "tanmancan",
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rule{
				Label:       tt.fields.Label,
				HeadRules:   tt.fields.HeadRules,
				BaseRules:   tt.fields.BaseRules,
				TitleRules:  tt.fields.TitleRules,
				BodyRules:   tt.fields.BodyRules,
				UserRules:   tt.fields.UserRules,
				NumberRules: tt.fields.NumberRules,
			}
			if got := r.MatchAllRules(tt.args.pr); got != tt.want {
				t.Errorf("Rule.MatchAllRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
