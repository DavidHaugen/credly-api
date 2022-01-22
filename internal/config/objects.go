package config

type Config struct {
	Marvel Marvel
	Credly Credly
}

type Marvel struct {
	PublicAPIKey  string
	PrivateAPIKey string
}

type Credly struct {
	OrganizationID string
	AuthToken      string
}
