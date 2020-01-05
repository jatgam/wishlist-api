package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jatgam/wishlist-api/models"
	"github.com/jatgam/wishlist-api/service/sgmail"
	"github.com/jatgam/wishlist-api/types"
	"github.com/jatgam/wishlist-api/utils"
)

func RegisterUser(username, password, email, firstname, lastname string, logger *logrus.Entry) error {
	username = strings.ToLower(username)
	email = strings.ToLower(email)
	userTaken, err := usernameTaken(username)
	if err != nil {
		logger.Errorf("User Registration Failed: %s", err)
		return types.ErrUserRegister
	}
	if userTaken {
		logger.Debugf("Username Taken: %s", username)
		return types.ErrUsernameTaken
	}
	emTaken, err := emailTaken(email)
	if err != nil {
		logger.Errorf("User Registration Failed: %s", err)
		return types.ErrUserRegister
	}
	if emTaken {
		logger.Debugf("EMail Taken: %s", email)
		return types.ErrEmailTaken
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		logger.Errorf("User Registration Password Hash Failed: %s", err)
		return types.ErrUserRegister
	}
	newUser := &models.UserModel{Username: username, PasswordHash: hash,
		EMail: email, FirstName: firstname, LastName: lastname, UserLevel: 1}
	err = models.CreateUser(newUser)
	if err != nil {
		logger.Errorf("User Registration DB Insert Failed: %s", err)
		return types.ErrUserRegister
	}
	return nil
}

func PasswordForgot(email, hostUrl string, logger *logrus.Entry) error {
	email = strings.ToLower(email)
	user, err := models.FindOneUser(&models.UserModel{EMail: email}, models.UserPassResetScope)
	if err != nil {
		logger.Errorf("PasswordForgot: Failed DB query: %s", email)
		return types.ErrPasswordForgot
	}
	if user == nil {
		logger.Errorf("PasswordForgot: No user record: %s", email)
		return types.ErrPasswordForgot
	}

	token, err := utils.RandomHex(20)
	if err != nil {
		logger.Errorf("PasswordForgot: Failed to Generate Reset Token: %s", err.Error())
		return types.ErrPasswordForgot
	}

	expirationTime := time.Now().Add(time.Hour)
	updates := models.UserModel{PasswordResetToken: &token, PasswordResetExpires: &expirationTime}
	err = models.UpdateUser(user, updates)
	if err != nil {
		logger.Errorf("PasswordForgot: Failed to Update User in DB: %s", err.Error())
		return types.ErrPasswordForgot
	}

	message := fmt.Sprintf("You are receiving this because you (or someone else) requested the reset of the password for your account.\n\n"+
		"Please click the following link, or paste into your browser to complete the process:\n\n"+
		"https://%s/password_reset/%s\n\n"+
		"Username: %s\n\n"+
		"If you did not request this, please ignore this email and your password will remain unchanged.",
		hostUrl, token, user.Username)

	mailer := sgmail.GetMailer()
	err = mailer.SendMail(email, "Jatgam Wishlist Password Reset", message, logger)
	if err != nil {
		logger.Errorf("PasswordForgot: Error Sending Reset Email: %s", err.Error())
		return types.ErrPasswordForgot
	}
	return nil
}

func PasswordResetTokenValidate(token string, logger *logrus.Entry) (bool, error) {
	user, err := models.FindOneUser(&models.UserModel{PasswordResetToken: &token}, models.UserPassResetScope)
	if err != nil {
		logger.Errorf("PasswordResetTokenValidate: Failed DB query: %s", token)
		return false, types.ErrPasswordResetValidate
	}
	if user == nil {
		logger.Errorf("PasswordResetTokenValidate: No user record: %s", token)
		return false, types.ErrPasswordResetValidate
	}

	if user.PasswordResetToken != nil && user.PasswordResetExpires != nil {
		if len(*user.PasswordResetToken) == 40 && *user.PasswordResetToken == token && user.PasswordResetExpires.After(time.Now()) {
			logger.Debugf("PasswordReset Token is Valid: %s", token)
			return true, nil
		}
	}
	logger.Debugf("PasswordReset Token is INVALID: %s", token)
	return false, nil
}

func PasswordReset(email, password, token string, logger *logrus.Entry) error {
	email = strings.ToLower(email)
	tokenValid, _ := PasswordResetTokenValidate(token, logger)
	if !tokenValid {
		logger.Debugf("Password Reset for %s:%s had invalid token.", email, token)
		return types.ErrPasswordResetValidate
	}

	user, err := models.FindOneUser(&models.UserModel{EMail: email}, models.UserAuthScope)
	if err != nil {
		logger.Errorf("PasswordReset: Failed DB query: %s", email)
		return types.ErrPasswordReset
	}
	if user == nil {
		logger.Errorf("PasswordReset: No user record: %s", email)
		return types.ErrPasswordReset
	}

	if user.PasswordResetToken != nil && *user.PasswordResetToken != token {
		logger.Errorf("PasswordReset: Potential multiple of same reset tokens exist")
		return types.ErrPasswordResetValidate
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		logger.Errorf("PasswordReset Password Hash Failed: %s", err)
		return types.ErrPasswordReset
	}
	updates := map[string]interface{}{"PasswordHash": hash, "PasswordReset": false, "PasswordResetToken": nil, "PasswordResetExpires": nil}
	// updates := models.UserModel{PasswordHash: hash, PasswordReset: false, PasswordResetToken: nil, PasswordResetExpires: nil}
	err = models.UpdateUserWithMap(user, updates)
	if err != nil {
		logger.Errorf("PasswordReset: Failed to Update User in DB: %s", err.Error())
		return types.ErrPasswordReset
	}

	return nil
}

func usernameTaken(username string) (bool, error) {
	user, err := models.FindOneUser(&models.UserModel{Username: username})
	if err != nil {
		return true, err
	}
	if user != nil {
		return true, nil
	}
	return false, nil
}

func emailTaken(email string) (bool, error) {
	user, err := models.FindOneUser(&models.UserModel{EMail: email})
	if err != nil {
		return true, err
	}
	if user != nil {
		return true, nil
	}
	return false, nil
}
