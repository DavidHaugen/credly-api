package marvel

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/DavidHaugen/golang-boilerplate/internal"
)

func (s Service) GetUsers() ([]internal.User, error) {
	baseURL := "https://gateway.marvel.com/v1/public/characters"
	res, err := http.Get(s.getRequestURL(baseURL))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	userResponse := &GetUserResponse{}
	err = json.Unmarshal(body, userResponse)
	switch {
	case err != nil:
		return nil, err
	case userResponse == nil:
		return nil, errors.New("user response is nil")
	}

	return mapUsers(*userResponse), nil
}

func mapUsers(userResponse GetUserResponse) []internal.User {
	users := []internal.User{}
	return users
}
