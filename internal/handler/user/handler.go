package userhandler

import "net/http"

type Handler interface {
	Registration() http.HandlerFunc
	Login() http.HandlerFunc
}
