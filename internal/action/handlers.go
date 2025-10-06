package action

import (
	"context"
	"tmux-session-launcher/internal/client"
	"tmux-session-launcher/internal/constants"

	"github.com/urfave/cli/v3"
)

func HandlerActionNextMode(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(constants.SockAddress)
	action := NewAction(c)

	return action.NextMode(ctx)
}

func HandlerActionPrevMode(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(constants.SockAddress)
	action := NewAction(c)

	return action.PrevMode(ctx)
}

func HandlerActionGetMode(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(constants.SockAddress)
	action := NewAction(c)
	return action.GetMode(ctx)
}

func HandlerActionGetContent(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(constants.SockAddress)
	action := NewAction(c)
	return action.GetContent(ctx)
}

func HandlerActionOpenIn(ctx context.Context, cmd *cli.Command) error {
	return nil
}
