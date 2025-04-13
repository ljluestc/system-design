mkdir -p internal/auth && echo "$(cat <<EOF
	package auth
	
	import (
		"errors"
		"time"
		"github.com/dgrijalva/jwt-go"
		"messenger/pkg/db"
	)
	
	type Service struct {
		db     *db.Cassandra
		secret string
	}
	
	func NewService(db *db.Cassandra) *Service {
		return &Service{db: db, secret: "my-secret-key"} // Hardcoded for simplicity
	}
	
	func (s *Service) Login(email, password string) (string, error) {
		user, err := s.db.GetUserByEmail(email)
		if err != nil || user.Password != password {
			return "", errors.New("invalid credentials")
		}
	
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})
		return token.SignedString([]byte(s.secret))
	}
	
	func (s *Service) ValidateToken(tokenStr string) (string, error) {
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secret), nil
		})
		if err != nil || !token.Valid {
			return "", errors.New("invalid token")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return "", errors.New("invalid claims")
		}
		return claims["user_id"].(string), nil
	}
	EOF
	)" > internal/auth/auth.go