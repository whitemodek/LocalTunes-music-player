package mapper

import responses "backend/internal/handler/dto/response"

func BadRequestErrors(error string) responses.BadRequestError {
	return responses.BadRequestError{
		Error: error,
	}
}

func InternalServerError() responses.InternalServerError {
	return responses.InternalServerError{
		Error: "Internal server error",
	}
}
