package core

import (
	"testing"

	"github.com/heedy/pipescript"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	require.NoError(t, Filter.Register())
	pipescript.TestCase{
		Pipescript: "filter",
		Parsed:     "error",
	}.Run(t)

	pipescript.TestCase{
		Pipescript: "filter($)",
		Input: []pipescript.Datapoint{
			{Timestamp: 1, Data: 1},
			{Timestamp: 2, Data: true},
			{Timestamp: 3, Data: false},
			{Timestamp: 4, Data: "hi"},
		},
		Output: []pipescript.Datapoint{
			{Timestamp: 1, Data: 1},
			{Timestamp: 2, Data: true},
		},

		OutputError: true,
	}.Run(t)

	pipescript.TestCase{
		Pipescript: "filter $",
		Input: []pipescript.Datapoint{
			{Timestamp: 1, Data: 1},
			{Timestamp: 2, Data: true},
			{Timestamp: 3, Data: false},
			{Timestamp: 4, Data: "hi"},
		},
		Output: []pipescript.Datapoint{
			{Timestamp: 1, Data: 1},
			{Timestamp: 2, Data: true},
		},

		OutputError: true,
	}.Run(t)

	pipescript.TestCase{
		Pipescript: "filter $ < 5",
		Input: []pipescript.Datapoint{
			{Timestamp: 1, Data: 1},
			{Timestamp: 2, Data: 8},
			{Timestamp: 3, Data: false},
		},
		Output: []pipescript.Datapoint{
			{Timestamp: 1, Data: 1},
			{Timestamp: 3, Data: false},
		},
	}.Run(t)

	pipescript.TestCase{
		Pipescript: "filter $ < 5 | $ >= 3",
		Input: []pipescript.Datapoint{
			{Timestamp: 1, Data: 1},
			{Timestamp: 2, Data: 10},
			{Timestamp: 3, Data: 7},
			{Timestamp: 4, Data: 1.0},
			{Timestamp: 5, Data: 3},
			{Timestamp: 6, Data: 2.0},
			{Timestamp: 7, Data: 3.14},
		},
		Output: []pipescript.Datapoint{
			{Timestamp: 1, Data: false},
			{Timestamp: 4, Data: false},
			{Timestamp: 5, Data: true},
			{Timestamp: 6, Data: false},
			{Timestamp: 7, Data: true},
		},
	}.Run(t)

	pipescript.TestCase{
		Pipescript: "filter($ < 5):($ >= 3)",
		Input: []pipescript.Datapoint{
			{Timestamp: 1, Data: 1},
			{Timestamp: 2, Data: 10},
			{Timestamp: 3, Data: 7},
			{Timestamp: 4, Data: 1.0},
			{Timestamp: 5, Data: 3},
			{Timestamp: 6, Data: 2.0},
			{Timestamp: 7, Data: 3.14},
		},
		Output: []pipescript.Datapoint{
			{Timestamp: 1, Data: false},
			{Timestamp: 4, Data: false},
			{Timestamp: 5, Data: true},
			{Timestamp: 6, Data: false},
			{Timestamp: 7, Data: true},
		},
	}.Run(t)

	pipescript.TestCase{
		// This tests order of prescedence: ":" pipes are high prescedence, and will be executed first
		Pipescript: "filter ($['test']:$ < 5) | $['test']",
		Input: []pipescript.Datapoint{
			{Timestamp: 1, Data: map[string]interface{}{"test": 4}},
			{Timestamp: 2, Data: map[string]interface{}{"test": 8}},
			{Timestamp: 3, Data: map[string]interface{}{"test": 3}},
		},
		Output: []pipescript.Datapoint{
			{Timestamp: 1, Data: 4},
			{Timestamp: 3, Data: 3},
		},
	}.Run(t)

	pipescript.TestCase{
		// This tests order of prescedence: ":" pipes are high prescedence, and will be executed first
		Pipescript: "filter $['test']:$ < 5 | $['test']",
		Input: []pipescript.Datapoint{
			{Timestamp: 1, Data: map[string]interface{}{"test": 4}},
			{Timestamp: 2, Data: map[string]interface{}{"test": 8}},
			{Timestamp: 3, Data: map[string]interface{}{"test": 3}},
		},
		Output: []pipescript.Datapoint{
			{Timestamp: 1, Data: 4},
			{Timestamp: 3, Data: 3},
		},
	}.Run(t)

	pipescript.TestCase{
		// This tests order of prescedence: ":" pipes are high prescedence, and will be executed first
		Pipescript: "filter $:5 > $['test']:$:$:$ | $['test']:$",
		Input: []pipescript.Datapoint{
			{Timestamp: 1, Data: map[string]interface{}{"test": 4}},
			{Timestamp: 2, Data: map[string]interface{}{"test": 8}},
			{Timestamp: 3, Data: map[string]interface{}{"test": 3}},
		},
		Output: []pipescript.Datapoint{
			{Timestamp: 1, Data: 4},
			{Timestamp: 3, Data: 3},
		},
	}.Run(t)

}
