package asl

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/peterstace/simplefeatures/geom"
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

func (c Client) SurfaceV2(req *SurfaceV2Req) (*Resp[SurfaceV2Req], error) {
	urlBase := c.BaseURL
	if urlBase == "" {
		// point to prod as a default
		urlBase = "https://airhub-api.airspacelink.com"
	}

	uri, err := url.Parse(urlBase + "/v2/surface")
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return apiReq[SurfaceV2Req](&c.HTTPClient, &http.Request{
		Method: http.MethodPost,
		URL:    uri,
		Header: c.makeHeaders(),
		Body:   io.NopCloser(bytes.NewBuffer(buf)),
	})
}
