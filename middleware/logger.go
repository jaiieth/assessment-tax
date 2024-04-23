package middleware

import (
	"github.com/labstack/echo/v4/middleware"
)

var Logger = middleware.LoggerWithConfig(middleware.LoggerConfig{
	Format:           "${time_custom} [${status}]: ${method} ${uri}\n",
	CustomTimeFormat: "2006-01-02 15:04:05",
})
