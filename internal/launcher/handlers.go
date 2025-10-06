package launcher

import (
	"context"
	"fmt"
	"tmux-session-launcher/internal/constants"
	"tmux-session-launcher/internal/fuzzyfinder"
	"tmux-session-launcher/internal/mode"

	"emperror.dev/errors"
)

func (l *Launcher) setupHandlers() {
	l.Server.RegisterHandler(constants.RouteNextMode, func(ctx context.Context, message string) (string, error) {
		m := mode.Next()

		err := fuzzyfinder.UpdateContentAndHeader(ctx)
		if err != nil {
			return "", errors.WrapIf(err, "failed to update fzf content and header")
		}

		return fmt.Sprintf("current mode: %s", m), nil
	})

	l.Server.RegisterHandler(constants.RoutePrevMode, func(ctx context.Context, message string) (string, error) {
		m := mode.Prev()

		err := fuzzyfinder.UpdateContentAndHeader(ctx)
		if err != nil {
			return "", errors.WrapIf(err, "failed to update fzf content and header")
		}

		return fmt.Sprintf("current mode: %s", m), nil
	})

	l.Server.RegisterHandler(constants.RouteGetMode, func(ctx context.Context, message string) (string, error) {
		m := mode.Get()

		return fmt.Sprintf("current mode: %s", m), nil
	})

	l.Server.RegisterHandler(constants.RouteGetContent, func(ctx context.Context, message string) (string, error) {
		return fuzzyfinder.GetContent(ctx)
	})

	l.Server.RegisterHandler(constants.RouteOpenIn, func(ctx context.Context, message string) (string, error) {
		return "", fuzzyfinder.OpenIn(ctx, message)
	})
}
