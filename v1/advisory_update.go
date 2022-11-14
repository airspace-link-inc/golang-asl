package aslv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UpdateAdvisoryOptions struct {
	Name               *string    `json:"name,omitempty"`
	Tags               *[]string  `json:"tags,omitempty"`
	GeoID              *string    `json:"geoID,omitempty"`
	CountryGeoID       *string    `json:"countryGeoID,omitempty"`
	Geometry           any        `json:"geom,omitempty"`
	AltitudeUpper      *int       `json:"altitudeUpper,omitempty"`
	AltitudeLower      *int       `json:"altitudeLower,omitempty"`
	ReferenceNumber    *string    `json:"referenceNumber,omitempty"`
	URL                *string    `json:"url,omitempty"`
	AdvisoryCategoryID *string    `json:"advisoryCategory,omitempty"`
	Published          *bool      `json:"published,omitempty"`
	ContactEmail       *string    `json:"contactEmail,omitempty"`
	ContactPhone       *string    `json:"contactPhone,omitempty"`
	TimezoneName       *string    `json:"timezoneName,omitempty"`
	StartTime          *time.Time `json:"startTime,omitempty"`
	EndTime            *time.Time `json:"endTime,omitempty"`
}

func (service *AdvisoryService) Update(ctx context.Context, advisoryID string, opts UpdateAdvisoryOptions) (*Advisory, error) {
	url := fmt.Sprintf("%s/v4/advisory/%s", service.client.baseUrl, advisoryID)

	status, body, err := service.client.patch(ctx, url, opts, AdvisoryRead.String())
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
