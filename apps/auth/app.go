package auth

import (
	"github.com/vertex-center/vertex/apps/auth/adapter"
	"github.com/vertex-center/vertex/apps/auth/core/service"
	"github.com/vertex-center/vertex/apps/auth/database"
	"github.com/vertex-center/vertex/apps/auth/handler"
	"github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/core/types/storage"
	"github.com/vertex-center/vertex/pkg/router"
)

// docapi:auth title Vertex Auth
// docapi:auth description An authentication service for Vertex.
// docapi:auth version 0.0.0
// docapi:auth filename auth

// docapi:auth url http://{ip}:{port-kernel}/api
// docapi:auth urlvar ip localhost The IP address of the server.
// docapi:auth urlvar port-kernel 7502 The port of the server.

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
	return meta.Meta
}

func (a *App) Initialize(r *router.Group) error {
	r.Use(middleware.ReadAuth)

	db, err := storage.NewDB(storage.DBParams{
		ID:         meta.Meta.ID,
		SchemaFunc: database.GetSchema,
		Migrations: database.Migrations,
	})
	if err != nil {
		return err
	}

	authAdapter := adapter.NewAuthDbAdapter(db)
	emailAdapter := adapter.NewEmailDbAdapter(db)

	authService := service.NewAuthService(authAdapter)
	emailService := service.NewEmailService(emailAdapter)
	userService := service.NewUserService(authAdapter)

	authHandler := handler.NewAuthHandler(authService)
	// docapi:auth route /app/auth/login auth_login
	r.POST("/login", authHandler.Login)
	// docapi:auth route /app/auth/register auth_register
	r.POST("/register", authHandler.Register)
	// docapi:auth route /app/auth/logout auth_logout
	r.POST("/logout", middleware.Authenticated, authHandler.Logout)
	// docapi:auth route /app/auth/verify auth_verify
	r.POST("/verify", authHandler.Verify)

	userHandler := handler.NewUserHandler(userService)
	user := r.Group("/user", middleware.Authenticated)
	// docapi:auth route /app/auth/user auth_get_current_user
	user.GET("", userHandler.GetCurrentUser)
	// docapi:auth route /app/auth/user auth_patch_current_user
	user.PATCH("", userHandler.PatchCurrentUser)
	// docapi:auth route /app/auth/credentials auth_get_current_user_credentials
	user.GET("/credentials", userHandler.GetCurrentUserCredentials)

	emailHandler := handler.NewEmailHandler(emailService)
	email := user.Group("/email", middleware.Authenticated)
	// docapi:auth route /app/auth/user/email auth_current_user_create_email
	email.POST("", emailHandler.CreateCurrentUserEmail)
	// docapi:auth route /app/auth/user/email auth_current_user_delete_email
	email.DELETE("", emailHandler.DeleteCurrentUserEmail)

	emails := user.Group("/emails", middleware.Authenticated)
	// docapi:auth route /app/auth/user/emails auth_current_user_get_emails
	emails.GET("", emailHandler.GetCurrentUserEmails)

	return nil
}
