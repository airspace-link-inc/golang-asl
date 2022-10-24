package asl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AnthonyHewins/dpm"
	"github.com/peterstace/simplefeatures/geom"
)

// QueryAdvisoriesArgs is the payload passed to query advisories
// via a POST or GET. Prefer using GET if possible, since it leverages
// bounding box intersection which has a cheaper computation cost
// and a faster turnaround time
type QueryAdvisoriesArgs struct {
	Geom          geom.Geometry `json:"geometry"`
	AltitudeUpper float64       `json:"altitudeUpper"`
	AltitudeLower float64       `json:"altitudeLower"`
	StartTime     time.Time     `json:"startTime"`
	EndTime       time.Time     `json:"endTime"`
	GeoIDs        []string      `json:"geoIDs"`
}

//go:generate enumer -type AdvisoryCategoryType -json -text -transform lower
type AdvisoryCategoryType byte

const (
	Emergency AdvisoryCategoryType = iota + 1
	Recreational
	Admin
)

// Advisories represent geographic areas where special consideration must be
// made before operating drones. Examples of advisories may range from restricted
// airspace - where it's illegal to operate a drone - to warnings where it's important
// for you to understand the context of where you're flying to maximize safety.
// Additional advisory types will be documented as they're added
type Advisory struct {
	ID               string               `json:"id"`
	GeoID            string               `json:"geoID"`
	AdvisoryCategory AdvisoryCategoryType `json:"advisoryCategory"`
	Name             string               `json:"name"`
	Tags             []string             `json:"tags"`

	// Geospatial/temporal fields
	AltitudeLower float64       `json:"altitudeLower"`
	AltitudeUpper float64       `json:"altitudeUpper"`
	Geometry      geom.Geometry `json:"geometry"` // mark this as unmarshallable for json.RawMessage
	StartTime     time.Time     `json:"startTime"`
	EndTime       time.Time     `json:"endTime"`
	TimezoneName  string        `json:"timezoneName"`

	// Contact metadata
	ContactEmail    *string `json:"contactEmail"`
	ContactPhone    *string `json:"contactPhone"`
	CountryGeoID    string  `json:"countryGeoID"`
	URLString       *string `json:"url"`
	ReferenceNumber *string `json:"referenceNumber"`

	// Access control metadata
	CreatedBy    string `json:"createdBy"`
	LastEditedBy string `json:"lastEditedBy"`
	Published    bool   `json:"published"`

	// ASTM access control metadata
	OVN     string `json:"ovn"`
	Version int    `json:"version"`
}

func (a *Advisory) UnmarshalJSON(buf []byte) error {
	// unmarshal the geometry from the geojson feature
	// and then capture the properties into the struct:
	//
	// {
	//    "type": "Feature" <- ignore this,
	//    "geometry": {geojson geometry...}, <- snag this,
	//    "properties": {...} <- but don't unmarshal this yet
	// }
	type asGJFeature struct {
		Geom  geom.Geometry   `json:"geometry"`
		Props json.RawMessage `json:"properties"`
	}

	var gjFeature asGJFeature
	if err := json.Unmarshal(buf, &gjFeature); err != nil {
		return err
	}

	// now capture properties. Cast to anonymous struct first to avoid
	// recursion loop
	if err := json.Unmarshal(gjFeature.Props, any(a)); err != nil {
		return err
	}

	// throw geojson geometry on top of it, and we're done
	a.Geometry = gjFeature.Geom
	return nil
}

func (a *Advisory) MarshalJSON() ([]byte, error) {
	return json.Marshal(geom.GeoJSONFeature{
		Properties: dpm.Params(a).Tag("json").Omit("geometry").FilterZero().Map(),
		Geometry:   a.Geometry,
	})
}

func (a *Client) QueryAdvisoriesByGeom(ctx context.Context, args *QueryAdvisoriesArgs) (*Resp[[]Advisory], error) {
	req, err := a.makeJSONReq(ctx, "POST", "/v4/advisories", args)
	if err != nil {
		return nil, err
	}

	return apiReq[[]Advisory](&a.HTTPClient, req)
}
