package rooms

import (
	"fmt"

	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRoomService interface {
	GetRooms() ([]mongodb.Room, error)
	GetAnyRoom() (mongodb.Room, error)
	AddRoom(name string) (string, error)
}
type roomService struct {
	repo mongodb.IMongoDB
}

func NewRoomService() *roomService {
	r := mongodb.NewMongoDB()
	return &roomService{repo: r}
}

func (s *roomService) GetRooms() ([]mongodb.Room, error) {
	return s.repo.GetRooms()
}

func (s *roomService) GetAnyRoom() (mongodb.Room, error) {
	anyRoom, err := s.repo.GetAnyRoom()
	if err != nil {
		return anyRoom, err
	}
	fmt.Println("inside room_io_handler-getRoom anyRoom!: ", anyRoom)
	objID, err := primitive.ObjectIDFromHex(anyRoom.ID)
	if err != nil {
		// panic(err)
		return anyRoom, err
	}

	filter := bson.M{"_id": objID}
	room, err := s.repo.GetRoom(filter)
	if err != nil {
		return room, err
	}
	fmt.Println("inside room_io_handler-getRoom!: ", room)
	return room, nil
}

func (s *roomService) AddRoom(name string) (string, error) {
	return s.repo.AddRoom(name)
}
