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
	"github.com/wI2L/fizz"
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
	return meta.Meta
}

func (a *App) Initialize(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth())

	db, err := storage.NewDB(storage.DBParams{
		ID:         meta.Meta.ID,
		SchemaFunc: database.GetSchema,
		Migrations: database.Migrations,
	})
	if err != nil {
		return err
	}

	var (
		authAdapter  = adapter.NewAuthDbAdapter(db)
		emailAdapter = adapter.NewEmailDbAdapter(db)

		authService  = service.NewAuthService(authAdapter)
		emailService = service.NewEmailService(emailAdapter)
		userService  = service.NewUserService(authAdapter)

		authHandler  = handler.NewAuthHandler(authService)
		userHandler  = handler.NewUserHandler(userService)
		emailHandler = handler.NewEmailHandler(emailService)

		user   = r.Group("/user", "User", "", middleware.Authenticated())
		email  = user.Group("/email", "Email", "User emails", middleware.Authenticated())
		emails = user.Group("/emails", "Emails", "User emails", middleware.Authenticated())
	)

	r.POST("/login", authHandler.LoginInfo(), authHandler.Login())
	r.POST("/register", authHandler.RegisterInfo(), authHandler.Register())
	r.POST("/logout", authHandler.LogoutInfo(), middleware.Authenticated(), authHandler.Logout())
	r.POST("/verify", authHandler.VerifyInfo(), authHandler.Verify())

	user.GET("", userHandler.GetCurrentUserInfo(), userHandler.GetCurrentUser())
	user.PATCH("", userHandler.PatchCurrentUserInfo(), userHandler.PatchCurrentUser())
	user.GET("/credentials", userHandler.GetCurrentUserCredentialsInfo(), userHandler.GetCurrentUserCredentials())

	email.POST("", emailHandler.CreateCurrentUserEmailInfo(), emailHandler.CreateCurrentUserEmail())
	email.DELETE("", emailHandler.DeleteCurrentUserEmailInfo(), emailHandler.DeleteCurrentUserEmail())

	emails.GET("", emailHandler.GetCurrentUserEmailsInfo(), emailHandler.GetCurrentUserEmails())

	return nil
}
