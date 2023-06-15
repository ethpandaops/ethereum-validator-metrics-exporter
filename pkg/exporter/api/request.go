package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (c *client) get(ctx context.Context, path, url string) ([]byte, error) {
	start := time.Now()

	httpMethod := "GET"

	c.metrics.ObserveRequest(httpMethod, path)

	var rsp *http.Response

	var err error

	defer func() {
		rspCode := "none"
		if rsp != nil {
			rspCode = fmt.Sprintf("%d", rsp.StatusCode)
		}

		c.metrics.ObserveResponse(httpMethod, path, rspCode, time.Since(start))
	}()

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)

	if err != nil {
		return nil, err
	}

	if c.apikey != "" {
		req.Header.Set("apikey", c.apikey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *client) getValidators(ctx context.Context, pubkeys []string) (*Response[[]Validator], error) {
	path := "api/v1/validator"
	url := fmt.Sprintf("%s/%s/%s", c.url, path, strings.Join(pubkeys, ","))

	data, err := c.get(ctx, path, url)
	if err != nil {
		return nil, err
	}

	resp := new(Response[[]Validator])
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) getValidator(ctx context.Context, pubkey string) (*Response[Validator], error) {
	path := "api/v1/validator"
	url := fmt.Sprintf("%s/%s/%s", c.url, path, pubkey)

	data, err := c.get(ctx, path, url)
	if err != nil {
		return nil, err
	}

	resp := new(Response[Validator])
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
