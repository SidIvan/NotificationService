package repo

import (
	"NotificationService/internal/dto"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var messageCollection *mongo.Collection

func CreateSuccessMessage(message *dto.Message) string {
	message.Status = dto.OkStatus
	return createMessage(message)
}

func CreateFailedMessage(message *dto.Message) string {
	message.Status = dto.NotOkStatus
	return createMessage(message)
}

func createMessage(message *dto.Message) string {
	res, err := messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		log.Println(err)
		return ""
	}
	return res.InsertedID.(primitive.ObjectID).Hex()
}

func GetStatus(dId primitive.ObjectID, cId primitive.ObjectID) dto.MessageStatus {
	cur, err := messageCollection.Find(context.TODO(), bson.D{
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
