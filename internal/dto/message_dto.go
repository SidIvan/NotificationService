package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageStatus string

const (
	OkStatus    MessageStatus = "success"
	NotOkStatus MessageStatus = "failed"
)

type Message struct {
	SendedAt        string        `bson:"sendedAt"`
	Status          MessageStatus `bson:"status"`
	DisctributionId string        `bson:"distributionId"`
	ClientId        string        `bson:"clientId"`
}

type MessageWithId struct {
	mongoId         primitive.ObjectID `bson:"_id"`
	SendedAt        string             `bson:"sendedAt"`
	Status          MessageStatus      `bson:"status"`
	DisctributionId string             `bson:"distributionId"`
	ClientId        string             `bson:"clientId"`
}

func (m *Message) IsValid() bool {
	return IsDateTimeValid(m.SendedAt)
}

func (m *Message) WithId(hexId string) *MessageWithId {
	id, _ := primitive.ObjectIDFromHex(hexId)
	mWI := MessageWithId{
		mongoId:         id,
		SendedAt:        m.SendedAt,
		Status:          m.Status,
		DisctributionId: m.DisctributionId,
		ClientId:        m.ClientId,
	}
	return &mWI
}
