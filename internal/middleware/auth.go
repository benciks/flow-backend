package middleware

import (
	"context"
	"database/sql"
	"errors"
	"github.com/benciks/flow-backend/internal/database/db"
	"github.com/labstack/echo/v4"
)

type contextKey string

var UserCtxKey = contextKey("user")

type UserContext struct {
	db.User
	Token string
}

func Auth(conn *db.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			reqContext := c.Request().Context()

			userContext := AuthMiddlewareFunction(reqContext, conn, token)
			if userContext == nil {
				return next(c)
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, UserCtxKey, userContext)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func AuthMiddlewareFunction(ctx context.Context, conn *db.Queries, token string) *UserContext {
	// Allow unauthenticated users in
	if token == "" {
		return nil
	}

	session, err := conn.GetSessionByToken(ctx, token)
	if err != nil {
		return nil
	}

	user, err := conn.FindUserById(ctx, session.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return &UserContext{
		user,
		session.Token,
	}
}

func GetUser(ctx context.Context) (*UserContext, bool) {
	raw, ok := ctx.Value(UserCtxKey).(*UserContext)
	return raw, ok && raw != nil
}
