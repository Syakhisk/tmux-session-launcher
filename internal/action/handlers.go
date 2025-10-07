package action

import (
	"context"
	"tmux-session-launcher/internal/client"
	"tmux-session-launcher/internal/rpc"

	"github.com/urfave/cli/v3"
)

func HandlerNextMode(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(rpc.SockAddress)
	action := NewAction(c)

	return action.NextMode(ctx)
}

func HandlerPrevMode(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(rpc.SockAddress)
	action := NewAction(c)

	return action.PrevMode(ctx)
}

func HandlerGetMode(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(rpc.SockAddress)
	action := NewAction(c)

	return action.GetMode(ctx)
}

func HandlerGetContent(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(rpc.SockAddress)
	action := NewAction(c)

	return action.GetContent(ctx)
}

func HandlerOpenIn(ctx context.Context, cmd *cli.Command) error {
	c := client.NewClient(rpc.SockAddress)
	action := NewAction(c)

	args := cmd.Args().Slice()
	if len(args) < 1 || len(args) > 2 {
		return cli.Exit("invalid number of arguments", 1)
	}

	return action.OpenIn(ctx, args[0])
}
