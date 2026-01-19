package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type SenderService interface {
	SendPost(context.Context, string, []byte) ([]byte, *http.Response, error)
	SendGet(context.Context, string) ([]byte, *http.Response, error)
	SendGetFile(ctx context.Context, url string) (*http.Response, error)
	SendFile(context.Context, string, io.Reader, string) ([]byte, *http.Response, error)
	SetToken(string)
}

type senderService struct {
	client        *http.Client
	serverAddress string
	token         string
}

func NewSenderService(client *http.Client, serverAddress string) *senderService {
	return &senderService{
		client:        client,
		serverAddress: serverAddress,
		token:         "",
	}
}

// SetToken - устанавливаем токен
// При его наличии добавляться в запросы по умолчанию
func (s *senderService) SetToken(token string) {
	s.token = token
}

func (s *senderService) SendPost(ctx context.Context, url string, json []byte) ([]byte, *http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.serverAddress+url, bytes.NewBuffer(json))
	if err != nil {
		return nil, nil, err
	}

	if s.token != "" {
		request.Header.Set("Authorization", s.token)
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

func (s *senderService) SendGet(ctx context.Context, url string) ([]byte, *http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, s.serverAddress+url, nil)
	if err != nil {
		return nil, nil, err
	}

	if s.token != "" {
		request.Header.Set("Authorization", s.token)
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

func (s *senderService) SendFile(ctx context.Context, url string, b io.Reader, fileName string) ([]byte, *http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.serverAddress+url, b)
	if err != nil {
		return nil, nil, err
	}

	if s.token != "" {
		request.Header.Set("Authorization", s.token)
	}
	request.Header.Set("filename", fileName)

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

func (s *senderService) SendGetFile(ctx context.Context, url string) (*http.Response, error) {
	// TODO добавить скачивание Range
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, s.serverAddress+url, nil)
	if err != nil {
		return nil, err
	}

	if s.token != "" {
		request.Header.Set("Authorization", s.token)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil

}
