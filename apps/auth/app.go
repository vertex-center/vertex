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
	r.Use(middleware.ReadAuth)

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

		user   = r.Group("/user", "User", "", middleware.Authenticated)
		email  = user.Group("/email", "Email", "User emails", middleware.Authenticated)
		emails = user.Group("/emails", "Emails", "User emails", middleware.Authenticated)
	)

	// Auth

	r.POST("/login", []fizz.OperationOption{
		fizz.ID("login"),
		fizz.Summary("Login"),
		fizz.Description("Login with username and password"),
	}, authHandler.Login())

	r.POST("/register", []fizz.OperationOption{
		fizz.ID("register"),
		fizz.Summary("Register"),
		fizz.Description("Register a new user with username and password"),
	}, authHandler.Register())

	r.POST("/verify", []fizz.OperationOption{
		fizz.ID("verify"),
		fizz.Summary("Verify"),
		fizz.Description("Verify a token"),
	}, authHandler.Verify())

	r.POST("/logout", []fizz.OperationOption{
		fizz.ID("logout"),
		fizz.Summary("Logout"),
		fizz.Description("Logout a user"),
	}, middleware.Authenticated, authHandler.Logout())

	// User

	user.GET("", []fizz.OperationOption{
		fizz.ID("getCurrentUser"),
		fizz.Summary("Get user"),
		fizz.Description("Retrieve the logged-in user"),
	}, userHandler.GetCurrentUser())

	user.PATCH("", []fizz.OperationOption{
		fizz.ID("patchCurrentUser"),
		fizz.Summary("Patch user"),
		fizz.Description("Update the logged-in user"),
	}, userHandler.PatchCurrentUser())

	user.GET("/credentials", []fizz.OperationOption{
		fizz.ID("getCurrentUserCredentials"),
		fizz.Summary("Get user credentials"),
		fizz.Description("Retrieve the logged-in user credentials"),
	}, userHandler.GetCurrentUserCredentials())

	// Emails

	emails.GET("", []fizz.OperationOption{
		fizz.ID("getCurrentUserEmails"),
		fizz.Summary("Get emails"),
		fizz.Description("Retrieve the emails of the logged-in user"),
	}, emailHandler.GetCurrentUserEmails())

	email.POST("", []fizz.OperationOption{
		fizz.ID("createCurrentUserEmail"),
		fizz.Summary("Create email"),
		fizz.Description("Create a new email for the logged-in user"),
	}, emailHandler.CreateCurrentUserEmail())

	email.DELETE("", []fizz.OperationOption{
		fizz.ID("deleteCurrentUserEmail"),
		fizz.Summary("Delete email"),
		fizz.Description("Delete an email from the logged-in user"),
	}, emailHandler.DeleteCurrentUserEmail())

	return nil
}
