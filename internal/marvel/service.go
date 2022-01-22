package marvel

// Service :
type Service struct {
	PublicAPIKey  string
	PrivateAPIKey string
}

// NewService :
func NewService(publicKey, privateKey string) *Service {
	return &Service{publicKey, privateKey}
}
