package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Hibiki-music-app/hibiki/backend/internal/repository"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

func BeginRegistration(w *webauthn.WebAuthn, c echo.Context) error {
	log.Default().Printf("[INFO] begin registration ----------------------\\")
	log.Default().Printf("Origin reçu : %s", c.Request().Header.Get("Origin"))

	username, err := getUsername(c)
	if err != nil {
		log.Default().Printf("[ERRO] can't get user name: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, _ := repository.GetOrCreateUser(username) // Find or create the new user

	options, session, err := w.BeginRegistration(user)
	if err != nil {
		msg := fmt.Sprintf("can't begin registration: %s", err.Error())
		log.Default().Printf("[ERRO] %s", msg)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}

	// Make a session key and store the sessionData values
	token, err := repository.GenSessionID()
	if err != nil {
		log.Default().Printf("[ERRO] can't generate session id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	repository.SaveSession(token, *session)

	c.SetCookie(&http.Cookie{
		Name:     "sid",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // SameSiteStrictMode maybe?
	})
	// return the options generated with the session key
	// options.publicKey contain our registration options
	return c.JSON(http.StatusOK, options)
}

func FinishRegistration(c echo.Context) error {
	// Get the session key from cookie
	sid, err := c.Cookie("sid")
	if err != nil {
		log.Default().Printf("[ERRO] can't get session id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing or invalid Session ID"})
	}

	// Get the session data stored from the function above
	session, _ := datastore.GetSession(sid.Value) // FIXME: cover invalid session

	// In out example username == userID, but in real world it should be different
	user := datastore.GetOrCreateUser(string(session.UserID)) // Get the user

	credential, err := webAuthn.FinishRegistration(user, session, c.Request())
	if err != nil {
		msg := fmt.Sprintf("can't finish registration: %s", err.Error())
		log.Default().Printf("[ERRO] %s", msg)
		// clean up sid cookie
		c.SetCookie(&http.Cookie{
			Name:  "sid",
			Value: "",
		})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}

	// If creation was successful, store the credential object
	user.AddCredential(credential)
	datastore.SaveUser(user)
	// Delete the session data
	datastore.DeleteSession(sid.Value)
	c.SetCookie(&http.Cookie{
		Name:  "sid",
		Value: "",
	})

	log.Default().Printf("[INFO] finish registration ----------------------/")

	return JSONResponse(c, "Registration Success", http.StatusOK) // Handle next steps
}

func BeginLogin(c echo.Context) error {
	log.Default().Printf("[INFO] begin login ----------------------\\")
	log.Default().Printf("Origin reçu : %s", c.Request().Header.Get("Origin"))

	username, err := getUsername(c)
	if err != nil {
		log.Default().Printf("[ERRO]can't get user name: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid username",
		})
	}

	user := datastore.GetOrCreateUser(username) // Find the user

	options, session, err := webAuthn.BeginLogin(user)
	if err != nil {
		msg := fmt.Sprintf("can't begin login: %s", err.Error())
		log.Default().Printf("[ERRO] %s", msg)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": msg,
		})
	}

	// Make a session key and store the sessionData values
	t, err := datastore.GenSessionID()
	if err != nil {
		log.Default().Printf("[ERRO] can't generate session id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate session ID",
		})
	}
	datastore.SaveSession(t, *session)

	c.SetCookie(&http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "/",
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // SameSiteStrictMode maybe?
	})

	log.Default().Printf("[DEBUG] BeginLogin saved session token=%s user=%s", t, username)
	return c.JSON(http.StatusOK, options) // return the options generated with the session key
	// options.publicKey contain our registration options
}

func FinishLogin(c echo.Context) error {
	// Get the session key from cookie
	sid, err := c.Cookie("sid")
	if err != nil {
		log.Default().Printf("[ERRO] no sid cookie: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// Get the session data stored from the function above
	session, ok := datastore.GetSession(sid.Value)
	if !ok {
		log.Default().Printf("[ERRO] session not found for token=%s", sid.Value)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid or expired session"})
	}

	if len(session.UserID) == 0 {
		log.Default().Printf("[ERRO] session.UserID empty for token=%s", sid.Value)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session user id"})
	}

	// In out example username == userID, but in real world it should be different
	user := datastore.GetOrCreateUser(string(session.UserID)) // Get the user
	if user == nil {
		log.Default().Printf("[ERRO] no user found for id=%v", session.UserID)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	credential, err := webAuthn.FinishLogin(user, session, c.Request())
	if err != nil {
		log.Default().Printf("[ERRO] can't finish login: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Handle credentialog.Default().Authenticator.CloneWarning
	if credentialog.Default().Authenticator.CloneWarning {
		log.Default().Printf("[WARN] can't finish login: %s", "CloneWarning")
	}

	// If login was successful, update the credential object
	user.UpdateCredential(credential)
	datastore.SaveUser(user)

	// Delete the login session data
	datastore.DeleteSession(sid.Value)
	c.SetCookie(&http.Cookie{
		Name:   "sid",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Add the new session cookie
	t, err := datastore.GenSessionID()
	if err != nil {
		log.Default().Printf("[ERRO] can't generate session id: %s", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate session ID",
		})
	}

	datastore.SaveSession(t, webauthn.SessionData{
		Expires: time.Now().Add(time.Hour),
		UserID:  []byte(user.WebAuthnName()),
	})
	c.SetCookie(&http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "/",
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // SameSiteStrictMode maybe?
	})

	log.Default().Printf("[INFO] finish login ----------------------/")

	return JSONResponse(c, "Login Success", http.StatusOK) // handle redirect?
}

func PrivatePage(c echo.Context) error {
	// test if everything works
	return c.String(http.StatusOK, "It works!")
}

// JSONResponse is a helper function to send json response
func JSONResponse(c echo.Context, data interface{}, status int) error {
	return c.JSON(status, data)
}

// getUsername is a helper function to extract the username from json request
func getUsername(c echo.Context) (string, error) {
	type Username struct {
		Username string `json:"username"`
	}
	var u Username
	if err := c.Bind(&u); err != nil {
		return "", err
	}
	return u.Username, nil
}
