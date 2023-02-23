package bot

import (
	"errors"

	utopiago "github.com/Sagleft/utopialib-go/v2"
)

type uBot struct {
	conn utopiago.Client
}

func NewUtopiaBot(cfg utopiago.Config) (Bot, error) {
	b := &uBot{
		conn: utopiago.NewUtopiaClient(cfg),
	}

	if !b.conn.CheckClientConnection() {
		return nil, errors.New("failed to connect to Utopia client")
	}

	return b, nil
}
