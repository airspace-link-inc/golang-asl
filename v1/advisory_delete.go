package aslv1

import (
	"context"
	"fmt"
	"net/http"
)

// Delete an advisory by its ID.
// Advisory ids are in UUID V4 format.
func (service *AdvisoryService) Delete(ctx context.Context, advisoryID string) (bool, error) {
	if advisoryID == "" {
		return false, fmt.Errorf("airspacelink/v1: advisory id must be a non-empty string")
	}

	url := fmt.Sprintf("%s/v4/advisory/%s", service.client.baseUrl, advisoryID)

	status, body, err := service.client.delete(ctx, url, AdvisoryRead.String())
	if err != nil {
		return false, err
	}

	if status != http.StatusOK {
		return false, fmt.Errorf("airspacelink/v1: status %d response: %s", status, body)
	}

	return true, nil
}
