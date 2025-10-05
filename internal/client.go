package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	address string
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

func (c *Client) Send(route, payload string) (string, error) {
	conn, err := net.Dial("unix", c.address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// send route
	fmt.Fprint(conn, route+"\n")

	// send payload (multi-line)
	if payload != "" {
		fmt.Fprintf(conn, "%s\n", payload)
	}

	// Send empty line to terminate
	fmt.Fprint(conn, "\n")

	// Read response
	scanner := bufio.NewScanner(conn)
	var payloadLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // empty line terminates response
			break
		}

		payloadLines = append(payloadLines, line)
	}

	if len(payloadLines) == 0 {
		return "", nil
	}

	response := strings.Join(payloadLines, "\n")
	if strings.HasPrefix(response, "ERROR:") {
		return "", fmt.Errorf("%s", strings.TrimPrefix(response, "ERROR: "))
	}

	return response, nil
}
