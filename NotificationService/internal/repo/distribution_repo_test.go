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
	distributions = []dto.Distribution{
		{
			StartAt: "2001-01-01T00:00:00",
			Message: "message1",
			Filter: dto.Filter{
				OpCode: "opCode1",
				Tag:    "tag1",
			},
			EndAt: "2001-01-01T01:00:00",
		},
		{
			StartAt: "2002-01-01T00:00:00",
			Message: "message2",
			Filter: dto.Filter{
				OpCode: "opCode2",
				Tag:    "tag2",
			},
			EndAt: "2002-01-01T01:00:00",
		},
		{
			StartAt: "2003-01-01T00:00:00",
			Message: "message3",
			Filter: dto.Filter{
				OpCode: "opCode3",
				Tag:    "tag3",
			},
			EndAt: "2003-01-01T01:00:00",
		},
		{
			StartAt: "2004-01-01T00:00:00",
			Message: "message4",
			Filter: dto.Filter{
				OpCode: "opCode4",
				Tag:    "tag4",
			},
			EndAt: "2004-01-01T01:00:00",
		},
	}
	distributionIds           []string
	notExistingDistributionId = "aaaaaaaaaaaaaaaaaaaaaaaa"
)

func defaultDistributionRepoTestDBFill() {
	var insertVals []interface{}
	for _, distribution := range distributions {
		insertVals = append(insertVals, distribution)
	}
	res, _ := DistributionRepo.collection.InsertMany(context.TODO(), insertVals)
	distributionIds = make([]string, 0)
	for _, id := range res.InsertedIDs {
		distributionIds = append(distributionIds, id.(primitive.ObjectID).Hex())
	}

}

func beforeDistributionRepoTest() {
	utils.PMan = utils.NewPman("test.properties")
	ConnectToMongo(context.Background(), "mongodb://"+utils.PMan.Get("mongo_host").(string)+":"+utils.PMan.Get("mongo_port").(string), utils.PMan.Get("mongo_db_name").(string))
	DropDb()
}

func TestIsDistributionExist(t *testing.T) {
	beforeDistributionRepoTest()
	defaultDistributionRepoTestDBFill()
	if !DistributionRepo.IsDistributionExist(distributionIds[0]) {
		t.Errorf("did not find existing distribution")
	}
	if DistributionRepo.IsDistributionExist(notExistingDistributionId) {
		t.Errorf("found not existing distribution")
	}
}

func TestFindExistingDistribution(t *testing.T) {
	beforeDistributionRepoTest()
	defaultDistributionRepoTestDBFill()
	for i, id := range distributionIds {
		if *DistributionRepo.FindDistributionById(id) != *distributions[i].WithId(id) {
			t.Errorf("did not find existing client")
		}
	}
}

func TestFindNotExistingDistribution(t *testing.T) {
	beforeDistributionRepoTest()
	defaultDistributionRepoTestDBFill()
	if DistributionRepo.FindDistributionById(notExistingDistributionId) != nil {
		t.Errorf("found not existing distribution")
	}
}

func TestFindAllDistributions(t *testing.T) {
	beforeDistributionRepoTest()
	defaultDistributionRepoTestDBFill()
	res := DistributionRepo.FindAllDistributions()
	if len(*res) != len(distributions) {
		t.Errorf("wrong num of found distributions")
		return
	}
	for _, distribution := range *res {
		foundF := false
		for i, testDistribution := range distributions {
			foundF = (*testDistribution.WithId(distributionIds[i]) == distribution)
			if foundF {
				break
			}
		}
		if !foundF {
			t.Errorf("existing distribution not found")
		}
	}
}

func TestCreateDistribution(t *testing.T) {
	beforeDistributionRepoTest()
	hexId := DistributionRepo.CreateDistribution(&distributions[0])
	id, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		t.Errorf("wrong format of returned id")
		return
	}
	var distribution dto.Distribution
	cur, err := DistributionRepo.collection.Find(context.TODO(), bson.D{
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
	if cur.Decode(&distribution) != nil {
		t.Errorf("decode error")
		return
	}
	if distribution != distributions[0] {
		t.Errorf("document info corrupted")
		return
	}
	if cur.Next(context.TODO()) {
		t.Errorf("not unique id")
	}
}

func TestDeleteDistribution(t *testing.T) {
	beforeDistributionRepoTest()
	defaultDistributionRepoTestDBFill()
	DistributionRepo.DeleteDistribution(distributionIds[0])
	res, _ := DistributionRepo.collection.Find(context.TODO(), bson.D{{}})
	var distribution dto.Distribution
	numDocs := 0
	for ; res.Next(context.TODO()); numDocs++ {
		res.Decode(&distribution)
		if distribution == distributions[0] {
			t.Errorf("distribution was not deleted")
			return
		}
	}
	if numDocs != len(distributionIds)-1 {
		t.Errorf("wrong distribution deleted")
	}
}

func TestUpdateDistribution(t *testing.T) {
	beforeDistributionRepoTest()
	res, _ := DistributionRepo.collection.InsertOne(context.TODO(), distributions[0])
	hexId := res.InsertedID.(primitive.ObjectID).Hex()
	DistributionRepo.UpdateDistribution(hexId, &distributions[1])
	cur, _ := DistributionRepo.collection.Find(context.TODO(), bson.D{{}})
	var distribution dto.Distribution
	if cur.Next(context.TODO()) {
		cur.Decode(&distribution)
		if distribution != distributions[1] {
			t.Errorf("invalid update")
		}
	} else {
		t.Errorf("id was update (should not)")
	}
}
