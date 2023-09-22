package web

import (
	"github.com/gofiber/fiber/v2"
	"hometask/internal/handler"
)

type Router struct {
	BirthdayHandler *handler.BirthdayHandler
}

func NewRouter(birthdayHandler *handler.BirthdayHandler) *Router {
	return &Router{
		BirthdayHandler: birthdayHandler,
	}
}

func (r *Router) Route(router fiber.Router) {
	router.Put("/:username<alpha>", r.BirthdayHandler.PutBirthday)
	router.Get("/:username<alpha>", r.BirthdayHandler.GetGreetings)
}
