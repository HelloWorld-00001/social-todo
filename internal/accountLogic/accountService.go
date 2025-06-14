package accountLogic

import (
	"errors"
	ulg "github.com/coderconquerer/go-login-app/internal/UserLoginLogic"
	"time"
)

type AccountService struct {
	repo *AccountRepository
}

func NewService(repo *AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Register(username, password string) error {
	acc, err := s.repo.GetByUsername(username)
	if err != nil {
		return err
	}

	if acc != nil {
		return errors.New("account already exists")
	}
	salt, err := GenerateSalt()
	if err != nil {
		return err
	}

	hashed, err := HashWithSalt(password, salt)
	if err != nil {
		return err
	}

	acc = &Account{
		Username:     username,
		HashPassword: hashed,
		Salt:         salt,
	}

	return s.repo.CreateAccount(acc)
}

func (s *AccountService) Login(username, password string) (*ulg.UserSession, error) {
	acc, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if acc == nil {
		return nil, nil
	}

	ok := VerifyWithSalt(acc.HashPassword, acc.Salt, password)
	if !ok {
		return nil, errors.New("Cannot verify password")
	}
	user := ulg.UserSession{
		ID:            acc.ID,
		Name:          acc.Username,
		PingCount:     0,
		ExpireSession: time.Now().Add(30 * time.Minute), // expires in 30 minutes
	}
	return &user, nil
}
