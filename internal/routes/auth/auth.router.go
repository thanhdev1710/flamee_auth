package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhdev1710/flamee_auth/internal/controllers"
	"github.com/thanhdev1710/flamee_auth/middlewares"
)

type AuthRouter struct {
}

func (ur *AuthRouter) InitAuthRouter(Router *gin.RouterGroup) {
	// public router
	AuthRouterPublic := Router.Group("/auth")
	{
		AuthRouterPublic.POST("/register",
			middlewares.RateLimitPerRouteAndIP(1, 30*time.Second, 3),
			controllers.NewAuthControllers().Register,
		)

		AuthRouterPublic.POST("/login",
			middlewares.RateLimitPerRouteAndIP(1, 10*time.Second, 5),
			controllers.NewAuthControllers().Login,
		)

		AuthRouterPublic.POST("/refresh-token",
			middlewares.RateLimitPerRouteAndIP(1, 5*time.Second, 10),
			controllers.NewAuthControllers().RefreshToken,
		)

		AuthRouterPublic.POST("/logout",
			middlewares.RateLimitPerRouteAndIP(1, 5*time.Second, 5),
			controllers.NewAuthControllers().Logout,
		)

		AuthRouterPublic.POST("/send-email/:email",
			middlewares.RateLimitPerRouteAndIP(1, 60*time.Second, 2),
			controllers.NewAuthControllers().SendVerifyEmail,
		)

		AuthRouterPublic.GET("/verify-email/:token",
			middlewares.RateLimitPerRouteAndIP(1, 60*time.Second, 2),
			controllers.NewAuthControllers().VerifyEmail,
		)

		AuthRouterPublic.POST("/reset-password/:email",
			middlewares.RateLimitPerRouteAndIP(1, 60*time.Second, 2),
			controllers.NewAuthControllers().SendResetPassword,
		)

		AuthRouterPublic.POST("/change-password/:token",
			middlewares.RateLimitPerRouteAndIP(1, 20*time.Second, 3),
			controllers.NewAuthControllers().ResetPassword,
		)

	}

	// private router
	AuthRouterPrivate := Router.Group("/auth").
		Use(middlewares.AuthMiddleware())
	{
		AuthRouterPrivate.POST("/delete-account",
			middlewares.RateLimitPerRouteAndIP(1, 20*time.Second, 3),
			middlewares.VerifyAccount(),
			controllers.NewAuthControllers().DeleteAccount,
		)

		AuthRouterPrivate.POST("/restore-account",
			middlewares.RateLimitPerRouteAndIP(1, 20*time.Second, 3),
			controllers.NewAuthControllers().RestoreAccount,
		)
	}
}
