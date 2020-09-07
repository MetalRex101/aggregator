package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const maxOutRequests = 4
const urlsLimit = 20

type scriptError struct {
	Code    int
	Message string
}

type Script interface {
	Run(ctx context.Context, body []byte) (ResponseBody, *scriptError)
}

type ResponseBody map[string]string

func NewScript(client *http.Client) *script {
	return &script{client: client}
}

type script struct {
	client *http.Client
}

type stopCh chan struct{}

func (c stopCh) Close() {
	select {
	case <-c:
	default:
		close(c)
	}
}

type urlResponse struct {
	url  string
	body []byte
}

func (s *script) Run(ctx context.Context, body []byte) (ResponseBody, *scriptError) {
	var urls Urls
	if err := json.Unmarshal(body, &urls); err != nil {
		return nil, &scriptError{Code: http.StatusInternalServerError, Message: "failed to unmarshal request body"}
	}

	if len(urls) > urlsLimit {
		return nil, &scriptError{Code: http.StatusBadRequest, Message: "maximum 20 urls per request"}
	}

	stopCh := make(stopCh)
	urlResponseCh := make(chan urlResponse, len(urls))

	if requestErr := s.requestUrls(ctx, urls, stopCh, urlResponseCh); requestErr != nil {
		return nil, &scriptError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to perform request to one of the urls: %s", requestErr),
		}
	}

	resp := make(ResponseBody, len(urls))
	for urlResponse := range urlResponseCh {
		resp[urlResponse.url] = string(urlResponse.body)
	}

	return resp, nil
}

func (s *script) requestUrls(ctx context.Context, urls []string, stopCh stopCh, respCh chan urlResponse) (requestErr error) {
	p := NewPool(maxOutRequests)

	defer func() {
		if requestErr != nil {
			stopCh.Close()
		}
		close(respCh)
	}()

	for _, url := range urls {
		select {
		case <-stopCh:
			break
		case <-ctx.Done():
			break
		default:
			p.Run(func() {
				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					requestErr = fmt.Errorf("failed to create request %w", err)
					return
				}

				req.Header.Set("content-type", "application/json")

				resp, err := s.client.Do(req)
				if err != nil {
					requestErr = err
					return
				}

				if resp.StatusCode > 299 {
					requestErr = fmt.Errorf("%s returned non 2** http Code: %d", url, resp.StatusCode)
					return
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					requestErr = err
					return
				}

				respCh <- urlResponse{url: url, body: body}
			})
		}
	}

	return
}
