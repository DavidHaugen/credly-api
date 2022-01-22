package internal

type MarvelService interface {
	GetUsers() ([]User, error)
}
