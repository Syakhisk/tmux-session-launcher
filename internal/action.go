package internal

import (
	"context"
	"tmux-session-launcher/pkg/logger"

	"emperror.dev/errors"
	"github.com/urfave/cli/v3"
)

func HandlerActionNextMode(ctx context.Context, cmd *cli.Command) error {
	client := NewClient(SockAddress)
	action := NewAction(client)

	return action.NextMode(ctx)
}

func HandlerActionPrevMode(ctx context.Context, cmd *cli.Command) error {
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
	res, err := a.client.Send(RouteNextMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	logger.Infof("Got response: %s", res)

	return nil
}

func (a *Action) PrevMode(ctx context.Context) error {
	res, err := a.client.Send(RoutePrevMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	logger.Infof("Got response: %s", res)

	return nil
}
