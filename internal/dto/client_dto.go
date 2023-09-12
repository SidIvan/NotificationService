package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ClientInterface interface {
	IsValid() bool
	WithId(string) *ClientWithId
	GetId() string
}

type Client struct {
	PhoneNumber string `bson:"phoneNumber"`
	OpCode      string `bson:"opCode"`
	Tag         string `bson:"tag"`
	Utc         string `bson:"gmt"`
}

type ClientWithId struct {
	PhoneNumber string             `bson:"phoneNumber"`
	OpCode      string             `bson:"opCode"`
	Tag         string             `bson:"tag"`
	Utc         string             `bson:"gmt"`
	Id          primitive.ObjectID `bson:"_id" json:"id"`
}

func (c *Client) IsValid() bool {
	return IsOppCodeValid(c.OpCode) && IsPhoneNumberValid(c.PhoneNumber) && IsUtcValid(c.Utc) && c.PhoneNumber[1:4] == c.OpCode[:]
}

func (c *Client) WithId(hexId string) *ClientWithId {
	id, _ := primitive.ObjectIDFromHex(hexId)
	cWI := ClientWithId{
		PhoneNumber: c.PhoneNumber,
		OpCode:      c.OpCode,
		Tag:         c.Tag,
		Utc:         c.Utc,
		Id:          id,
	}
	return &cWI
}

func (c *ClientWithId) GetId() string {
	return c.Id.Hex()
}
