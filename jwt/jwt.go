package jwt

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/jatgam/wishlist-api/metrics"
	"github.com/jatgam/wishlist-api/microservice"
	"github.com/jatgam/wishlist-api/models"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

// JwtPayload is the data struct to be encoded/decoded from a jwt token.
type JwtPayload struct {
	ID            int  `json:"id"`
	PasswordReset bool `json:"passwordreset"`
	UserLevel     uint `json:"userlevel"`
}

// CreateJWTMiddleware creates the handlers for JWT authentication for use in
// Gin routes.
func CreateJWTMiddleware(jwtSecretKey string, realmName string) *jwt.GinJWTMiddleware {
	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            realmName,
		Key:              []byte(jwtSecretKey),
		SigningAlgorithm: "HS512",
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		IdentityKey:      identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*JwtPayload); ok {
				return jwt.MapClaims{
					identityKey:     v.ID,
					"passwordreset": v.PasswordReset,
					"userlevel":     v.UserLevel,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			mylogger := microservice.GetLogger(c)
			claims := jwt.ExtractClaims(c)
			mylogger.Debugf("JWT: Found Claims: %s", claims)
			return &JwtPayload{
				ID:            int(claims[identityKey].(float64)),
				PasswordReset: claims["passwordreset"].(bool),
				UserLevel:     uint(claims["userlevel"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			mylogger := microservice.GetLogger(c)
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				mylogger.Error("JWT: Missing Login Values")
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			foundUser, err := models.FindOneUser(&models.UserModel{Username: userID}, models.UserAuthScope)
			if err != nil {
				mylogger.Errorf("JWT: Failed to lookup user: %s", err.Error())
				metrics.FailedLogin.WithLabelValues(metrics.LoginFailedUser).Inc()
				return nil, jwt.ErrFailedAuthentication
			}
			if foundUser.ValidatePassword(password) {
				mylogger.Debugf("JWT: Successful Login: %s", userID)
				return &JwtPayload{
					ID:            foundUser.ID,
					PasswordReset: foundUser.PasswordReset,
					UserLevel:     foundUser.UserLevel,
				}, nil
			}
			mylogger.Errorf("JWT: Login Failed: %s", userID)
			metrics.FailedLogin.WithLabelValues(metrics.LoginFailedPassword).Inc()
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			mylogger := microservice.GetLogger(c)
			mylogger.Error("JWT: Invalid Login")
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": "Authentication Successful",
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": "Authentication Refresh Successful",
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
			})
		},

		TokenLookup: "header: Authorization",

		// TokenHeadName is a string in the header. The jwt module will
		// automatically set it if you try to disable.
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		logrus.Fatal("JWT: Error: " + err.Error())
	}

	return jwtMiddleware
}
