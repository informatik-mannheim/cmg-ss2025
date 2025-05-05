// core/service.go
package core

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/alexedwards/argon2id"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/model"
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
	hash, err := argon2id.CreateHash(plainSecret, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	s.users[id] = model.User{
		ID:     id,
		Role:   role,
		Secret: hash,
	}

	s.saveUsersToFile()
	return plainSecret, nil
}

func (s *Service) Authenticate(secret string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		match, err := argon2id.ComparePasswordAndHash(secret, user.Secret)
		if err == nil && match {
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
