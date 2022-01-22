package credly

// Service :
type Service struct {
	OrganizationID string
	AuthToken      string
}

// NewService :
func NewService(organizationID, authToken string) *Service {
	return &Service{organizationID, authToken}
}
