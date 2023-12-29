package impl

import (
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Jourloy/Go-Budget-Service/internal/auth"
	"github.com/Jourloy/Go-Budget-Service/internal/auth/impl/strategy"
	"github.com/Jourloy/Go-Budget-Service/internal/storage"
)

type authService struct {
	storage *storage.Storage
	local   strategy.LocalStrategy
}

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[auth]`,
		Level:  log.DebugLevel,
	})
)

var (
	Secret string
	Domain string
)

func parseENV() {
	if env, exist := os.LookupEnv(`SECRET`); !exist {
		logger.Fatal(`Error loading SECRET from .env file`)
	} else {
		Secret = env
	}

	if env, exist := os.LookupEnv(`DOMAIN`); !exist {
		logger.Fatal(`Error loading DOMAIN from .env file`)
	} else {
		Domain = env
	}
}

func CreateAuthService(s *storage.Storage) auth.AuthService {
	parseENV()

	localStategy := strategy.CreateLocalStrategy(s)

	return &authService{
		storage: s,
		local:   *localStategy,
	}
}

// REGISTER LOGIC //

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *authService) Login(c *gin.Context) {
	s.local.Login(c)
}

func (s *authService) Register(c *gin.Context) {
	s.local.Register(c)
}

// LOGIN LOGIC //

func (s *authService) GetUserData(c *gin.Context) {
	// Get cookies
	a, err := c.Cookie(`access_token`)
	if err != nil {
		logger.Error(`failed to get access token`, `err`, err)
		c.String(http.StatusBadRequest, `failed to get access token`)
		return
	}
	r, err := c.Cookie(`refresh_token`)
	if err != nil {
		logger.Error(`failed to get refresh token`, `err`, err)
		c.String(http.StatusBadRequest, `failed to get refresh token`)
		return
	}

	if a == `` || r == `` {
		logger.Error(`failed to get user data`)
		c.String(http.StatusBadRequest, `failed to get user data`)
		return
	}

	// Verify and decode
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(a, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	}); err != nil {
		logger.Error(`failed to verify access token`, `err`, err)
		c.String(http.StatusBadRequest, `failed to get user data`)
		return
	}

	// Check username and role
	var username string
	var role string

	for key, val := range claims {
		if key == `username` {
			username = val.(string)
		}
		if key == `role` {
			role = val.(string)
		}
	}

	if username == `` || role == `` {
		logger.Error(`failed to get user data`)
		c.String(http.StatusBadRequest, `failed to get user data`)
		return
	}

	// Get actual user info
	user, err := s.storage.User.GetUserByUsername(username)
	if err != nil {
		logger.Error(`failed to get user`, `err`, err)
		c.String(http.StatusInternalServerError, `failed to get user`)
		return
	}

	if user == nil || user.Username == `` {
		logger.Error(`failed to get user data`)
		c.String(http.StatusBadRequest, `failed to get user data`)
		return
	}

	// Add JWT tokens to cookies
	if err := s.addJWTCookies(UserData{Username: user.Username}, c); err != nil {
		logger.Error(`failed to add JWT tokens to cookies`, `err`, err)
		c.String(http.StatusInternalServerError, `failed to add JWT tokens to cookies`)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		`username`: user.Username,
		`role`:     user.Role,
	})
}

func (s *authService) addJWTCookies(body UserData, c *gin.Context) error {
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

func (s *authService) generateJWTTokens(username string, role string) (string, string) {
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
