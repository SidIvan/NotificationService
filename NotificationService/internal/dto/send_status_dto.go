package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type SendStatus struct {
	DistributionId primitive.ObjectID `bson:"distributionId"`
	ClientId       primitive.ObjectID `bson:"clientId"`
	Status         MessageStatus      `bson:"status"`
}
