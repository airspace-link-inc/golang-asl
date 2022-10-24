package asl

import (
	"encoding/json"

	"github.com/uber/h3-go"
)

// HexFeature is an h3 equivalent of GeoJSON. It comes in very similar to a GeoJSON
// geometry, with a properties member (shortened to props) and a "geometry":
//
//	{
//		"hexes": ["892ab2c106bffff", ... ],
//		"properties": {...}
//	}
type HexFeature struct {
	Hexes map[h3.H3Index]bool
	Props map[string]any
}

func (s *HexFeature) UnmarshalJSON(buf []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(buf, &m); err != nil {
		return err
	} else if m == nil {
		return nil
	}

	var props map[string]any
	if err := json.Unmarshal(m["props"], &props); err != nil {
		return err
	}

	var hexSlice []string
	if err := json.Unmarshal(m["hexes"], &hexSlice); err != nil {
		return err
	}

	hexes := make(map[h3.H3Index]bool, len(hexSlice))
	for _, v := range hexSlice {
		hexes[h3.FromString(v)] = true
	}

	s.Props = props
	s.Hexes = hexes
	return nil
}

func (s *HexFeature) MarshalJSON() ([]byte, error) {
	hexStrings, i := make([]string, len(s.Hexes)), 0
	for k := range s.Hexes {
		hexStrings[i] = h3.ToString(k)
		i++
	}

	return json.Marshal(map[string]any{
		"props": s.Props,
		"hexes": hexStrings,
	})
}
