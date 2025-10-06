package action

import (
	"context"
	"fmt"
	"tmux-session-launcher/internal/client"
	"tmux-session-launcher/internal/constants"
	"tmux-session-launcher/pkg/logger"

	"emperror.dev/errors"
)

type Action struct {
	client *client.Client
}

func NewAction(client *client.Client) *Action {
	return &Action{
		client: client,
	}
}

func (a *Action) NextMode(ctx context.Context) error {
	logger := logger.WithPrefix("action.NextMode")

	res, err := a.client.Send(constants.RouteNextMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	logger.Debugf("Raw response: %s", res)

	return nil
}

func (a *Action) PrevMode(ctx context.Context) error {
	log := logger.WithPrefix("action.PrevMode")

	res, err := a.client.Send(constants.RoutePrevMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	log.Debugf("Raw response: %s", res)

	return nil
}

func (a *Action) GetMode(ctx context.Context) error {
	res, err := a.client.Send(constants.RouteGetMode, "")
	if err != nil {
		return errors.Wrap(err, "failed to get current mode")
	}

	fmt.Println("Current mode:", res)

	return nil
}

func (a *Action) GetContent(ctx context.Context) error {
	res, err := a.client.Send(constants.RouteGetContent, "")
	if err != nil {
		return errors.Wrap(err, "failed to get current content")
	}

	fmt.Println(res)

	return nil
}
