package aslv1

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/peterstace/simplefeatures/geom"
)

type AviationService service

//go:generate enumer -type=FaaType -transform=snake -json -text

// FAA airspace types represented by a constant string
type FaaType uint8

const (
	// Controlled Airspace Classification.
	ControlledAirspace FaaType = iota
	// UAS Facility Management Flight Ceiling.
	UasfmCeiling
	// Both sua_prohibited and sua_restricted.
	Sua
	// Washington DC Flight Restricted Zone.
	WashingtonFrz
	// Part Time National Security UAS Flight Restriction.
	NsufrPt
	// Full Time National Security UAS Flight Restriction.
	NsufrFt
	// Select stadiums containing temporary flight restriction. The stadium points are buffered by 3 nautical miles.
	Stadium
	// Points designated for landing or takeoff. Returns all airports within 3 nautical miles of input geometry by default.
	Airports
	// Controlled airspace schedule for select airports across the country (geometry has no bearing on the result).
	AirspaceSchedule
	// Temporary Flight Restrictions imposed by the FAA to restrict aircraft operations within designated areas.
	Tfr
)

type AviationOptions struct {
	// Note: only used when type includes airports.
	// Buffer specifies the buffer around the airport in nautical miles.
	// If not specified, the default is 3.
	Buffer *uint `url:"buffer,omitempty" json:"-"`
	// GeoJSON formatted geometry. Available geometries may be a point, line, or polygon.
	// The resulting bounding box of the geometry can't have a side length greater than 30 nautical miles.
	Geometry any `json:"geometry"`
	// An array of zero or more types from FAA types below.
	// Specifying type will limit the results to only include the types that you specify.
	// If you specify no type then all types will be returned.
	Type []FaaType `json:"type"`
}

// Returns intersecting aviation data from the FAA.
func (as *AviationService) IntersectionDataFAA(ctx context.Context, opts AviationOptions) ([]geom.GeoJSONFeature, error) {
	queryParams, err := query.Values(opts)
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: unable to parse query params: %w", err)
	}

	u, err := url.ParseRequestURI(fmt.Sprintf("%s/v1/aviation", as.client.baseUrl))
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: unable to construct URL: %w", err)
	}

	u.RawQuery = queryParams.Encode()

	status, body, err := as.client.post(ctx, u.String(), opts, AviationRead.String())
	if err != nil {
		return nil, err
	}

	// TODO(Anthony): This fails since some of the features returned have null geometries...  
	var features AslServerResp[[]geom.GeoJSONFeature]
	err = unmarshalStatusOk(status, body, &features)
	if err != nil {
		return nil, err
	}

	return features.Data, nil
}
