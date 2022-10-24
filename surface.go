package asl

import (
	"context"
	"net/http"

	"github.com/peterstace/simplefeatures/geom"
)

type Layer struct {
	Alias  string   `json:"alias"`
	Fields []string `json:"fields"`
	Code   string   `json:"code"`
	Where  []any    `json:"where"`
	Score  float64  `json:"score"`
}

type SurfaceReq struct {
	// A GeoJSON geometry representing the area you want to query
	Geometry geom.Geometry `json:"geometry"`

	// Layers you want to query
	Layers []Layer `json:"layers"`

	// Resolution you want to use for hexagon indexing. The highest
	// possible is 16, but
	Resolution uint8 `json:"resolution"`
}

func (c Client) Surface(ctx context.Context, req *SurfaceReq) (*Resp[[]HexFeature], error) {
	httpReq, err := c.makeJSONReq(ctx, http.MethodPost, "/v2/surface", req)
	if err != nil {
		return nil, err
	}

	return apiReq[[]HexFeature](&c.HTTPClient, httpReq)
}
