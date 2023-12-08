package port

import (
	"github.com/gin-gonic/gin"
)

type (
	ChecksHandler interface {
		Check() gin.HandlerFunc
	}

	DatabaseHandler interface {
		GetCurrentDbms() gin.HandlerFunc
		MigrateTo() gin.HandlerFunc
	}

	HardwareHandler interface {
		GetHost() gin.HandlerFunc
		GetCPUs() gin.HandlerFunc
		Reboot() gin.HandlerFunc
	}

	HardwareKernelHandler interface {
		Reboot() gin.HandlerFunc
	}

	SettingsHandler interface {
		Get() gin.HandlerFunc
		Patch() gin.HandlerFunc
	}

	SshHandler interface {
		Get() gin.HandlerFunc
		Add() gin.HandlerFunc
		Delete() gin.HandlerFunc
		GetUsers() gin.HandlerFunc
	}

	SshKernelHandler interface {
		Get() gin.HandlerFunc
		Add() gin.HandlerFunc
		Delete() gin.HandlerFunc
		GetUsers() gin.HandlerFunc
	}

	UpdateHandler interface {
		Get() gin.HandlerFunc
		Install() gin.HandlerFunc
	}
)
