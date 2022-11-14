package aslv1

import (
	"context"
	"fmt"
	"time"
)

type CheckService service

type CheckAirspaceOptions struct {
	// GeoJSON formatted geometry. Available geometries may be a polygon.
	Geometry any `json:"geometry"`
	// RFC3339 compliant start timestamp of the operation.
	StartTime time.Time `json:"startTime"`
	// RFC3339 compliant start timestamp of the operation.
	EndTime time.Time `json:"endTime"`
}

type AirspaceDetails struct {
	// Represents airspace that is monitored and managed by air traffic control. If true, rules implemented by air traffic apply.
	Controlled bool `json:"controlled"`
	// Applies to controlled grids/areas with altitude ceilings. At times, these grids may be disabled.
	// If true, the altitude ceilings are in effect and require authorization in accordance with the rules and ceilings.
	// If false, alternative authorization is likely required. For example, the Drone Zone in the US.
	Enabled bool `json:"enabled"`
	// Denotes airspace that is subject to limitations or may be prohibited entirely.
	// If true, drone flights cannot be authorized during a time window or in some cases into perpetuity.
	Restricted bool `json:"restricted"`
}

// Returns a payload containing properties that inform whether Airspace Authorization is required
// Examples:
//  1. Airspace is uncontrolled and does not require authorization { "controlled": false, "enabled": true | false, "restricted": false }
//  2. Airspace is controlled with enabled grids/areas and requires authorization { "controlled": true, "enabled": true, "restricted": false }
//  3. Airspace is restricted and drone flights are not allowed { "controlled": false, "enabled": false, "restricted": true }
//  4. Airspace is controlled but grids/areas are not enabled and alternative authorization is required { "controlled": true, "enabled": false, "restricted": false }
func (cs *CheckService) Airspace(ctx context.Context, opts CheckAirspaceOptions) (*AirspaceDetails, error) {
	url := fmt.Sprintf("%s/v1/check/airspace", cs.client.baseUrl)

	status, body, err := cs.client.post(ctx, url, opts, AviationRead.String())
	if err != nil {
		return nil, err
	}

	var airspaceDetails AirspaceDetails
	err = unmarshalStatusOk(status, body, &airspaceDetails)
	if err != nil {
		return nil, err
	}

	return &airspaceDetails, nil
}
