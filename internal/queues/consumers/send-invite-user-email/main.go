package sendinviteuseremailconsumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	sendresetpasswordemailcmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/notification/email/send-reset-pwd-email"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type SendInviteUserEmailMessage struct {
	UserToken string
}

type consumer struct {
	adapters *adapters.Adapters
}

func New(adapters *adapters.Adapters) consumer {
	return consumer{adapters: adapters}
}

func (c consumer) Process(ctx context.Context, d amqp091.Delivery) {
	var body SendInviteUserEmailMessage
	if err := json.Unmarshal(d.Body, &body); err != nil {
		log.Printf("json unmarshal error for queue %q: %v", constants.Queues.ResetPasswordEmail, err)
		_ = d.Reject(false)
		return
	}

	params := sendresetpasswordemailcmd.Params{
		ResetPasswordToken: uuid.MustParse(body.UserToken),
	}

	log.Printf("Processing password reset email message: %+v", params)

	useCase := sendresetpasswordemailcmd.New(c.adapters)
	res, err := useCase.Execute(params)
	if err != nil {
		log.Printf("Error processing message for password reset email: %v. Response: %+v", err, res)
		_ = d.Reject(true)
		return
	}

	log.Printf("Successfully processed password reset email for token: %s", body.UserToken)
	_ = d.Ack(false)
}
