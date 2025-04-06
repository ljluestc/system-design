package auth

import (
    "errors"
    "uber/internal/user"
    "uber/internal/utils"
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
    "time"
)

type AuthService struct {
    userSvc *user.UserService
    secret  string
    logger  *utils.Logger
}

// NewAuthService initializes the auth service
func NewAuthService(us *user.UserService, secret string, logger *utils.Logger) *AuthService {
    return &AuthService{userSvc: us, secret: secret, logger: logger}
}

// Login authenticates a user and returns a JWT token
func (as *AuthService) Login(email, password string) (string, error) {
    user, err := as.userSvc.GetUserByEmail(email)
    if err != nil {
        as.logger.Error("User not found:", err)
        return "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        as.logger.Error("Password mismatch:", err)
        return "", errors.New("invalid credentials")
    }

    token, err := GenerateJWT(user.ID, as.secret)
    if err != nil {
        as.logger.Error("Failed to generate token:", err)
        return "", err
    }

    as.logger.Info("User logged in:", user.ID)
    return token, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (as *AuthService) ValidateToken(tokenString string) (string, error) {
    userID, err := ValidateJWT(tokenString, as.secret)
    if err != nil {
        as.logger.Error("Token validation failed:", err)
        return "", err
    }
    return userID, nil
}

// GenerateJWT creates a JWT token
func GenerateJWT(userID, secret string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString([]byte(secret))
}

// ValidateJWT validates a JWT token
func ValidateJWT(tokenString, secret string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims["user_id"].(string), nil
    }
    return "", errors.New("invalid token")
}