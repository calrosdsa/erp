	package middlewares

import (
	"context"
	"erp/internal/app/connection"
	"fmt"
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

func Authorization(adapter *connection.Adapter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			user := "rober"
			method := c.Request().Method
			path := c.Request().URL.Path
			log.Println(path)

			ok, err := enforce(ctx, user, path, method, adapter)
            if err != nil || !ok {
                return &echo.HTTPError{
                    Code:    http.StatusForbidden,
                    Message: "not allowed",
                }
            }
            if !ok {
                return err
            }
            return next(c)
		}
	}
}

func enforce(ctx context.Context,sub string,obj string,act string,adapter *connection.Adapter)(bool,error){
	enforcer,err := casbin.NewEnforcer("./static/casbin.conf",adapter)
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	err = enforcer.LoadPolicy()
    if err != nil {
        return false, fmt.Errorf("error in policy: %w", err)
    }
    // Verify
    ok, err := enforcer.Enforce(sub, obj, act)
    if err != nil {
        return false, fmt.Errorf("error in policy: %w", err)
    }
	return ok, nil
}