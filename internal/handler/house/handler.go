package househandler

import "net/http"

type Handler interface {
	Create() http.HandlerFunc
	Houses() http.HandlerFunc
}
