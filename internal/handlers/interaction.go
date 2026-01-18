package handlers

import (
	"content-flow/internal/pkgs/apierrors"
	"content-flow/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CreateCommentRequest struct {
	Body string `json:"body"`
}

// AddComment godoc
// @Summary Add a comment
// @Description Add a comment to a content
// @Tags Engagement
// @Accept json
// @Produce json
// @Param id path int true "Content ID"
// @Param request body CreateCommentRequest true "Comment Body"
// @Success 200 {object} models.Comment
// @Failure 400 {object} apierrors.AppError
// @Security Bearer
// @Router /api/content/{id}/comments [post]
func AddComment(c *fiber.Ctx) error {
	contentID, _ := strconv.Atoi(c.Params("id"))
	userID := uint(c.Locals("user_id").(float64))

	req := new(CreateCommentRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON")
	}

	if req.Body == "" {
		return apierrors.BadRequest("Comment body cannot be empty")
	}

	comment, err := services.AddComment(userID, uint(contentID), req.Body)
	if err != nil {
		return apierrors.Internal("Failed to add comment: " + err.Error())
	}

	return c.JSON(comment)
}

// GetComments godoc
// @Summary Get comments
// @Description Get all comments for a content
// @Tags Engagement
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {array} models.Comment
// @Router /api/content/{id}/comments [get]
func GetComments(c *fiber.Ctx) error {
	contentID, _ := strconv.Atoi(c.Params("id"))
	comments, err := services.GetComments(uint(contentID))
	if err != nil {
		return apierrors.Internal("Failed to fetch comments")
	}
	return c.JSON(comments)
}

type LikeResponse struct {
	Liked      bool  `json:"liked"`
	TotalLikes int64 `json:"total_likes"`
}

// ToggleLike godoc
// @Summary Like/Unlike content
// @Description Toggle like status for content
// @Tags Engagement
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {object} LikeResponse
// @Security Bearer
// @Router /api/content/{id}/like [post]
func ToggleLike(c *fiber.Ctx) error {
	contentID, _ := strconv.Atoi(c.Params("id"))
	userID := uint(c.Locals("user_id").(float64))

	liked, err := services.ToggleLike(userID, uint(contentID))
	if err != nil {
		return apierrors.Internal("Failed to update like")
	}

	count, _ := services.GetLikeCount(uint(contentID))

	return c.JSON(LikeResponse{
		Liked:      liked,
		TotalLikes: count,
	})
}
