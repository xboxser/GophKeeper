package command

import (
	"fmt"
	"gophkeeper/internal/client/service"
	commonService "gophkeeper/internal/service"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/cobra"
)

type FileHandler struct {
	FileService       service.FileService
	UserMasterService service.UserMasterService
}

func NewFileHandler(s service.FileService, u service.UserMasterService) *FileHandler {
	return &FileHandler{
		FileService:       s,
		UserMasterService: u,
	}
}

func (s *FileHandler) GetCommands() []*cobra.Command {
	addFileCmd := &cobra.Command{
		Use:     "file-add",
		Short:   "Добавить новый файл",
		Example: `  todo file-add -f=name.txt -m=secret`,
		Run:     s.addFileRun,
	}

	addFileCmd.Flags().StringP("file", "f", "", "путь к файлу (обязательный)")
	addFileCmd.MarkFlagRequired("file")
	addFileCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	addFileCmd.MarkFlagRequired("masterPass")

	listFileCmd := &cobra.Command{
		Use:     "file-list",
		Short:   "Получить список файлов",
		Example: `  todo file-list`,
		Run:     s.listFileRun,
	}

	downloadFileCmd := &cobra.Command{
		Use:     "file-get",
		Short:   "Скачать файл с сервера",
		Example: `  todo file-get -f=name.txt`,
		Run:     s.downloadFileRun,
	}

	downloadFileCmd.Flags().StringP("file", "f", "", "наименование файла для скачивания с сервера (Обязательный)")
	downloadFileCmd.MarkFlagRequired("file")

	return []*cobra.Command{
		addFileCmd,
		listFileCmd,
		downloadFileCmd,
	}
}

func (s *FileHandler) addFileRun(cmd *cobra.Command, _ []string) {

	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}

	filePath, err := cmd.Flags().GetString("file")
	if err != nil || filePath == "" {
		fmt.Println("❌ Ошибка: путь до файла выгрузки (--file или -f)")
		return
	}

	err = s.UserMasterService.Master(masterPass)
	if err != nil {
		fmt.Println("❌ Ошибка проверки мастер пароля:", err)
		return
	}

	err = s.FileService.AddFile(filePath, masterPass)

	if err != nil {
		fmt.Println("❌ Ошибка загрузки файла:", err)
		return
	}

}

func (s *FileHandler) listFileRun(_ *cobra.Command, _ []string) {

	files, err := s.FileService.ListFile()

	if err != nil {
		fmt.Println("❌ Ошибка получение списка файлов:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("📦 Список файлов пуст")
		return
	}

	fmt.Println("📁 Список файлов:")

	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "Имя файла", "Размер"})

	for i, file := range files {
		t.AppendRow(table.Row{i + 1, file.Name, commonService.FormatFileSize(file.Size)})
	}

	t.SetStyle(table.StyleLight)
	fmt.Println(t.Render())
}

func (s *FileHandler) downloadFileRun(cmd *cobra.Command, _ []string) {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil || filePath == "" {
		fmt.Println("❌ Ошибка: укажите наименование файла для скачивания (--file или -f)")
		return
	}

	err = s.FileService.DownloadFile(filePath)

	if err != nil {
		fmt.Println("❌ Ошибка скачивания файла:", err)
		return
	}

	fmt.Println("✅ Файл", filePath, "скачан")
}
