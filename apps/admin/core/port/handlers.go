package port

import (
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
)

type (
	ChecksHandler interface {
		Check() gin.HandlerFunc
		CheckInfo() []fizz.OperationOption
	}

	DatabaseHandler interface {
		GetCurrentDbms() gin.HandlerFunc
		GetCurrentDbmsInfo() []fizz.OperationOption

		MigrateTo() gin.HandlerFunc
		MigrateToInfo() []fizz.OperationOption
	}

	HardwareHandler interface {
		GetHost() gin.HandlerFunc
		GetHostInfo() []fizz.OperationOption

		GetCPUs() gin.HandlerFunc
		GetCPUsInfo() []fizz.OperationOption

		Reboot() gin.HandlerFunc
		RebootInfo() []fizz.OperationOption
	}

	HardwareKernelHandler interface {
		Reboot() gin.HandlerFunc
		RebootInfo() []fizz.OperationOption
	}

	SettingsHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		Patch() gin.HandlerFunc
		PatchInfo() []fizz.OperationOption
	}

	SshHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		Add() gin.HandlerFunc
		AddInfo() []fizz.OperationOption

		Delete() gin.HandlerFunc
		DeleteInfo() []fizz.OperationOption

		GetUsers() gin.HandlerFunc
		GetUsersInfo() []fizz.OperationOption
	}

	SshKernelHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		Add() gin.HandlerFunc
		AddInfo() []fizz.OperationOption

		Delete() gin.HandlerFunc
		DeleteInfo() []fizz.OperationOption

		GetUsers() gin.HandlerFunc
		GetUsersInfo() []fizz.OperationOption
	}

	UpdateHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		Install() gin.HandlerFunc
		InstallInfo() []fizz.OperationOption
	}
)
