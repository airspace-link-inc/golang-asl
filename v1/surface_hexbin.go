package aslv1

import (
	"context"
	"encoding/json"
	"fmt"

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

type SurfaceHexbinOptions struct {
	// A GeoJSON geometry representing the area you want to query
	Geometry geom.Geometry `json:"geometry"`

	// Layers you want to query
	Layers []Layer `json:"layers"`

	// Resolution you want to use for hexagon indexing. The highest
	// possible is 16, but
	Resolution uint8 `json:"resolution"`
}

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

// Generate a surface with any data sources you request.
// Returns a hexbin payload that includes hexes contained within the area you specify in the geometry parameter.
// For a GeoJSON response, use the Surface GeoJSON endpoint.
func (ss *SurfaceService) Surface(ctx context.Context, opts *SurfaceHexbinOptions) ([]HexFeature, error) {
	url := fmt.Sprintf("%s/v2/surface", ss.client.baseUrl)

	status, body, err := ss.client.post(ctx, url, opts, SurfaceTier1.String())
	if err != nil {
		return nil, err
	}

	var features []HexFeature
	err = unmarshalStatusOk(status, body, &features)
	if err != nil {
		return nil, err
	}

	return features, nil
}
