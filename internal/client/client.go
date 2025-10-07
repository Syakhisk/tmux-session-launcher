package client

import (
	"context"
	"net"
	"strings"
	"tmux-session-launcher/internal/rpc"

	"emperror.dev/errors"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

type Client struct {
	address string
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

func (c *Client) Call(ctx context.Context, method string, params, result interface{}) error {
	conn, err := net.Dial("unix", c.address)
	if err != nil {
		return errors.Wrap(err, "failed to connect to server")
	}
	defer conn.Close()

	// Create channel for the connection
	ch := channel.Line(conn, conn)

	// Create jrpc2 client
	cli := jrpc2.NewClient(ch, nil)
	defer cli.Close()

	// Make the RPC call
	_, err = cli.Call(ctx, method, params)
	if err != nil {
		return errors.Wrap(err, "RPC call failed")
	}

	return nil
}

func (c *Client) CallWithResult(ctx context.Context, method string, params, result interface{}) error {
	conn, err := net.Dial("unix", c.address)
	if err != nil {
		return errors.Wrap(err, "failed to connect to server")
	}
	defer conn.Close()

	// Create channel for the connection
	ch := channel.Line(conn, conn)

	// Create jrpc2 client
	cli := jrpc2.NewClient(ch, nil)
	defer cli.Close()

	// Make the RPC call
	response, err := cli.Call(ctx, method, params)
	if err != nil {
		return errors.Wrap(err, "RPC call failed")
	}

	// Unmarshal result if provided
	if result != nil {
		err = response.UnmarshalResult(result)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal result")
		}
	}

	return nil
}

// Convenience methods for specific RPC calls

func (c *Client) NextMode(ctx context.Context) (*rpc.ModeResponse, error) {
	var result rpc.ModeResponse
	err := c.CallWithResult(ctx, rpc.MethodModeNext, rpc.EmptyParams{}, &result)
	return &result, err
}

func (c *Client) PrevMode(ctx context.Context) (*rpc.ModeResponse, error) {
	var result rpc.ModeResponse
	err := c.CallWithResult(ctx, rpc.MethodModePrev, rpc.EmptyParams{}, &result)
	return &result, err
}

func (c *Client) GetMode(ctx context.Context) (*rpc.ModeResponse, error) {
	var result rpc.ModeResponse
	err := c.CallWithResult(ctx, rpc.MethodModeGet, rpc.EmptyParams{}, &result)
	return &result, err
}

func (c *Client) GetContent(ctx context.Context) (*rpc.ContentResponse, error) {
	var result rpc.ContentResponse
	err := c.CallWithResult(ctx, rpc.MethodContentGet, rpc.EmptyParams{}, &result)
	return &result, err
}

func (c *Client) OpenIn(ctx context.Context, selection string) error {
	split := strings.Split(selection, "|")
	if len(split) != 2 {
		return errors.Errorf("invalid selection format: %s", selection)
	}

	params := rpc.OpenInParams{
		Category: strings.TrimSpace(split[0]),
		Path:     strings.TrimSpace(split[1]),
	}

	return c.Call(ctx, rpc.MethodLauncherOpenIn, params, nil)
}
