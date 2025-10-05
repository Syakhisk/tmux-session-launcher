package internal

import (
	"context"
	"fmt"

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

	fmt.Printf("Got response: %s\n", res)

	return nil
}

func (a *Action) PrevMode(ctx context.Context) error {
	res, err := a.client.Send(RoutePrevMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	fmt.Printf("Got response: %s\n", res)

	return nil
}

func (a *Action) GetMode(ctx context.Context) error {
	res, err := a.client.Send(RouteGetMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get current mode")
	}

	fmt.Printf("Got response: %s\n", res)

	return nil
}
