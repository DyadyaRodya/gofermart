package interfaces

type PasswordService interface {
	Hash(password, salt string) (string, error)
	Validate(password string) bool
	GenerateRandomSalt(saltSize int) (string, error)
	Compare(currPassword, hashedPassword, salt string) (bool, error)
}
