package command

import (
	"fmt"
	"gophkeeper/internal/client/service"

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
		Use:     "addFile",
		Short:   "Добавить новый файл",
		Example: `  todo addFile -f=name.txt -m=secret`,
		Run:     s.addFileRun,
	}

	addFileCmd.Flags().StringP("file", "f", "", "Пароль для шифрования (Обязательный)")
	addFileCmd.MarkFlagRequired("file")
	addFileCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	addFileCmd.MarkFlagRequired("masterPass")

	return []*cobra.Command{
		addFileCmd,
	}
}

func (s *FileHandler) addFileRun(cmd *cobra.Command, args []string) {

	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}

	filePath, err := cmd.Flags().GetString("file")
	if err != nil || filePath == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--file или -f)")
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
