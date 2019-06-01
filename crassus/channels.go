package crassus

type Channel struct {
	ID            string
	Channel       chan Message
	Subscriptions []Client
}

func (c *Channel) Subscribe(client Client) {
	c.Subscriptions = append(c.Subscriptions, client) // this is the worst possible way to do it
}
