package dto

import interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"

type CredentialsDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (c CredentialsDTO) ConvertToInteractorsDTO() *interactorsdto.Credentials {
	return &interactorsdto.Credentials{
		Login:    c.Login,
		Password: c.Password,
	}
}
