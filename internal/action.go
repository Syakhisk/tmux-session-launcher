package internal

import (
	"context"
	"fmt"
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

func HandlerActionGetMode(ctx context.Context, cmd *cli.Command) error {
	client := NewClient(SockAddress)
	action := NewAction(client)
	return action.GetMode(ctx)
}

func HandlerActionGetContent(ctx context.Context, cmd *cli.Command) error {
	client := NewClient(SockAddress)
	action := NewAction(client)
	return action.GetContent(ctx)
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
	logger := logger.WithPrefix("action.NextMode")

	res, err := a.client.Send(RouteNextMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	logger.Debugf("Raw response: %s", res)

	return nil
}

func (a *Action) PrevMode(ctx context.Context) error {
	logger := logger.WithPrefix("action.PrevMode")

	res, err := a.client.Send(RoutePrevMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	logger.Debugf("Raw response: %s", res)

	return nil
}

func (a *Action) GetMode(ctx context.Context) error {
	res, err := a.client.Send(RouteGetMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get current mode")
	}

	fmt.Println("Current mode:", res)

	return nil
}

func (a *Action) GetContent(ctx context.Context) error {
	res, err := a.client.Send(RouteGetContent, "")
	if err != nil {
		return errors.Wrap(err, "failed to get current content")
	}

	fmt.Println(res)

	return nil
}
