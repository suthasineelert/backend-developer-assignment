package controllers

import (
	"backend-developer-assignment/app/services"

	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// BannerController handles HTTP requests for banner operations
type BannerController struct {
	bannerService services.BannerService
}

// NewBannerController creates a new banner controller
func NewBannerController(bannerService services.BannerService) *BannerController {
	return &BannerController{
		bannerService: bannerService,
	}
}

// GetBannerByID handles the request to get a banner by ID
func (c *BannerController) GetBanner(ctx *fiber.Ctx) error {
	bannerID := ctx.Params("id")
	if bannerID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "Banner ID is required")
	}

	banner, err := c.bannerService.GetBannerByID(bannerID)
	if err != nil {
		logger.Error("Failed to get banner", zap.String("banner_id", bannerID), zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get banner")
	}

	if banner == nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "Banner not found")
	}

	return ctx.Status(fiber.StatusOK).JSON(banner)
}

// ListBanners returns all banners for the current user
func (c *BannerController) ListBanners(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)

	// Get banners from service
	banners, err := c.bannerService.GetBannersByUserID(userID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(banners)
}
