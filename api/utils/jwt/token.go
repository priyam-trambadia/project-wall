package jwt

import "github.com/priyam-trambadia/project-wall/internal/logger"

func GenerateAccessToken(userID int64) (string, error) {
	logger := logger.Logger{Caller: "GenerateAccessToken utils/jwt"}

	token, err := generateToken(userID, AccessTokenExpiry, accessTokenSecret)
	if err != nil {
		err = logger.AppendError(err)
	}

	return token, err
}

func ValidateAccessToken(tokenString string) (int64, bool, error) {
	logger := logger.Logger{Caller: "ValidateAccessToken utils/jwt"}

	userID, ok, err := validateToken(tokenString, accessTokenSecret)
	if err != nil {
		err = logger.AppendError(err)
	}

	return userID, ok, err
}

func GenerateRefreshToken(userID int64) (string, error) {
	logger := logger.Logger{Caller: "GenerateRefreshToken utils/jwt"}

	token, err := generateToken(userID, RefreshTokenExpiry, refreshTokenSecret)
	if err != nil {
		err = logger.AppendError(err)
	}

	return token, err
}

func ValidateRefreshToken(tokenString string) (int64, bool, error) {
	logger := logger.Logger{Caller: "ValidateRefreshToken utils/jwt"}

	userID, ok, err := validateToken(tokenString, refreshTokenSecret)
	if err != nil {
		err = logger.AppendError(err)
	}

	return userID, ok, err
}

func GenerateUserActivationToken(userID int64) (string, error) {
	logger := logger.Logger{Caller: "GenerateUserActivationToken utils/jwt"}

	token, err := generateToken(userID, userActivationTokenExpiry, userActivationTokenSecret)
	if err != nil {
		err = logger.AppendError(err)
	}

	return token, err
}

func ValidateUserActivationToken(tokenString string) (int64, bool, error) {
	logger := logger.Logger{Caller: "ValidateUserActivationToken utils/jwt"}

	userID, ok, err := validateToken(tokenString, userActivationTokenSecret)
	if err != nil {
		err = logger.AppendError(err)
	}

	return userID, ok, err
}

func GenerateUserPasswordResetToken(userID int64) (string, error) {
	logger := logger.Logger{Caller: "GenerateUserPasswordResetToken utils/jwt"}

	token, err := generateToken(userID, userPasswordResetTokenExpiy, userPasswordResetTokenSecret)
	if err != nil {
		err = logger.AppendError(err)
	}

	return token, err
}

func ValidateUserPasswordResetToken(tokenString string) (int64, bool, error) {
	logger := logger.Logger{Caller: "ValidateUserPasswordResetToken utils/jwt"}

	userID, ok, err := validateToken(tokenString, userPasswordResetTokenSecret)
	if err != nil {
		err = logger.AppendError(err)
	}

	return userID, ok, err
}
