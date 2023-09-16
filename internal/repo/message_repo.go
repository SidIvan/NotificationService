package repo

import (
	"NotificationService/internal/dto"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepoImpl struct {
	collection *mongo.Collection
}

type MessageRepoInterface interface {
	CreateSuccessMessage(*dto.Message) string
	CreateFailedMessage(*dto.Message) string
	GetStatus(primitive.ObjectID, primitive.ObjectID) dto.MessageStatus
}

func (r MessageRepoImpl) CreateSuccessMessage(message *dto.Message) string {
	message.Status = dto.OkStatus
	return r.createMessage(message)
}

func (r MessageRepoImpl) CreateFailedMessage(message *dto.Message) string {
	message.Status = dto.NotOkStatus
	return r.createMessage(message)
}

func (r MessageRepoImpl) createMessage(message *dto.Message) string {
	res, err := r.collection.InsertOne(context.TODO(), message)
	if err != nil {
		log.Println(err)
		return ""
	}
	return res.InsertedID.(primitive.ObjectID).Hex()
}

func (r MessageRepoImpl) GetStatus(dId primitive.ObjectID, cId primitive.ObjectID) dto.MessageStatus {
	cur, err := r.collection.Find(context.TODO(), bson.D{
		{Key: "distributionId", Value: dId.Hex()},
		{Key: "clientId", Value: cId.Hex()},
	})
	if err != nil {
		log.Println(err)
		return ""
	}
	var message dto.Message
	for cur.Next(context.TODO()) {
		cur.Decode(&message)
		if message.Status == dto.OkStatus {
			return dto.OkStatus
		}
	}
	return dto.NotOkStatus
}
