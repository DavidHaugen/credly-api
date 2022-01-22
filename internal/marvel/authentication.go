package marvel

import (
	"crypto/md5"
	"fmt"

	"github.com/google/uuid"
)

func (s Service) GetAuthParams() string {
	uuid := uuid.New()
	str := fmt.Sprintf(`%s%s%s`, uuid, s.PrivateAPIKey, s.PublicAPIKey)
	hash := md5.Sum([]byte(str))

	return fmt.Sprintf(`ts=%s&apikey=%s&hash=%s`, uuid, s.PublicAPIKey, hash)
}
