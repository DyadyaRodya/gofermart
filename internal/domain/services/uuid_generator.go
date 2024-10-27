package services

import uuidPkg "github.com/google/uuid"

type UUID4Generator struct{}

func NewUUID4Generator() *UUID4Generator {
	return &UUID4Generator{}
}

func (u *UUID4Generator) Generate() (string, error) {
	uuid, err := uuidPkg.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
