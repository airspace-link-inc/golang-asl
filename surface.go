package asl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

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
