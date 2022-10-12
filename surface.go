package asl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/peterstace/simplefeatures/geom"
	"github.com/uber/h3-go/v3"
)

type Layer struct {
	Alias  string   `json:"alias"`
	Fields []string `json:"fields"`
	Code   string   `json:"code"`
	Where  []any    `json:"where"`
	Score  float64  `json:"score"`
}

type SurfaceV2Req struct {
	// A GeoJSON geometry representing the area you want to query
	Geometry geom.Geometry `json:"geometry"`

	// Layers you want to query
	Layers []Layer `json:"layers"`

	// Resolution you want to use for hexagon indexing. The highest
	// possible is 16, but
	Resolution uint8 `json:"resolution"`
}

type SurfaceV2HexResp struct {
	Hexes []h3.H3Index   `json:"hexes"`
	Props map[string]any `json:"props"`
}

func (s *SurfaceV2HexResp) UnmarshalJSON(buf []byte) error {
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

	hexes := make([]h3.H3Index, len(stringSlice))
	for i, v := range stringSlice {
		hexString, ok := v.(string)
		if !ok {
			return &json.MarshalerError{
				Type: reflect.TypeOf(""),
				Err:  fmt.Errorf("all elements in hexes must be strings"),
			}
		}

		hexes[i] = h3.FromString(hexString)
	}
	s.Hexes = hexes

	return nil
}

func (c Client) SurfaceV2(ctx context.Context, req *SurfaceV2Req) (*Resp[[]SurfaceV2HexResp], error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := c.makeReq(ctx, http.MethodPost, "/v2/surface", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	return apiReq[[]SurfaceV2HexResp](&c.HTTPClient, httpReq)
}
