package repo

import (
	"NotificationService/internal/dto"
	"NotificationService/internal/utils"
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	client1PhoneNumber = "712334567890"

	client1 = dto.Client{
		PhoneNumber: client1PhoneNumber,
		OpCode:      "123",
		Tag:         "tag1",
		Utc:         "+01:00",
	}
	filterTag    = "filterTag"
	filterOpCode = "903"
	client2      = dto.Client{
		PhoneNumber: "79039999999",
		OpCode:      filterOpCode,
		Tag:         filterTag,
		Utc:         "+02:00",
	}
	client3 = dto.Client{
		PhoneNumber: "79031236666",
		OpCode:      filterOpCode,
		Tag:         "tag3",
		Utc:         "+03:00",
	}
	client4 = dto.Client{
		PhoneNumber: "79041237777",
		OpCode:      "904",
		Tag:         filterTag,
		Utc:         "+03:00",
	}
	notExistingPhoneNumber = "70000000000"
	clientIds              [4]string
)

func defaultClientTestDBFill() {
	res, _ := ClientRepo.collection.InsertMany(context.TODO(), []interface{}{
		client1,
		client2,
		client3,
		client4,
	})
	for i, id := range res.InsertedIDs {
		clientIds[i] = id.(primitive.ObjectID).Hex()
	}
}

func beforeClientRepoTest() {
	utils.PMan = utils.NewPman("test.properties")
	ConnectToMongo(context.Background(), "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	DropDb()
}

func TestIsClientExist(t *testing.T) {
	beforeClientRepoTest()
	defaultClientTestDBFill()
	if !ClientRepo.IsClientExist(client1.PhoneNumber) {
		t.Errorf("did not find existing client")
	}
	if ClientRepo.IsClientExist(notExistingPhoneNumber) {
		t.Errorf("found not existing client")
	}
}

func TestFindExistingClient(t *testing.T) {
	beforeClientRepoTest()
	defaultClientTestDBFill()
	if *ClientRepo.FindClient(client1PhoneNumber) != client1 {
		t.Errorf("did not find existing client")
	}
}

func TestFindNotExistingClient(t *testing.T) {
	beforeClientRepoTest()
	defaultClientTestDBFill()
	if ClientRepo.FindClient(notExistingPhoneNumber) != nil {
		t.Errorf("found not existing client")
	}
}

func TestFindClientsByFilter(t *testing.T) {
	beforeClientRepoTest()
	defaultClientTestDBFill()
	filter := dto.Filter{
		OpCode: filterOpCode,
		Tag:    filterTag,
	}
	clients := ClientRepo.FindClientsByFilter(&filter)
	if len(*clients) != 1 && (*clients)[0] != *client2.WithId(clientIds[1]) {
		t.Errorf("wrong find by filter result")
	}
}

func TestCreateClient(t *testing.T) {
	beforeClientRepoTest()
	hexId := ClientRepo.CreateClient(&client1)
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		t.Errorf("wrong format of returned id")
		return
	}
	var client dto.Client
	cur, err := ClientRepo.collection.Find(context.TODO(), bson.D{
		{Key: "_id", Value: id},
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !cur.Next(context.TODO()) {
		t.Errorf("document does not inserted")
		return
	}
	if cur.Decode(&client) != nil {
		t.Errorf("decode error")
		return
	}
	if client != client1 {
		t.Errorf("document info corrupted")
		return
	}
	if cur.Next(context.TODO()) {
		t.Errorf("not unique id")
	}
}

func TestDeleteClient(t *testing.T) {
	beforeClientRepoTest()
	defaultClientTestDBFill()
	ClientRepo.DeleteClient(client1PhoneNumber)
	res, _ := ClientRepo.collection.Find(context.TODO(), bson.D{{}})
	var client dto.Client
	numDocs := 0
	for ; res.Next(context.TODO()); numDocs++ {
		res.Decode(&client)
		if client.PhoneNumber == client1PhoneNumber {
			t.Errorf("client was not deleted")
			return
		}
	}
	if numDocs != len(clientIds)-1 {
		t.Errorf("wrong client deleted")
	}
}

func TestUpdateClient(t *testing.T) {
	beforeClientRepoTest()
	res, _ := ClientRepo.collection.InsertOne(context.TODO(), client1)
	hexId := res.InsertedID.(primitive.ObjectID).Hex()
	ClientRepo.UpdateClient(hexId, &client2)
	cur, _ := ClientRepo.collection.Find(context.TODO(), bson.D{{}})
	var client dto.Client
	if cur.Next(context.TODO()) {
		cur.Decode(&client)
		if client != client2 {
			t.Errorf("invalid update")
		}
	} else {
		t.Errorf("id was update (should not)")
	}
}
