package sse

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"
)

func SubscribeToSse(
	ctx context.Context,
	url string,
	onData func(string) bool,
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// Essential header for SSE
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	// Scanner reads the body line by line
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// SSE protocol: data is prefixed with "data: "
		if after, ok := strings.CutPrefix(line, "data: "); ok {
			data := after
			if !onData(data) {
				break
			}
		}

		// An empty line in SSE signals the end of a single event block
		if line == "" {
			// Logic for handling end of event block goes here
		}
	}

	return scanner.Err()
}
