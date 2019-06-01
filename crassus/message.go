package crassus

type Message struct {
	Channel string   `json:"channel"`
	From    string   `json:from"`
	Action  string   `json:"command"`
	Args    []string `json:"args"`
}
