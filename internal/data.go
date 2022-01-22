package internal

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   struct {
		Path      string `json:"path"`
		Extension string `json:"extension"`
	} `json:"thumbnail"`
}
