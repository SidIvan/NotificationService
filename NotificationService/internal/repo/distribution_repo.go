package repo

import (
	"NotificationService/internal/dto"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DistributionRepoImpl struct {
	collection *mongo.Collection
}

type DistributionRepoInterface interface {
	CreateDistribution(*dto.Distribution) string
	IsDistributionExist(string) bool
	FindDistributionById(string) *dto.DistributionWithId
	FindAllDistributions() *[]dto.DistributionWithId
	DeleteDistribution(string)
	UpdateDistribution(string, *dto.Distribution)
}

func (r DistributionRepoImpl) CreateDistribution(distribution *dto.Distribution) string {
	res, err := r.collection.InsertOne(context.TODO(), distribution)
	if err != nil {
		log.Println(err)
		return ""
	}
	return res.InsertedID.(primitive.ObjectID).Hex()
}

func (r DistributionRepoImpl) IsDistributionExist(hexId string) bool {
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		log.Println(err)
		return false
	}
	res := r.collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}})
	return res.Err() != mongo.ErrNoDocuments
}

func (r DistributionRepoImpl) FindDistributionById(hexId string) *dto.DistributionWithId {
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		log.Println(err)
		return nil
	}
	res := r.collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}})
	if res.Err() == mongo.ErrNoDocuments {
		return nil
	}
	var disctribution dto.DistributionWithId
	err = res.Decode(&disctribution)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &disctribution
}

func (r DistributionRepoImpl) FindAllDistributions() *[]dto.DistributionWithId {
	var distributions []dto.DistributionWithId
	cursor, err := r.collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Println(err)
		return nil
	}
	for i := 0; cursor.Next(context.TODO()); i++ {
		distributions = append(distributions, dto.DistributionWithId{})
		err = cursor.Decode(&distributions[i])
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	return &distributions
}

func (r DistributionRepoImpl) DeleteDistribution(hexId string) {
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = r.collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		log.Println(err)
	}
}

func (r DistributionRepoImpl) UpdateDistribution(hexId string, newDistributionInfo *dto.Distribution) {
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = r.collection.UpdateByID(context.TODO(), id, bson.D{{Key: "$set", Value: newDistributionInfo}})
	if err != nil {
		log.Println(err)
	}
}
