package action

import (
	"context"
	"fmt"
	"tmux-session-launcher/internal/client"
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

	res, err := a.client.NextMode(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get next mode")
	}

	logger.Debugf("Raw response: %s", res.Mode)

	return nil
}

func (a *Action) PrevMode(ctx context.Context) error {
	log := logger.WithPrefix("action.PrevMode")

	res, err := a.client.PrevMode(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get previous mode")
	}

	log.Debugf("Raw response: %s", res.Mode)

	return nil
}

func (a *Action) GetMode(ctx context.Context) error {
	res, err := a.client.GetMode(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get current mode")
	}

	fmt.Println(res.Mode)

	return nil
}

func (a *Action) GetContent(ctx context.Context) error {
	res, err := a.client.GetContent(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get current content")
	}

	fmt.Println(res.Content)

	return nil
}

func (a *Action) OpenIn(ctx context.Context, selectionString string) error {
	log := logger.WithPrefix("action.OpenIn")

	err := a.client.OpenIn(ctx, selectionString)
	if err != nil {
		return errors.Wrap(err, "failed to open in")
	}

	log.Debugf("Successfully opened selection: %s", selectionString)

	return nil
}
