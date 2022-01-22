package credly

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (s Service) GetBadges() error {
	client := &http.Client{}
	req, err := s.getBadgesRequest()
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

// If this API were serving multiple orgs, would change to accept org ID as argument
func (s Service) getBadgesRequest() (*http.Request, error) {
	url := fmt.Sprintf(`https://sandbox-api.credly.com/v1/organizations/%s/badges`, s.OrganizationID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.AuthToken, "")
	return req, nil
}
