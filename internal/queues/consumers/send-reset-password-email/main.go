package sendresetpasswordemailconsumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
	sendresetpasswordemailcmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/notification/email/send-reset-pwd-email"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type SendPasswordResetEmailMessage struct {
	ResetPasswordToken string
	Token              string
}

type consumer struct {
	repositories *repositories.Provider
}

func New(repositories *repositories.Provider) consumer {
	return consumer{repositories: repositories}
}

func (c consumer) Process(ctx context.Context, d amqp091.Delivery) {
	var body SendPasswordResetEmailMessage
	if err := json.Unmarshal(d.Body, &body); err != nil {
		log.Printf("json unmarshal error for queue %q: %v", constants.Queues.ResetPasswordEmail, err)
		_ = d.Reject(false)
		return
	}

	params := sendresetpasswordemailcmd.Params{
		PasswordResetToken: uuid.MustParse(body.ResetPasswordToken),
		Token:              body.Token,
	}

	log.Printf("Processing password reset email message: %+v", params)

	useCase := sendresetpasswordemailcmd.New(c.repositories)
	res, err := useCase.Execute(params)
	if err != nil {
		log.Printf("Error processing message for password reset email: %v. Response: %+v", err, res)
		_ = d.Reject(true)
		return
	}

	log.Printf("Successfully processed password reset email for token: %s", body.ResetPasswordToken)
	_ = d.Ack(false)
}
