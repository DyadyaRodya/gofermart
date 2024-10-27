package interfaces

type LoginService interface {
	Validate(login string) error
}
