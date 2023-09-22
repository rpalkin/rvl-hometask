package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"hometask/internal/service"
	"hometask/internal/structures"
)

type BirthdayHandler struct {
	BirthdayService service.BirthdayServiceInterface
}

func NewBirthdayHandler(birthdayService service.BirthdayServiceInterface) *BirthdayHandler {
	return &BirthdayHandler{
		BirthdayService: birthdayService,
	}
}

func (b *BirthdayHandler) PutBirthday(ctx *fiber.Ctx) error {
	var birthdayRequest structures.PutBirthdayRequest
	if err := ctx.BodyParser(&birthdayRequest); err != nil {
		zap.L().Error("failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(structures.ErrorResponse{
			CommonResponse: structures.CommonResponse{Message: "failed to parse request body"},
		})
	}

	if err := birthdayRequest.Validate(); err != nil {
		zap.L().Error("failed to validate request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(structures.ErrorResponse{
			CommonResponse: structures.CommonResponse{Message: err.Error()},
		})
	}

	if err := b.BirthdayService.PutBirthday(ctx.Params("username"), birthdayRequest.DateOfBirth.Time); err != nil {
		zap.L().Error("failed to validate request body", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(structures.ErrorResponse{
			CommonResponse: structures.CommonResponse{Message: "failed to store birthday"},
		})
	}

	return ctx.Status(fiber.StatusNoContent).Send(nil)
}

func (b *BirthdayHandler) GetGreetings(ctx *fiber.Ctx) error {
	days, err := b.BirthdayService.GetDaysToBirthday(ctx.Params("username"))
	if err != nil {
		zap.L().Error("failed to get days to birthday", zap.Error(err))
		if errors.Is(err, service.UserNotFoundError) {
			return ctx.Status(fiber.StatusNotFound).JSON(structures.ErrorResponse{
				CommonResponse: structures.CommonResponse{Message: "Sorry, user not found"},
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(structures.ErrorResponse{
			CommonResponse: structures.CommonResponse{Message: "failed to get days to birthday"},
		})
	}

	if days < 1 {
		return ctx.Status(fiber.StatusOK).JSON(structures.SuccessResponse{
			CommonResponse: structures.CommonResponse{Message: fmt.Sprintf("Hello, %s! Happy birthday!", ctx.Params("username"))},
		})
	}

	suffix := ""
	if days > 1 {
		suffix = "s"
	}
	return ctx.Status(fiber.StatusOK).JSON(structures.SuccessResponse{
		CommonResponse: structures.CommonResponse{Message: fmt.Sprintf("Hello, %s! Your birthday is in %d day%s", ctx.Params("username"), days, suffix)},
	})
}
