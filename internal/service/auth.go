package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

var secretoJWT = []byte("codigogarra-api-secreto-demo")

// AuthService gestiona registro, login y validación de tokens JWT.
type AuthService struct {
	repo storage.UserRepository
}

func NuevoAuthService(repo storage.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}
	return s.repo.CrearUsuario(models.Usuario{
		Email:        email,
		PasswordHash: string(hash),
		CreadoEn:     time.Now(),
	})
}

func (s *AuthService) Login(email, password string) (string, error) {
	u, ok := s.repo.BuscarUsuarioPorEmail(email)
	if !ok {
		return "", ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.Email,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(secretoJWT)
}

func (s *AuthService) ValidarToken(tokenStr string) bool {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return secretoJWT, nil
	})
	return err == nil && t.Valid
}
