package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	webAuthn *webauthn.WebAuthn
	err      error

	datastore PasskeyStore
	//sessions  SessionStore
	l Logger
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type PasskeyUser interface {
	webauthn.User
	AddCredential(*webauthn.Credential)
	UpdateCredential(*webauthn.Credential)
}

type PasskeyStore interface {
	GetOrCreateUser(userName string) PasskeyUser
	SaveUser(PasskeyUser)
	GenSessionID() (string, error)
	GetSession(token string) (webauthn.SessionData, bool)
	SaveSession(token string, data webauthn.SessionData)
	DeleteSession(token string)
}

func main() {
	l = log.Default()

	proto := getEnv("PROTO", "http")
	host := getEnv("HOST", "localhost")
	port := getEnv("PORT", ":8080")
	origin := fmt.Sprintf("%s://%s%s", proto, host, port)

	// TODO: serverTLS config
	// serverCRT := getEnv("CRT", "")
	// serverKEY := getEnv("KEY", "")

	l.Printf("[INFO] make webauthn config")
	wconfig := &webauthn.Config{
		RPDisplayName: "HibikiGo Webauthn", // Display Name for your site
		RPID:          host,                // Generally the FQDN for your site
		RPOrigins:     []string{origin},    // The origin URLs allowed for WebAuthn
	}

	l.Printf("[INFO] create webauthn")
	if webAuthn, err = webauthn.New(wconfig); err != nil {
		fmt.Printf("[FATA] %s", err.Error())
		os.Exit(1)
	}

	l.Printf("[INFO] create datastore") //TODO: change to persist db
	datastore = NewInMem(l)

	e := echo.New()

	l.Printf("[INFO] setup CORS")
	e.Use(middleware.CORS()) // CORSWithConfig to enable the api endpoint only

	l.Printf("[INFO] register routes")

	// Add auth the routes
	api := e.Group("/api/passkey")
	api.POST("/registerStart", BeginRegistration)
	api.POST("/registerFinish", FinishRegistration)
	api.POST("/loginStart", BeginLogin)
	api.POST("/loginFinish", FinishLogin)

	// protected route
	api.GET("/private", PrivatePage, LoggedInMiddleware)

	// Start the server
	l.Printf("[INFO] start server at %s", origin)
	if err := e.Start(port); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func BeginRegistration(c echo.Context) error {
	l.Printf("[INFO] begin registration ----------------------\\")

	username, err := getUsername(c)
	if err != nil {
		l.Printf("[ERRO] can't get user name: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user := datastore.GetOrCreateUser(username) // Find or create the new user

	options, session, err := webAuthn.BeginRegistration(user)
	if err != nil {
		msg := fmt.Sprintf("can't begin registration: %s", err.Error())
		l.Printf("[ERRO] %s", msg)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}

	// Make a session key and store the sessionData values
	t, err := datastore.GenSessionID()
	if err != nil {
		l.Printf("[ERRO] can't generate session id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	datastore.SaveSession(t, *session)

	c.SetCookie(&http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "api/passkey/registerStart",
		MaxAge:   3600,
		Secure:   true,
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
		l.Printf("[ERRO] can't get session id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing or invalid Session ID"})
	}

	// Get the session data stored from the function above
	session, _ := datastore.GetSession(sid.Value) // FIXME: cover invalid session

	// In out example username == userID, but in real world it should be different
	user := datastore.GetOrCreateUser(string(session.UserID)) // Get the user

	credential, err := webAuthn.FinishRegistration(user, session, c.Request())
	if err != nil {
		msg := fmt.Sprintf("can't finish registration: %s", err.Error())
		l.Printf("[ERRO] %s", msg)
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

	l.Printf("[INFO] finish registration ----------------------/")

	return JSONResponse(c, "Registration Success", http.StatusOK) // Handle next steps
}

func BeginLogin(c echo.Context) error {
	l.Printf("[INFO] begin login ----------------------\\")

	username, err := getUsername(c)
	if err != nil {
		l.Printf("[ERRO]can't get user name: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid username",
		})
	}

	user := datastore.GetOrCreateUser(username) // Find the user

	options, session, err := webAuthn.BeginLogin(user)
	if err != nil {
		msg := fmt.Sprintf("can't begin login: %s", err.Error())
		l.Printf("[ERRO] %s", msg)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": msg,
		})
	}

	// Make a session key and store the sessionData values
	t, err := datastore.GenSessionID()
	if err != nil {
		l.Printf("[ERRO] can't generate session id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate session ID",
		})
	}
	datastore.SaveSession(t, *session)

	c.SetCookie(&http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "api/passkey/loginStart",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // SameSiteStrictMode maybe?
	})

	return c.JSON(http.StatusOK, options) // return the options generated with the session key
	// options.publicKey contain our registration options
}

func FinishLogin(c echo.Context) error {
	// Get the session key from cookie
	sid, err := c.Cookie("sid")
	if err != nil {
		l.Printf("[ERRO] can't get session id: %s", err.Error())

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// Get the session data stored from the function above
	session, _ := datastore.GetSession(sid.Value) // FIXME: cover invalid session

	// In out example username == userID, but in real world it should be different
	user := datastore.GetOrCreateUser(string(session.UserID)) // Get the user

	credential, err := webAuthn.FinishLogin(user, session, c.Request())
	if err != nil {
		l.Printf("[ERRO] can't finish login: %s", err.Error())
		panic(err)
	}

	// Handle credential.Authenticator.CloneWarning
	if credential.Authenticator.CloneWarning {
		l.Printf("[WARN] can't finish login: %s", "CloneWarning")
	}

	// If login was successful, update the credential object
	user.UpdateCredential(credential)
	datastore.SaveUser(user)

	// Delete the login session data
	datastore.DeleteSession(sid.Value)
	c.SetCookie(&http.Cookie{
		Name:  "sid",
		Value: "",
	})

	// Add the new session cookie
	t, err := datastore.GenSessionID()
	if err != nil {
		l.Printf("[ERRO] can't generate session id: %s", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate session ID",
		})
	}

	datastore.SaveSession(t, webauthn.SessionData{
		Expires: time.Now().Add(time.Hour),
	})
	c.SetCookie(&http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // SameSiteStrictMode maybe?
	})

	l.Printf("[INFO] finish login ----------------------/")

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

// getEnv is a helper function to get the environment variable
func getEnv(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return def
}

// LoggedInMiddleware check if user is authenticated before
func LoggedInMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		sid, err := c.Cookie("sid")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		session, ok := datastore.GetSession(sid.Value)
		if !ok || session.Expires.Before(time.Now()) {
			return c.Redirect(http.StatusSeeOther, "/")
		}
		return next(c)
	}
}
