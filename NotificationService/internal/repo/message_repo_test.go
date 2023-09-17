package repo

import (
	"NotificationService/internal/dto"
	"NotificationService/internal/utils"
	"context"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	testSendedAt          = "2023-05-23T10:00:00"
	testDistributionId, _ = primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	testClientId, _       = primitive.ObjectIDFromHex("507f191e810c19729de860ea")
	testMessageDto        = dto.Message{
		SendedAt:        testSendedAt,
		DisctributionId: testDistributionId.Hex(),
		ClientId:        testClientId.Hex(),
	}
)

func beforeMessageRepoTest() {
	utils.PMan = utils.NewPman("test.properties")
	ConnectToMongo(context.Background(), "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	DropDb()
}

func TestCreateSuccessMessage(t *testing.T) {
	beforeMessageRepoTest()
	testMessage := testMessageDto
	messageId, err := primitive.ObjectIDFromHex(MessageRepo.CreateSuccessMessage(&testMessage))
	if err != nil {
		t.Errorf("wrong id format of created message")
		return
	}
	cur, err := MessageRepo.collection.Find(context.TODO(), bson.D{
		{Key: "_id", Value: messageId},
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	var savedMessage dto.Message
	if !cur.Next(context.TODO()) {
		t.Errorf("message has not saved")
		return
	}
	if cur.Decode(&savedMessage) != nil {
		t.Errorf("cursor decoding error")
		return
	}
	if savedMessage.Status != dto.OkStatus {
		t.Errorf("wrong message status (should be OK)")
		return
	}
	if !reflect.DeepEqual(savedMessage, testMessage) {
		t.Errorf("message dto corrupted")
		return
	}
	if cur.Next(context.TODO()) {
		t.Errorf("not unique message id")
		return
	}
}

func TestCreateFailedMessage(t *testing.T) {
	beforeMessageRepoTest()
	testMessage := testMessageDto
	messageId, err := primitive.ObjectIDFromHex(MessageRepo.CreateFailedMessage(&testMessage))
	if err != nil {
		t.Errorf("wrong id format of created message")
		return
	}
	cur, err := MessageRepo.collection.Find(context.TODO(), bson.D{
		{Key: "_id", Value: messageId},
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	var savedMessage dto.Message
	if !cur.Next(context.TODO()) {
		t.Errorf("message has not saved")
		return
	}
	if cur.Decode(&savedMessage) != nil {
		t.Errorf("cursor decoding error")
		return
	}
	if savedMessage.Status != dto.NotOkStatus {
		t.Errorf("wrong message status (should be not OK)")
		return
	}
	if !reflect.DeepEqual(savedMessage, testMessage) {
		t.Errorf("message dto corrupted")
		return
	}
	if cur.Next(context.TODO()) {
		t.Errorf("not unique message id")
		return
	}
}

func TestGetStatusOk(t *testing.T) {
	beforeMessageRepoTest()
	MessageRepo.collection.InsertMany(context.TODO(), []interface{}{
		dto.Message{
			SendedAt:        "2000-12-12T19:30:00",
			Status:          dto.NotOkStatus,
			DisctributionId: testDistributionId.Hex(),
			ClientId:        testClientId.Hex(),
		},
		dto.Message{
			SendedAt:        "2000-12-12T19:31:00",
			Status:          dto.NotOkStatus,
			DisctributionId: testDistributionId.Hex(),
			ClientId:        testClientId.Hex(),
		},
		dto.Message{
			SendedAt:        "2000-12-12T19:32:00",
			Status:          dto.OkStatus,
			DisctributionId: testDistributionId.Hex(),
			ClientId:        testClientId.Hex(),
		},
	})
	if MessageRepo.GetStatus(testDistributionId, testClientId) != dto.OkStatus {
		t.Errorf("wrong message status detection (should be ok)")
	}
}

func TestGetStatusNotOk(t *testing.T) {
	beforeMessageRepoTest()
	MessageRepo.collection.InsertMany(context.TODO(), []interface{}{
		dto.Message{
			SendedAt:        "2000-12-12T19:30:00",
			Status:          dto.NotOkStatus,
			DisctributionId: testDistributionId.Hex(),
			ClientId:        testClientId.Hex(),
		},
		dto.Message{
			SendedAt:        "2000-12-12T19:31:00",
			Status:          dto.NotOkStatus,
			DisctributionId: testDistributionId.Hex(),
			ClientId:        testClientId.Hex(),
		},
		dto.Message{
			SendedAt:        "2000-12-12T19:32:00",
			Status:          dto.NotOkStatus,
			DisctributionId: testDistributionId.Hex(),
			ClientId:        testClientId.Hex(),
		},
	})
	if MessageRepo.GetStatus(testDistributionId, testClientId) != dto.NotOkStatus {
		t.Errorf("wrong message status detection (should be not ok)")
	}
}
