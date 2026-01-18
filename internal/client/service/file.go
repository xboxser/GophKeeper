package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileService interface {
	AddFile(filePath, masterPass string) error
	InitToken(TokenService)
}

type fileService struct {
	EncryptionService EncryptionService
	SenderService     SenderService
	TokenService      TokenService
}

func NewFileService(senderService SenderService, encryptionService EncryptionService) *fileService {
	return &fileService{
		EncryptionService: encryptionService,
		SenderService:     senderService,
	}
}

func (s *fileService) InitToken(tokenService TokenService) {
	s.TokenService = tokenService
}

func (s *fileService) AddFile(filePath, masterPass string) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	s.EncryptionService.SetMasterPass(masterPass)

	fileName := filepath.Base(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendFile(ctx, "/api/files/add", file, fileName)

	if err != nil {
		return err
	}

	fmt.Println(string(body))
	fmt.Println("code", response.StatusCode)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error add file, %v", string(body))
	}

	return nil
}
