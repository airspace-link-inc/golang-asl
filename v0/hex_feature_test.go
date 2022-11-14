package asl

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/h3-go/v3"
)

func TestHexFeatureMarshal(mainTest *testing.T) {
	testCases := []struct {
		name        string
		arg         HexFeature
		expected    []byte
		expectedErr error
	}{
		{
			name:     "base case",
			expected: []byte(`{"hexes":[],"props":null}`),
		},
		{
			name: "single hex (note-for many hexes, when marshalling into a string array, the order isn't deterministic, we don't have more than 1 hex in this test for this reason)",
			arg: HexFeature{Hexes: map[h3.H3Index]bool{
				172893179283: true,
			}},
			expected: []byte(`{"hexes":["28413c8d93"],"props":null}`),
		},
		{
			name: "properties only",
			arg: HexFeature{Props: map[string]any{
				"prop": 1,
			}},
			expected: []byte(`{"hexes":[],"props":{"prop":1}}`),
		},
		{
			name: "hexes and properties",
			arg: HexFeature{
				Hexes: map[h3.H3Index]bool{
					182938102: true,
				},
				Props: map[string]any{
					"prop": 1,
				},
			},
			expected: []byte(`{"hexes":["ae769f6"],"props":{"prop":1}}`),
		},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		actual, actualErr := tc.arg.MarshalJSON()
		t.Equal(tc.expected, actual, tc.name)
		t.Equal(tc.expectedErr, actualErr, tc.name)
	}
}

func TestHexFeatureUnmarshal(mainTest *testing.T) {
	testCases := []struct {
		name        string
		arg         []byte
		expected    *HexFeature
		expectedErr error
	}{
		{
			name: "base case",
			arg:  []byte(`null`),
		},
		{
			name: "single hex (note-for many hexes, when marshalling into a string array, the order isn't deterministic, we don't have more than 1 hex in this test for this reason)",
			arg:  []byte(`{"hexes":["28413c8d93"],"props":null}`),
			expected: &HexFeature{Hexes: map[h3.H3Index]bool{
				172893179283: true,
			}},
		},
		{
			name: "properties only",
			arg:  []byte(`{"hexes":[],"props":{"prop":1}}`),
			expected: &HexFeature{
				Hexes: map[h3.H3Index]bool{},
				Props: map[string]any{
					"prop": 1.0,
				},
			},
		},
		{
			name: "hexes and properties",
			arg:  []byte(`{"hexes":["ae769f6"],"props":{"prop":1}}`),
			expected: &HexFeature{
				Hexes: map[h3.H3Index]bool{
					182938102: true,
				},
				Props: map[string]any{
					"prop": 1.0,
				},
			},
		},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		var actual *HexFeature
		actualErr := json.Unmarshal(tc.arg, &actual)
		if tc.expectedErr == nil {
			if t.Nil(actualErr, tc.name) {
				t.Equal(tc.expected, actual, tc.name)
			}
			continue
		}

		if t.Nil(actual, tc.name) {
			t.Equal(tc.expectedErr, actualErr, tc.name)
		}
	}
}
