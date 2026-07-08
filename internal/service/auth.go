package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

var secretoJWT = []byte("codigogarra-api-secreto-demo")

// AuthService gestiona registro, login y validación de tokens JWT con roles.
type AuthService struct {
	repo storage.UserRepository
}

func NuevoAuthService(repo storage.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Registrar crea un usuario con rol "veterinario" por defecto.
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
		Rol:          "veterinario",
		CreadoEn:     time.Now(),
	})
}

// RegistrarAdmin crea un usuario con rol "admin". Solo para uso interno (seeder).
func (s *AuthService) RegistrarAdmin(email, password string) (models.Usuario, error) {
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
		Rol:          "admin",
		CreadoEn:     time.Now(),
	})
}

// Login autentica al usuario y devuelve un JWT con el claim "rol" incluido.
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
		"rol": u.Rol,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(secretoJWT)
}

// ValidarToken parsea y valida el token. Devuelve los claims y si es válido.
func (s *AuthService) ValidarToken(tokenStr string) (jwt.MapClaims, bool) {
	t, err := jwt.Parse(tokenStr, func(*jwt.Token) (any, error) {
		return secretoJWT, nil
	})
	if err != nil || !t.Valid {
		return nil, false
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	return claims, ok
}
