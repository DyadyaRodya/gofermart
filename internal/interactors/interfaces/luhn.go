package interfaces

type LuhnService interface {
	Validate(number string) bool
}
