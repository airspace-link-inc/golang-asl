package aslv1

type AslServerResp[T any] struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}
