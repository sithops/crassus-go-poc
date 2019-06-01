package crassus

const (
	SUBSCRIBE = "SUBSCRIBE"
	TELL      = "TELL"
	PART      = "UNSUBSCRIBE"
)

func NewRouter() *Router {
	return new(Router)
}
