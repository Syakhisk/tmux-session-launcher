package internal

import (
	"context"
	"tmux-session-launcher/pkg/logger"

	"emperror.dev/errors"
	"github.com/urfave/cli/v3"
)

func ActionNextModeHandler(ctx context.Context, cmd *cli.Command) error {
	client := NewClient(SockAddress)
	action := NewAction(client)

	return action.NextMode(ctx)
}

func ActionPrevModeHandler(ctx context.Context, cmd *cli.Command) error {
	client := NewClient(SockAddress)
	action := NewAction(client)

	return action.PrevMode(ctx)
}

type Action struct {
	client *Client
}

func NewAction(client *Client) *Action {
	return &Action{
		client: client,
	}
}

func (a *Action) NextMode(ctx context.Context) error {
	res, err := a.client.Send("next", "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	logger.Infof("Switched to next mode: %s", res)

	return nil
}

func (a *Action) PrevMode(ctx context.Context) error {
	res, err := a.client.Send("prev", "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	logger.Infof("Switched to prev mode: %s", res)

	return nil
}
