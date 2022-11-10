package asl

// resp is the generic response you get back from any
// API call
type Resp[X any] struct {
	Status int    `json:"statusCode"`
	Msg    string `json:"message"`
	Data   X      `json:"data"`
}
