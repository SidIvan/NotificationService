package repo

import (
	"NotificationService/internal/dto"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientRepoImpl struct {
	collection *mongo.Collection
}

type ClientRepoInterface interface {
	IsClientExist(string) bool
	FindClient(string) *dto.Client
	FindClientsByFilter(*dto.Filter) *[]dto.ClientWithId
	CreateClient(*dto.Client) string
	DeleteClient(string)
	UpdateClient(string, *dto.Client)
}

func (r ClientRepoImpl) IsClientExist(phoneNumber string) bool {
	res := r.collection.FindOne(context.TODO(), bson.D{{Key: "phoneNumber", Value: phoneNumber}})
	return res.Err() != mongo.ErrNoDocuments
}

func (r ClientRepoImpl) FindClient(phoneNumber string) *dto.Client {
	res := r.collection.FindOne(context.TODO(), bson.D{{Key: "phoneNumber", Value: phoneNumber}})
	if res.Err() == mongo.ErrNoDocuments {
		return nil
	}
	var client dto.Client
	err := res.Decode(&client)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &client
}

func (r ClientRepoImpl) FindClientsByFilter(filter *dto.Filter) *[]dto.ClientWithId {
	cursor, err := r.collection.Find(context.TODO(), *filter)
	if err != nil {
		log.Println(err)
		return nil
	}
	var clients []dto.ClientWithId
	for i := 0; cursor.Next(context.TODO()); i++ {
		clients = append(clients, dto.ClientWithId{})
		err = cursor.Decode(&clients[i])
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	return &clients
}

func (r ClientRepoImpl) CreateClient(client *dto.Client) string {
	res, err := r.collection.InsertOne(context.TODO(), client)
	if err != nil {
		log.Println(err)
		return ""
	}
	return res.InsertedID.(primitive.ObjectID).Hex()
}

func (r ClientRepoImpl) DeleteClient(phoneNumber string) {
	_, err := r.collection.DeleteOne(context.TODO(), bson.D{{Key: "phoneNumber", Value: phoneNumber}})
	if err != nil {
		log.Println(err)
	}
}

func (r ClientRepoImpl) UpdateClient(hexId string, newClientInfo *dto.Client) {
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = r.collection.UpdateByID(context.TODO(), id, bson.D{{Key: "$set", Value: newClientInfo}})
	if err != nil {
		log.Println(err)
	}
}
