package relay

type Error struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
