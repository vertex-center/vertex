package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	ChecksHandler interface {
		Check() gin.HandlerFunc
		CheckInfo() []oapi.Info
	}

	DatabaseHandler interface {
		GetCurrentDbms() gin.HandlerFunc
		GetCurrentDbmsInfo() []oapi.Info

		MigrateTo() gin.HandlerFunc
		MigrateToInfo() []oapi.Info
	}

	HardwareHandler interface {
		GetHost() gin.HandlerFunc
		GetHostInfo() []oapi.Info

		GetCPUs() gin.HandlerFunc
		GetCPUsInfo() []oapi.Info

		Reboot() gin.HandlerFunc
		RebootInfo() []oapi.Info
	}

	HardwareKernelHandler interface {
		Reboot() gin.HandlerFunc
		RebootInfo() []oapi.Info
	}

	SettingsHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		Patch() gin.HandlerFunc
		PatchInfo() []oapi.Info
	}

	SshHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		Add() gin.HandlerFunc
		AddInfo() []oapi.Info

		Delete() gin.HandlerFunc
		DeleteInfo() []oapi.Info

		GetUsers() gin.HandlerFunc
		GetUsersInfo() []oapi.Info
	}

	SshKernelHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		Add() gin.HandlerFunc
		AddInfo() []oapi.Info

		Delete() gin.HandlerFunc
		DeleteInfo() []oapi.Info

		GetUsers() gin.HandlerFunc
		GetUsersInfo() []oapi.Info
	}

	UpdateHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		Install() gin.HandlerFunc
		InstallInfo() []oapi.Info
	}
)
