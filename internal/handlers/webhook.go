package handlers

import (
	"content-flow/internal/pkgs/apierrors"
	"content-flow/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CreateWebhookRequest struct {
	URL    string `json:"url"`
	Events string `json:"events"`
}

// CreateWebhook godoc
// @Summary Create a webhook
// @Description Register a new webhook for events
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param webhook body CreateWebhookRequest true "Webhook Config"
// @Success 200 {object} models.Webhook
// @Failure 400 {object} apierrors.AppError
// @Security Bearer
// @Router /api/webhooks [post]
func CreateWebhook(c *fiber.Ctx) error {
	req := new(CreateWebhookRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON")
	}

	webhook, err := services.CreateWebhook(req.URL, req.Events)
	if err != nil {
		return apierrors.Internal("Failed to create webhook: " + err.Error())
	}

	return c.JSON(webhook)
}

// GetAllWebhooks godoc
// @Summary List webhooks
// @Tags Webhooks
// @Produce json
// @Success 200 {array} models.Webhook
// @Security Bearer
// @Router /api/webhooks [get]
func GetAllWebhooks(c *fiber.Ctx) error {
	webhooks, err := services.GetAllWebhooks()
	if err != nil {
		return apierrors.Internal(err.Error())
	}
	return c.JSON(webhooks)
}
