package returns

type Api struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    any      `json:"data"`
	Code    int      `json:"code"`
	Errors  []string `json:"errors"`
}

func Success(message string, data any) Api {
	return Api{
		Success: true,
		Message: message,
		Data:    data,
		Code:    200,
	}
}

func InternalServerError(errs []string) Api {
	errors := []string{"internal server error"}
	errors = append(errors, errs...)

	return Api{
		Success: false,
		Message: "Ocorreu um erro interno.",
		Data:    nil,
		Code:    500,
		Errors:  errors,
	}
}

func BadRequest(message string) Api {
	return Api{
		Success: false,
		Message: message,
		Data:    nil,
		Code:    400,
		Errors:  []string{"Bad request"},
	}
}

func Unauthorized(message string) Api {
	return Api{
		Success: false,
		Message: message,
		Data:    nil,
		Code:    401,
		Errors:  []string{"Authentication required"},
	}
}

func Forbidden(message string) Api {
	return Api{
		Success: false,
		Message: message,
		Data:    nil,
		Code:    403,
		Errors:  []string{"Permission denied"},
	}
}

func NotFound(message string) Api {
	return Api{
		Success: false,
		Message: message,
		Data:    nil,
		Code:    404,
		Errors:  []string{"Not Found"},
	}
}
