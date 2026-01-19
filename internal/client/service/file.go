package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/repository"
	"gophkeeper/internal/model"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type FileService interface {
	AddFile(filePath, masterPass string) error
	DownloadFile(fileName string) error
	ListFile() ([]model.FileAPI, error)
	InitToken(TokenService)
}

type fileService struct {
	EncryptionService EncryptionService
	SenderService     SenderService
	TokenService      TokenService
	FileRepository    repository.FileRepository
}

func NewFileService(f repository.FileRepository, s SenderService, e EncryptionService) *fileService {
	return &fileService{
		EncryptionService: e,
		SenderService:     s,
		FileRepository:    f,
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

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	totalSize := stat.Size()

	// Прогресс-бар
	bar := pb.Full.Start64(totalSize).SetWidth(80)
	bar.Set("filename", filepath.Base(filePath))

	// Обёртка с прогрессом
	reader := bar.NewProxyReader(file)

	// TODO продумать контекст
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendFile(ctx, "/api/files/add", reader, fileName)

	if err != nil {
		return err
	}

	bar.Finish()

	fmt.Println(string(body))
	fmt.Println("code", response.StatusCode)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error add file, %v", string(body))
	}

	return nil
}

func (s *fileService) DownloadFile(fileName string) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	// TODO продумать контекст
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	response, err := s.SenderService.SendGetFile(ctx, "/api/files/download/"+fileName)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	fmt.Println("/api/files/download/" + fileName)
	fmt.Println("status", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error status download files")
	}

	err = s.FileRepository.DownloadFile(ctx, response)
	if err != nil {
		return err
	}
	return nil
}

func (s *fileService) ListFile() ([]model.FileAPI, error) {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return nil, err
	}
	s.SenderService.SetToken(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendGet(ctx, "/api/files/list")

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status get list files, %v", string(body))
	}

	var files []model.FileAPI
	if err = json.Unmarshal(body, &files); err != nil {
		return nil, err
	}

	return files, nil
}
