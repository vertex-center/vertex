package auth

import (
	"github.com/vertex-center/vertex/apps/auth/adapter"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/service"
	"github.com/vertex-center/vertex/apps/auth/handler"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

var (
	authAdapter port.AuthAdapter

	AuthService port.AuthService
	userService port.UserService
)

type App struct {
	ctx *apptypes.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx
}

func (a *App) Meta() apptypes.Meta {
	return apptypes.Meta{
		ID:          "auth",
		Name:        "Vertex Auth",
		Description: "Authentication app for Vertex",
		Icon:        "admin_panel_settings",
		Hidden:      true,
		DefaultPort: "7502",
	}
}

func (a *App) Initialize(r *router.Group) error {
	authAdapter = adapter.NewAuthDbAdapter(a.ctx.Db())
	emailAdapter := adapter.NewEmailDbAdapter(a.ctx.Db())

	AuthService = service.NewAuthService(authAdapter)
	emailService := service.NewEmailService(emailAdapter)
	userService = service.NewUserService(authAdapter)

	middleware.AuthService = AuthService

	authHandler := handler.NewAuthHandler(AuthService)
	// docapi:v route /app/auth/login auth_login
	r.POST("/login", authHandler.Login)
	// docapi:v route /app/auth/register auth_register
	r.POST("/register", authHandler.Register)
	// docapi:v route /app/auth/logout auth_logout
	r.POST("/logout", middleware.Authenticated, authHandler.Logout)

	userHandler := handler.NewUserHandler(userService)
	user := r.Group("/user")
	// docapi:v route /app/auth/user auth_get_current_user
	user.GET("", middleware.Authenticated, userHandler.GetCurrentUser)
	// docapi:v route /app/auth/user auth_patch_current_user
	user.PATCH("", middleware.Authenticated, userHandler.PatchCurrentUser)
	// docapi:v route /app/auth/credentials auth_get_current_user_credentials
	user.GET("/credentials", middleware.Authenticated, userHandler.GetCurrentUserCredentials)

	emailHandler := handler.NewEmailHandler(emailService)
	email := user.Group("/email")
	// docapi:v route /app/auth/user/email auth_current_user_create_email
	email.POST("", middleware.Authenticated, emailHandler.CreateCurrentUserEmail)
	// docapi:v route /app/auth/user/email auth_current_user_delete_email
	email.DELETE("", middleware.Authenticated, emailHandler.DeleteCurrentUserEmail)

	emails := user.Group("/emails")
	// docapi:v route /app/auth/user/emails auth_current_user_get_emails
	emails.GET("", middleware.Authenticated, emailHandler.GetCurrentUserEmails)

	return nil
}
