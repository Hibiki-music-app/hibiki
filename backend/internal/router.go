package internal

import (
	"net/http"
	"time"

	"github.com/Hibiki-music-app/hibiki/backend/internal/interfaces"
	"github.com/Hibiki-music-app/hibiki/backend/internal/repository"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

// LoggedInMiddleware check if user is authenticated before
func LoggedInMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		sid, err := c.Cookie("sid")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		session, ok := repository.GetSession(webauthn.SessionData, token)
		if !ok || session.Expires.Before(time.Now()) {
			return c.Redirect(http.StatusSeeOther, "/")
		}
		return next(c)
	}
}
