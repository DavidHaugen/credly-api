package internal

type User struct {
	MarvelID    int    `json:"marvelID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}
