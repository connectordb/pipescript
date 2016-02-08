package core

import (
	"testing"

	"github.com/connectordb/pipescript"
)

func TestMap(t *testing.T) {
	Register()
	pipescript.TestCase{
		Pipescript: "map $ {next}",
		ParseError: true,
	}.Run(t)

	pipescript.TestCase{
		// This tests order of prescedence: ":" pipes are high prescedence, and will be executed first
		Pipescript: "map $ > 5 {count}",
		Input: []pipescript.Datapoint{
			{1, 4},
			{2, 6},
			{3, 8},
		},
		Output: []pipescript.Datapoint{
			{3, map[string]interface{}{"false": int64(1), "true": int64(2)}},
		},
		SecondaryInput: []pipescript.Datapoint{
			{4, 8},
			{5, 4},
			{6, 6},
		},
		SecondaryOutput: []pipescript.Datapoint{
			{6, map[string]interface{}{"false": int64(1), "true": int64(2)}},
		},
	}.Run(t)

	pipescript.TestCase{
		// This tests order of prescedence: ":" pipes are high prescedence, and will be executed first
		Pipescript: "map($ > 5,count)",
		Input: []pipescript.Datapoint{
			{1, 4},
			{2, 6},
			{3, 8},
		},
		Output: []pipescript.Datapoint{
			{3, map[string]interface{}{"false": int64(1), "true": int64(2)}},
		},
		SecondaryInput: []pipescript.Datapoint{
			{4, 8},
			{5, 4},
			{6, 6},
		},
		SecondaryOutput: []pipescript.Datapoint{
			{6, map[string]interface{}{"false": int64(1), "true": int64(2)}},
		},
	}.Run(t)
}