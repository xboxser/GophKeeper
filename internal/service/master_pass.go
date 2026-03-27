package service

import "golang.org/x/crypto/bcrypt"

// ConvertMasterPassToHash - преобразует мастер пароль в хэш
func ConvertMasterPassToHash(masterPass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(masterPass), bcrypt.DefaultCost)
}

func ValidateHash(hash []byte, masterPass string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(masterPass))
}
