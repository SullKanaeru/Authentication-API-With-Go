package handler

import (
	"authentication_api/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserRepo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: repo}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memproses ID user",
		})
	}

	userID := uint(userIDFloat)

	user, err := h.UserRepo.FindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Data user tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data profil",
		"data":    user,
	})
}