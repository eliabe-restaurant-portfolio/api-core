package producers

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/connections"
	sendresetpasswordemailproducer "github.com/eliabe-portfolio/restaurant-app/internal/queues/producers/send-reset-password-email"
)

type Provider interface {
	SendPasswordResetEmail() sendresetpasswordemailproducer.Producer
}

type producers struct {
	connections *connections.Provider
}

func New(connections *connections.Provider) Provider {
	return producers{connections: connections}
}

func (p producers) SendPasswordResetEmail() sendresetpasswordemailproducer.Producer {
	return sendresetpasswordemailproducer.New(p.connections)
}
