package responses

type BadRequestError struct {
	Error string `json:"error"`
}
