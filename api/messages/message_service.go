package messages

import "github.com/pranotobudi/myslack-monolith-backend/mongodb"

type IMessageService interface {
	GetMessages(filter interface{}) ([]mongodb.Message, error)
}
type messageService struct {
	repo mongodb.IMongoDB
}

func NewMessageService() *messageService {
	r := mongodb.NewMongoDB()
	return &messageService{repo: r}
}

func (s *messageService) GetMessages(filter interface{}) ([]mongodb.Message, error) {
	messages, err := s.repo.GetMessages(filter)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
