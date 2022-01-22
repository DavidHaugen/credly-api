package internal

type MarvelService interface {
	GetUsers() ([]User, error)
}

type CredlyService interface {
	GetBadges() error
}
