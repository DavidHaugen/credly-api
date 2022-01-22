package marvel

import (
	"encoding/json"
	"errors"
	"fmt"
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

	fmt.Println("response: ", *userResponse)

	return mapUsers(*userResponse), nil
}

func mapUsers(userResponse GetUserResponse) []internal.User {
	users := []internal.User{}
	for _, v := range userResponse.Data.Results {
		users = append(users, internal.User{
			MarvelID:    v.ID,
			Name:        v.Name,
			Description: v.Description,
			Thumbnail:   fmt.Sprintf(`%s%s`, v.Thumbnail.Path, v.Thumbnail.Extension),
		})
	}
	return users
}
