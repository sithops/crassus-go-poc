package crassus

import (
	"fmt"
)

type Router struct{}

func (r *Router) Run() {
	fmt.Println("Crassus")
}

func NewRouter() *Router {
	return new(Router)
}
