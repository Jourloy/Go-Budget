package strategy

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/Jourloy/Go-Budget-Service/internal/storage"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LocalStrategy struct {
	storage *storage.Storage
}

var (
	localLogger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[auth-local]`,
		Level:  log.DebugLevel,
	})
)

var (
	Secret string
	Domain string
)

func parseLocalENV() {
	if env, exist := os.LookupEnv(`SECRET`); !exist {
		localLogger.Fatal(`Error loading SECRET from .env file`)
	} else {
		Secret = env
	}

	if env, exist := os.LookupEnv(`DOMAIN`); !exist {
		localLogger.Fatal(`Error loading DOMAIN from .env file`)
	} else {
		Domain = env
	}
}

func CreateLocalStrategy(s *storage.Storage) *LocalStrategy {
	parseLocalENV()

	return &LocalStrategy{
		storage: s,
	}
}

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *LocalStrategy) Login(c *gin.Context) {
	// Parse body
	var body UserData
	ok, err := s.parseBody(c, &body)
	if !ok || body.Username == `` || body.Password == `` {
		localLogger.Error(`failed to parse body`, `err`, err)
		c.String(400, `failed to parse body`)
		return
	}

	// Get user by username
	user, err := s.storage.User.GetUserByUsername(body.Username)
	if err != nil {
		localLogger.Error(`failed to get user`, `err`, err)
		c.String(500, `failed to get user`)
		return
	}

	if user == nil {
		// If user doesn't exist create
		// Register user
		if err := s.registerUser(body.Username, body.Password); err != nil {
			localLogger.Error(`failed to register user`, `err`, err)
			c.String(500, `failed to register user`)
			return
		}

		// Add JWT tokens to cookies
		if err := s.addJWTCookies(body, c); err != nil {
			localLogger.Error(`failed to add JWT tokens to cookies`, `err`, err)
			c.String(500, `failed to add JWT tokens to cookies`)
			return
		}

		// Get actual user info
		newUser, err := s.storage.User.GetUserByUsername(body.Username)
		if err != nil {
			localLogger.Error(`failed to get user`, `err`, err)
			c.String(500, `failed to get user`)
			return
		}

		c.JSON(200, gin.H{
			`username`: newUser.Username,
			`role`:     newUser.Role,
		})
		return
	} else {
		// Check password
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			localLogger.Error(`invalid credentials`, `err`, err)
			c.String(403, `invalid credentials`)
			return
		}
	}

	// Add JWT tokens to cookies
	if err := s.addJWTCookies(body, c); err != nil {
		localLogger.Error(`failed to add JWT tokens to cookies`, `err`, err)
		c.String(500, `failed to add JWT tokens to cookies`)
		return
	}

	localLogger.Debug(`logged in`, `username`, user.Username)

	c.JSON(200, gin.H{
		`username`: user.Username,
		`role`:     user.Role,
	})
}

func (s *LocalStrategy) Register(c *gin.Context) {
	// Check body
	if c.Request.Body == nil {
		localLogger.Error(`body not found`)
		c.String(400, `body not found`)
		return
	}

	// Read body
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		localLogger.Error(`failed to read body`, `err`, err)
		c.String(500, `failed to read body`)
		return
	}
	defer c.Request.Body.Close()

	// Unmarshal
	var body UserData
	if err := json.Unmarshal(b, &body); err != nil {
		localLogger.Error(`failed to unmarshal body`, `err`, err)
		c.String(500, `failed to unmarshal body`)
		return
	}

	// Register user
	if err := s.registerUser(body.Username, body.Password); err != nil {
		localLogger.Error(`failed to register user`, `err`, err)
		c.String(500, `failed to register user`)
		return
	}

	// Add JWT tokens to cookies
	if err := s.addJWTCookies(body, c); err != nil {
		localLogger.Error(`failed to add JWT tokens to cookies`, `err`, err)
		c.String(500, `failed to add JWT tokens to cookies`)
		return
	}

	// Get actual user info
	newUser, err := s.storage.User.GetUserByUsername(body.Username)
	if err != nil {
		localLogger.Error(`failed to get user`, `err`, err)
		c.String(500, `failed to get user`)
		return
	}

	c.JSON(200, gin.H{
		`username`: newUser.Username,
		`role`:     newUser.Role,
	})
}

func (s *LocalStrategy) registerUser(username string, password string) error {
	// Get user by username
	user, err := s.storage.User.GetUserByUsername(username)
	if err != nil {
		return errors.New(`failed to get user`)
	}

	if user != nil {
		return errors.New(`user already exists`)
	}

	// Create user
	if err := s.storage.User.CreateUser(&users.UserCreate{
		Username: username,
		Password: password,
		Role:     `user`,
	}); err != nil {
		return errors.New(`failed to create user`)
	}

	localLogger.Debug(`registered`, `username`, username)

	return nil
}

func (s *LocalStrategy) addJWTCookies(body UserData, c *gin.Context) error {
	a, r := s.generateJWTTokens(body.Username, `user`)

	user, err := s.storage.User.GetUserByUsername(body.Username)
	if err != nil {
		return err
	}
	user.RefreshTokens = append(user.RefreshTokens, r)
	err = s.storage.User.UpdateUser(user)
	if err != nil {
		return err
	}

	c.SetCookie(`access_token`, a, 60*60*24, `/`, `localhost`, true, true)
	c.SetCookie(`refresh_token`, r, 60*60*24, `/`, `localhost`, true, true)

	return nil
}

func (s *LocalStrategy) generateJWTTokens(username string, role string) (string, string) {
	// Generate tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		`username`: username,
		`role`:     role,
	})
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	// Sign tokens
	accessTokenString, _ := accessToken.SignedString([]byte(Secret))
	refreshTokenString, _ := refreshToken.SignedString([]byte(Secret))

	return accessTokenString, refreshTokenString
}

func (s *LocalStrategy) parseBody(c *gin.Context, body interface{}) (bool, error) {
	// Check body
	if c.Request.Body == nil {
		return false, errors.New(`body not found`)
	}

	// Read body
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return false, err
	}

	// Unmarshal
	if err := json.Unmarshal(b, &body); err != nil {
		return false, err
	}

	return true, nil
}
