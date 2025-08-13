package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	webAuthn *webauthn.WebAuthn
	err      error

	datastore PasskeyStore
	l         Logger
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

	// Initialize PostgreSQL connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "hibikidb")
	dbPortStr := getEnv("DB_PORT", "5432")
	dbPortInt, err := strconv.Atoi(dbPortStr)
	if err != nil {
		l.Printf("[FATA] Invalid DB_PORT: %s", err.Error())
		os.Exit(1)
	}

	l.Printf("[INFO] Initializing database connection")
	InitDB(dbHost, dbUser, dbPassword, dbName, dbPortInt)

	l.Printf("[INFO] make webauthn config")
	wconfig := &webauthn.Config{
		RPDisplayName: "HibikiGo Webauthn",               // Display Name for your site
		RPID:          host,                              // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:5173"}, // The origin URLs allowed for WebAuthn
	}

	l.Printf("[INFO] create webauthn")
	if webAuthn, err = webauthn.New(wconfig); err != nil {
		fmt.Printf("[FATA] %s", err.Error())
		os.Exit(1)
	}

	l.Printf("[INFO] create datastore")
	datastore = NewPostgreSQLStore(l)

	e := echo.New()

	l.Printf("[INFO] setup CORS")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Cookie"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Set-Cookie"},
	}))

	l.Printf("[INFO] register routes")

	// Add auth routes
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

// clearSessionCookie helper function to clear session cookie properly
func clearSessionCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// setSessionCookie helper function to set session cookie properly
func setSessionCookie(c echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     "sid",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func BeginRegistration(c echo.Context) error {
	l.Printf("[INFO] begin registration ----------------------\\")
	l.Printf("Origin reçu : %s", c.Request().Header.Get("Origin"))

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
	setSessionCookie(c, t)

	// return the options generated with the session key
	return c.JSON(http.StatusOK, options)
}

func FinishRegistration(c echo.Context) error {
	l.Printf("[INFO] finish registration ----------------------\\")

	// Get the session key from cookie
	sid, err := c.Cookie("sid")
	if err != nil {
		l.Printf("[ERRO] can't get session id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing or invalid Session ID"})
	}

	// Get the session data stored from the function above
	session, ok := datastore.GetSession(sid.Value)
	if !ok {
		l.Printf("[ERRO] session not found for token=%s", sid.Value)
		clearSessionCookie(c)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid or expired session"})
	}

	// In our example username == userID, but in real world it should be different
	user := datastore.GetOrCreateUser(string(session.UserID)) // Get the user

	credential, err := webAuthn.FinishRegistration(user, session, c.Request())
	if err != nil {
		msg := fmt.Sprintf("can't finish registration: %s", err.Error())
		l.Printf("[ERRO] %s", msg)
		// Clear session on error
		datastore.DeleteSession(sid.Value)
		clearSessionCookie(c)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}

	// If creation was successful, store the credential object
	user.AddCredential(credential)
	datastore.SaveUser(user)

	// Delete the registration session
	datastore.DeleteSession(sid.Value)
	clearSessionCookie(c)

	l.Printf("[INFO] finish registration ----------------------/")

	return JSONResponse(c, "Registration Success", http.StatusOK)
}

func BeginLogin(c echo.Context) error {
	l.Printf("[INFO] begin login ----------------------\\")
	l.Printf("Origin reçu : %s", c.Request().Header.Get("Origin"))

	username, err := getUsername(c)
	if err != nil {
		l.Printf("[ERRO] can't get user name: %s", err.Error())
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
	setSessionCookie(c, t)

	l.Printf("[DEBUG] BeginLogin saved session token=%s user=%s", t, username)
	return c.JSON(http.StatusOK, options)
}

func FinishLogin(c echo.Context) error {
	l.Printf("[INFO] finish login ----------------------\\")

	// Get the session key from cookie
	sid, err := c.Cookie("sid")
	if err != nil {
		l.Printf("[ERRO] no sid cookie: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get the session data stored from the function above
	session, ok := datastore.GetSession(sid.Value)
	if !ok {
		l.Printf("[ERRO] session not found for token=%s", sid.Value)
		clearSessionCookie(c)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid or expired session"})
	}

	if len(session.UserID) == 0 {
		l.Printf("[ERRO] session.UserID empty for token=%s", sid.Value)
		datastore.DeleteSession(sid.Value)
		clearSessionCookie(c)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session user id"})
	}

	// In our example username == userID, but in real world it should be different
	user := datastore.GetOrCreateUser(string(session.UserID)) // Get the user
	if user == nil {
		l.Printf("[ERRO] no user found for id=%v", session.UserID)
		datastore.DeleteSession(sid.Value)
		clearSessionCookie(c)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	credential, err := webAuthn.FinishLogin(user, session, c.Request())
	if err != nil {
		l.Printf("[ERRO] can't finish login: %s", err.Error())
		datastore.DeleteSession(sid.Value)
		clearSessionCookie(c)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Handle credential.Authenticator.CloneWarning
	if credential.Authenticator.CloneWarning {
		l.Printf("[WARN] clone warning detected for user: %s", user.WebAuthnName())
	}

	// If login was successful, update the credential object
	user.UpdateCredential(credential)
	datastore.SaveUser(user)

	// Delete the login session data
	datastore.DeleteSession(sid.Value)

	// Create a new authenticated session
	authToken, err := datastore.GenSessionID()
	if err != nil {
		l.Printf("[ERRO] can't generate auth session id: %s", err.Error())
		clearSessionCookie(c)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate session ID",
		})
	}

	// Save authenticated session
	authSession := webauthn.SessionData{
		Expires: time.Now().Add(time.Hour * 24), // 24 hour session
		UserID:  []byte(user.WebAuthnName()),
	}
	datastore.SaveSession(authToken, authSession)
	setSessionCookie(c, authToken)

	l.Printf("[INFO] finish login ----------------------/")

	return JSONResponse(c, "Login Success", http.StatusOK)
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

// LoggedInMiddleware check if user is authenticated
func LoggedInMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sid, err := c.Cookie("sid")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		}

		session, ok := datastore.GetSession(sid.Value)
		if !ok || session.Expires.Before(time.Now()) {
			clearSessionCookie(c)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session expired"})
		}

		return next(c)
	}
}
