package responses

type BadRequestError struct {
	Error string `json:"error"`
}

type InternalServerError struct {
	Error string `json:"error"`
}
