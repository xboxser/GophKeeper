package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type SenderService interface {
	SendPost(context.Context, string, []byte) ([]byte, *http.Response, error)
}

type senderService struct {
	client        *http.Client
	serverAddress string
}

func NewSenderService(client *http.Client, serverAddress string) *senderService {
	return &senderService{
		client:        client,
		serverAddress: serverAddress,
	}
}

func (s *senderService) SendPost(ctx context.Context, url string, json []byte) ([]byte, *http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.serverAddress+url, bytes.NewBuffer(json))
	if err != nil {
		return nil, nil, err
	}

	response, err := s.client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, response, nil
}
