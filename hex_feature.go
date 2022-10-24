package asl

import (
	"encoding/json"
	"fmt"
	"reflect"

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
	var m map[string]any
	if err := json.Unmarshal(buf, &m); err != nil {
		return err
	}

	props, ok := m["props"].(map[string]any)
	if !ok {
		return &json.UnmarshalTypeError{
			Value:  "JSON object",
			Type:   reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(props[""])),
			Struct: "SurfaceV2HexResp",
			Field:  "hexes",
		}
	}
	s.Props = props

	stringSlice, ok := m["hexes"].([]any)
	if !ok {
		return &json.UnmarshalTypeError{
			Value:  "array of h3 indices (as strings)",
			Type:   reflect.SliceOf(reflect.TypeOf("")),
			Struct: "SurfaceV2HexResp",
			Field:  "hexes",
		}
	}

	hexes := make(map[h3.H3Index]bool, len(stringSlice))
	for _, v := range stringSlice {
		hexString, ok := v.(string)
		if !ok {
			return &json.MarshalerError{
				Type: reflect.TypeOf(""),
				Err:  fmt.Errorf("all elements in hexes must be strings"),
			}
		}

		hexes[h3.FromString(hexString)] = true
	}

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
