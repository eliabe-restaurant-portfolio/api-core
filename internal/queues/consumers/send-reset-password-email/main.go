package sendresetpasswordemailconsumer

import (
	"encoding/json"
	"log"

	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	sendresetpasswordemailcmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/notification/email/send-reset-pwd-email"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type SendPasswordResetEmailMessage struct {
	ResetPasswordToken string
	RandomHash         string
	RandomHashs        string
}

type consumer struct {
	adapters *adapters.Adapters
}

func New(adapters *adapters.Adapters) consumer {
	return consumer{adapters}
}

func (c consumer) Process(d amqp091.Delivery) {
	var body SendPasswordResetEmailMessage
	if err := json.Unmarshal(d.Body, &body); err != nil {
		log.Printf("json unmarshal error for queue %q: %v", constants.Queues.ResetPasswordEmail, err)
		_ = d.Reject(false)
		return
	}

	params := sendresetpasswordemailcmd.Params{
		ResetPasswordToken: uuid.MustParse(body.ResetPasswordToken),
		RandomHash:         body.RandomHash,
	}

	log.Printf("Processing password reset email message: %+v", params)

	useCase := sendresetpasswordemailcmd.New(c.adapters)
	res, err := useCase.Execute(params)
	if err != nil {
		log.Printf("Error processing message for password reset email: %v. Response: %+v", err, res)
		_ = d.Reject(true)
		return
	}

	log.Printf("Successfully processed password reset email for token: %s", body.ResetPasswordToken)
	_ = d.Ack(false)
}
