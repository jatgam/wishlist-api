package v1

import (
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/jatgam/wishlist-api/metrics"
	"github.com/jatgam/wishlist-api/microservice"
	"github.com/jatgam/wishlist-api/service"
	"github.com/jatgam/wishlist-api/types"
)

type registerUserForm struct {
	Username  string `form:"username" binding:"required,alphanum,min=3,notblank"`
	Password  string `form:"password" binding:"required,min=10,notblank,passcomplexity"`
	Email     string `form:"email" binding:"required,notblank,email"`
	Firstname string `form:"firstname" binding:"required,notblank,alpha,min=1"`
	Lastname  string `form:"lastname" binding:"required,notblank,alpha,min=1"`
}

type passwordForgotForm struct {
	Email string `form:"email" binding:"required,notblank,email"`
}

type passwordResetURI struct {
	PWResetToken string `uri:"pwResetToken" binding:"required,alphanum,min=40,max=40,notblank"`
}
type passwordResetForm struct {
	Password string `form:"password" binding:"required,min=10,notblank,passcomplexity"`
	Email    string `form:"email" binding:"required,notblank,email"`
}

func registerUser(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	var newUser registerUserForm
	if err := c.ShouldBind(&newUser); err != nil {
		mylogger.Debug("RegisterUser Failed Data Validation")
		metrics.UserError.WithLabelValues(metrics.RequestDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err := service.RegisterUser(strings.TrimSpace(newUser.Username), newUser.Password, strings.TrimSpace(newUser.Email),
		strings.TrimSpace(newUser.Firstname), strings.TrimSpace(newUser.Lastname), mylogger)
	if err != nil {
		mylogger.Errorf("Failed to Register User: %s:%s", newUser.Username, newUser.Email)
		metrics.UserError.WithLabelValues(metrics.UserRegFailedError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Infof("Registered New User: %s:%s", newUser.Username, newUser.Email)
	metrics.UsersRegister.Inc()
	types.WriteResponse(c, http.StatusOK, "User Created.")
}

func passwordForgot(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	var pwForgot passwordForgotForm
	if err := c.ShouldBind(&pwForgot); err != nil {
		mylogger.Debug("Password Forgot Failed Data Validation")
		metrics.UserError.WithLabelValues(metrics.RequestDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err := service.PasswordForgot(strings.TrimSpace(pwForgot.Email), c.Request.Host, mylogger)
	if err != nil {
		mylogger.Errorf("Failed to start password reset process: %s", pwForgot.Email)
		metrics.UserError.WithLabelValues(metrics.UserPasswordForgotError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Infof("Started PW Forgot: %s", pwForgot.Email)
	types.WriteResponse(c, http.StatusOK, "Sending an Email to the provided address.")
}

func passwordResetTokenValidate(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	var pwReset passwordResetURI
	if err := c.ShouldBindUri(&pwReset); err != nil {
		mylogger.Debug("PasswordReset Token Validate Data Validation")
		metrics.UserError.WithLabelValues(metrics.RequestDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	valid, err := service.PasswordResetTokenValidate(strings.TrimSpace(pwReset.PWResetToken), mylogger)
	if err != nil && err == types.ErrPasswordResetValidateServerErr {
		mylogger.Errorf("PasswordResetValidateToken failed to validate: %s", err.Error())
		metrics.UserError.WithLabelValues(metrics.UserPasswordResetValidateError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !valid {
		mylogger.Info("Password Reset Token INVALID")
		types.WriteResponse(c, http.StatusBadRequest, "Password reset token is invalid or expired.")
		return
	}

	mylogger.Info("Password Reset Token Valid")
	types.WriteResponse(c, http.StatusOK, "Token Valid")
}

func passwordReset(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	var pwReset passwordResetForm
	var pwResetURI passwordResetURI
	if err := c.ShouldBindUri(&pwResetURI); err != nil {
		mylogger.Debug("PasswordReset Token Data Validation")
		metrics.UserError.WithLabelValues(metrics.RequestDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := c.ShouldBind(&pwReset); err != nil {
		mylogger.Debug("PasswordReset Form Data Validation")
		metrics.UserError.WithLabelValues(metrics.RequestDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err := service.PasswordReset(strings.TrimSpace(pwReset.Email), pwReset.Password, pwResetURI.PWResetToken, mylogger)
	if err != nil {
		mylogger.Errorf("PasswordReset failed to reset: %s : %s", pwReset.Email, err.Error())
		metrics.UserError.WithLabelValues(metrics.UserPasswordResetError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Infof("Password reset for: %s", pwReset.Email)
	types.WriteResponse(c, http.StatusOK, "Password Reset")

}

func setupUserRoutes(router *gin.RouterGroup, ginjwt *jwt.GinJWTMiddleware) {

	router.GET("/auth", ginjwt.MiddlewareFunc())
	router.POST("/auth", ginjwt.LoginHandler)
	router.POST("/auth/refresh", ginjwt.RefreshHandler)

	router.POST("/register", registerUser)
	router.POST("/password_forgot", passwordForgot)
	router.GET("/password_reset/:pwResetToken", passwordResetTokenValidate)
	router.POST("/password_reset/:pwResetToken", passwordReset)
}
