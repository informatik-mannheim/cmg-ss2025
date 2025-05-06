package core

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/model"
	"golang.org/x/crypto/argon2"
)

type Service struct {
	mu       sync.RWMutex
	users    map[string]model.User
	filePath string
}

func NewService() *Service {
	s := &Service{
		users:    make(map[string]model.User),
		filePath: "users.json",
	}
	s.loadUsersFromFile()
	return s
}

func (s *Service) AddUser(id string, role model.Role) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; exists {
		return "", errors.New("user already exists")
	}

	if role == model.JobScheduler {
		for _, u := range s.users {
			if u.Role == model.JobScheduler {
				return "", errors.New("a Job Scheduler already exists")
			}
		}
	}

	plainSecret := generateSecret()
	hashed := hashSecret(plainSecret)

	s.users[id] = model.User{
		ID:     id,
		Role:   role,
		Secret: hashed,
	}

	s.saveUsersToFile()
	return plainSecret, nil
}

func (s *Service) Authenticate(secret string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if verifySecret(secret, user.Secret) {
			return &user, nil
		}
	}

	return nil, errors.New("invalid secret")
}

func generateSecret() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("failed to generate secret: " + err.Error())
	}
	return hex.EncodeToString(key)
}

func hashSecret(secret string) string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		panic("failed to generate salt: " + err.Error())
	}

	hash := argon2.IDKey([]byte(secret), salt, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)
}

func verifySecret(secret, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 2 {
		return false
	}
	salt, err1 := base64.RawStdEncoding.DecodeString(parts[0])
	expected, err2 := base64.RawStdEncoding.DecodeString(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}

	hash := argon2.IDKey([]byte(secret), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, expected) == 1
}

func (s *Service) loadUsersFromFile() {
	file, err := os.Open(s.filePath)
	if err != nil {
		return // no file yet
	}
	defer file.Close()

	_ = json.NewDecoder(file).Decode(&s.users)
}

func (s *Service) saveUsersToFile() {
	file, err := os.Create(s.filePath)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(s.users)
}
