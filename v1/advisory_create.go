package aslv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateAdvisoryOptions struct {
	Name         string `json:"name"`
	GeoID        string `json:"geoId"`
	CountryGeoID string `json:"countryGeoID,omitempty"`
	// // GeoJSON formatted geometry. Must be a point, line, or polygon.
	// The resulting bounding box of the geometry can't have a side length greater than 100 miles.
	Geometry         any    `json:"geom"`
	AltitudeUpper    int    `json:"altitudeUpper"`
	AltitudeLower    int    `json:"altitudeLower"`
	AdvisoryCategory string `json:"advisoryCategory"`
	TimezoneName     string `json:"timezoneName"`
	Published        bool   `json:"published,omitempty"`
	StartTime        string `json:"startTime,omitempty"`
	EndTime          string `json:"endTime,omitempty"`
	URL              string `json:"url,omitempty"`
	ContactEmail     string `json:"contactEmail,omitempty"`
	ReferenceNumber  string `json:"referenceNumber,omitempty"`
	ContactPhone     string `json:"contactPhone,omitempty"`
}

func (service *AdvisoryService) Create(ctx context.Context, opts CreateAdvisoryOptions) (*Advisory, error) {
	url := fmt.Sprintf("%s/v4/advisory", service.client.baseUrl)

	status, body, err := service.client.post(ctx, url, opts, AdvisoryRead.String())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("airspacelink/v1: status %d response: %s", status, body)
	}

	var advisory Advisory
	err = json.Unmarshal(body, &advisory)
	if err != nil {
		return nil, fmt.Errorf("airspacelink/v1: unable to unmarshal response to advisory: %w", err)
	}

	return &advisory, nil
}
