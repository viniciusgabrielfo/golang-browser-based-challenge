package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"
	"go.uber.org/zap"
)

type service struct {
	repo   Repository
	logger *zap.SugaredLogger
}

func NewService(repo Repository, logger *zap.SugaredLogger) *service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateUser(nick, password string) error {
	if err := s.repo.Create(entity.NewUser(nick, password)); err != nil {
		s.logger.Error(err)
	}

	s.logger.Infow("new user created", "nick", nick)
	return nil
}

func (s *service) GetUser(id entity.UserID) (*entity.User, error) {
	return s.repo.Get(id)
}

func (s *service) Auth(nick, password string) (entity.UserID, error) {
	user, err := s.repo.GetByNick(nick)
	if err != nil {
		if !errors.Is(err, entity.ErrNotFoundEntity) {
			s.logger.Error(err)
		}
		return uuid.UUID{}, err
	}

	if user.Password != password {
		return user.ID, entity.ErrPassworNotMatch
	}

	return user.ID, nil
}
