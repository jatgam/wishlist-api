package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricLabelItemError  = "item_error"
	metricLabelLoginError = "login_error"
	metricLabelUserError  = "user_error"

	ItemAddError                   = "item_add_error"
	LoginFailedUser                = "login_invalid_user"
	LoginFailedPassword            = "login_invalid_password"
	RequestDataValidationError     = "data_validation_error"
	UserRegFailedError             = "user_reg_error"
	UserPasswordForgotError        = "user_pwforgot_error"
	UserPasswordResetValidateError = "user_pwresetvalidate_error"
	UserPasswordResetError         = "user_pwreset_error"
)

var (
	FailedLogin = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "wishlist_api_failed_login_total",
		Help: "The total number of failed logins",
	},
		[]string{metricLabelLoginError})

	UsersRegister = promauto.NewCounter(prometheus.CounterOpts{
		Name: "wishlist_api_users_registered_total",
		Help: "The total number of new registered users",
	})

	UserError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "wishlist_api_user_errors",
		Help: "The total number of errors encountered when dealing with users",
	},
		[]string{metricLabelUserError})

	ItemErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "wishlist_api_item_errors",
		Help: "Errors encountered when dealing with items",
	},
		[]string{metricLabelItemError})
)
