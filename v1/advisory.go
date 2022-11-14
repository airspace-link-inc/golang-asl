package aslv1

import (
	"time"

	"github.com/google/uuid"
)

type AdvisoryService service

// Advisories represent geographic areas where special consideration must be made before operating drones.
// Examples of advisories may range from restricted airspace where it's illegal to operate a drone
//
//	to warnings where it's important for you to understand the context of where you're flying to maximize safety.
type Advisory struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Tags               []string  `json:"tags"`
	GeoID              string    `json:"geoID"`
	CountryGeoID       string    `json:"countryGeoID"`
	Geometry           any       `json:"geom"`
	AltitudeUpper      int32     `json:"altitudeUpper"`
	AltitudeLower      int32     `json:"altitudeLower"`
	ReferenceNumber    *string   `json:"referenceNumber,omitempty"`
	URL                *string   `json:"url,omitempty"`
	AdvisoryCategoryID string    `json:"advisoryCategory"`
	Published          bool      `json:"published"`
	ContactEmail       *string   `json:"contactEmail,omitempty"`
	ContactPhone       *string   `json:"contactPhone,omitempty"`
	TimezoneName       string    `json:"timezoneName"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
}
