package sendresetpasswordemailcmd

import "github.com/eliabe-portfolio/restaurant-app/pkg/returns"

type MessageProvider struct {
	Success func() returns.Api
	Default func() returns.Api
}

var messageProvider = MessageProvider{
	Success: func() returns.Api { return returns.Success("Email enviado com sucesso.", nil) },
	Default: func() returns.Api { return returns.InternalServerError([]string{}) },
}
