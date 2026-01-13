package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/EyeMart07/scheduler/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignUpReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) setSessionCookie(c *gin.Context, customerId string) error {
	b := make([]byte, 32) // 256-bit
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error signing in"})
		return err
	}

	raw := base64.RawURLEncoding.EncodeToString(b) // cookie-safe
	hash := sha256.Sum256([]byte(raw))             // hash the raw token string

	if err := a.Store.CreateSession(store.Session{
		Customer:    customerId,
		SessionHash: hash,
	}); err != nil {
		return err
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    raw,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //SET TO TRUE IN PROD
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return nil
}

func (a *App) CheckAdmin(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		c.AbortWithStatus(http.StatusUnauthorized) //user is not authorized
		return
	}

	if admin := a.Store.CheckAdmin(fmt.Sprintf("%v", user)); !admin {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		c.AbortWithStatus(http.StatusUnauthorized) //user is not authorized
		return
	}
	c.Next()
}

func (a *App) CheckAuth(c *gin.Context) {
	token, err := c.Cookie("session")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		c.AbortWithStatus(http.StatusUnauthorized) //user is not authorized
	}

	hash := sha256.Sum256([]byte(token))

	if userid := a.Store.CheckAuth(hash); userid != "" {
		c.Set("user", userid)
		c.Next()
		return
	}
	c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
	c.AbortWithStatus(http.StatusUnauthorized) //user is not authorized
}

func (a *App) SignUp(c *gin.Context) {
	var req SignUpReq

	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password

	if len(email) == 0 || !strings.Contains(email, "@") {
		c.JSON(400, gin.H{"error": "invalid email"})
		return
	}
	if len(password) < 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 12 characters"})
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error hashing password"})
		return
	}

	var customerId string
	customerId, err = a.Store.CreateUser(store.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     email,
		Password:  string(bytes), //passes the hashed password
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error creating account"})
		return
	}

	if err := a.setSessionCookie(c, customerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error creating session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "account successfully created"})

}

func (a *App) SignIn(c *gin.Context) {
	var credentials AuthReq

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error signing in"})
		fmt.Println(err)
		return
	}

	customerId := a.Store.Authorize(store.AuthReq{
		Email:    credentials.Email,
		Password: credentials.Password,
	})

	if customerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	if err := a.setSessionCookie(c, customerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error signing in"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged in"})

	//otherwise set cookies and add to session
}
