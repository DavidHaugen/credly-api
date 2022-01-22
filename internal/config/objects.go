package config

type Config struct {
	Marvel Marvel
}

type Marvel struct {
	PublicAPIKey  string
	PrivateAPIKey string
}
