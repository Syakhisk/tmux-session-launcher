package launcher

import (
	"context"
	"tmux-session-launcher/internal/constants"
	"tmux-session-launcher/internal/fuzzyfinder"
	"tmux-session-launcher/internal/mode"
	"tmux-session-launcher/internal/rpc"

	"emperror.dev/errors"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
)

func (l *Launcher) setupHandlers() {
	l.Server.RegisterHandler(constants.MethodModeNext, handler.New(func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
		m := mode.Next()

		err := fuzzyfinder.UpdateContentAndHeader(ctx)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to update fzf content and header")
		}

		return rpc.ModeResponse{Mode: m.String()}, nil
	}))

	l.Server.RegisterHandler(constants.MethodModePrev, handler.New(func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
		m := mode.Prev()

		err := fuzzyfinder.UpdateContentAndHeader(ctx)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to update fzf content and header")
		}

		return rpc.ModeResponse{Mode: m.String()}, nil
	}))

	l.Server.RegisterHandler(constants.MethodModeGet, handler.New(func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
		m := mode.Get()
		return rpc.ModeResponse{Mode: m.String()}, nil
	}))

	l.Server.RegisterHandler(constants.MethodContentGet, handler.New(func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
		content, err := fuzzyfinder.GetContent(ctx)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get content")
		}
		return rpc.ContentResponse{Content: content}, nil
	}))

	l.Server.RegisterHandler(constants.MethodLauncherOpenIn, handler.New(func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
		var params rpc.OpenInParams
		if err := req.UnmarshalParams(&params); err != nil {
			return nil, errors.WrapIf(err, "failed to unmarshal parameters")
		}

		err := fuzzyfinder.OpenIn(ctx, params.Category, params.Path)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to open in")
		}

		return rpc.EmptyResponse{}, nil
	}))
}
