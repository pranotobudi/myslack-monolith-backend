package mongodb_test

import (
	"context"

	"github.com/benweissmann/memongo"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type mockCollection struct{}
type MyMockedObject struct {
	mock.Mock
}

var mongoServer memongo.Server

func (m *MyMockedObject) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	c := &mongo.InsertOneResult{}
	return c, nil
}

// func TestInsertDoc(t *testing.T) {
// 	// create an instance of our test object
// 	// testObj := new(MyMockedObject)
// 	// // setup expectations
// 	// testObj.On("InsertOne", 123).Return(true, nil)

// 	// mockCol := mockCollection{}
// 	memongo.RandomDatabase()
// 	connectAndDoStuff(mongoServer.URI(), memongo.RandomDatabase())

// 	mongoServer, err := memongo.Start("4.0.5")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer mongoServer.Stop()

// 	mongodb.Inser
// 	mongodb.insertDoc(mockCol, name, BadExpr)
// }
