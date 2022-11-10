package aslv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type QueryAdvisoryOptions struct {
	// GeoJSON formatted geometry. Must be a point, line, or polygon.
	// The resulting bounding box of the geometry can't have a side length greater than 100 miles.
	Geometry any `json:"geometry"`
	// Lower bound of the time range filter. Returns advisories that are active at the same time as, or after startTime.
	// Send as YYYY-MM-DDThh:mmZ.
	StartTime *time.Time `json:"startTime,omitempty"`
	//  	Upper bound of the time range filter. Returns advisories that are active at the same time as, or before endTime.
	// Send as YYYY-MM-DDThh:mmZ.
	EndTime *time.Time `json:"endTime,omitempty"`
	// Filter advisories that are flying at or above the specified altitude.
	AltitudeLower *int `json:"altitudeLower,omitempty"`
	// Filter advisories that are flying at or below the specified altitude.
	AltitudeUpper *int `json:"altitudeUpper,omitempty"`
	// Filter advisories that are present in the following geo-IDs.
	//Limit of 5 IDs.
	GeoIDs []string `json:"geoIDs,omitempty"`
}

// Find advisories that intersect the input geometry.
func (service *AdvisoryService) Query(ctx context.Context, opts QueryAdvisoryOptions) ([]Advisory, error) {

	url := fmt.Sprintf("%s/v4/advisory/query", service.client.baseUrl)

	status, body, err := service.client.post(ctx, url, opts, AdvisoryRead.String())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("airspacelink/v1: status %d response: %s", status, body)
	}

	var advisories []Advisory
	err = json.Unmarshal(body, &advisories)
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: unable to unmarshal response to advisory: %w", err)
	}

	return advisories, nil
}
