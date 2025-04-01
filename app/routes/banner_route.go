package routes

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
)

func BannerRoute(route fiber.Router, controller *controllers.Controller) {
	bannerRoutes := route.Group("/banners", middleware.AuthProtected()...)
	bannerRoutes.Get("/", controller.BannerController.ListBanners)
	bannerRoutes.Get("/:id", controller.BannerController.GetBanner)
}
