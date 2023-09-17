package route

import (
	"NotificationService/internal/utils"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestBodyHandler struct {
	Int     int
	String  string
	Bool    bool
	MongoId primitive.ObjectID `json:"_id"`
}

func TestSuccessParseBody(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)

	mongoIdHex := "5f82c256db7a68e6bc123456"
	r, err := http.NewRequest(http.MethodGet, "", strings.NewReader(fmt.Sprintf(
		`{
			"Int":42,
			"String":"string",
			"Bool":true,
			"_id": {
				"$oid": "%s"
			}
		}`,
		mongoIdHex),
	))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	var testBodyHandler TestBodyHandler
	if !parseBody(w, r, &testBodyHandler) {
		t.Errorf("Wrong boolean return value")
	}

	mongoId, _ := primitive.ObjectIDFromHex(mongoIdHex)
	if !reflect.DeepEqual(testBodyHandler, TestBodyHandler{
		Int:     42,
		String:  "string",
		Bool:    true,
		MongoId: mongoId,
	}) {
		t.Errorf("Wrong body parsing")
	}
}

func TestFailedParseBodyUnmarshalErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusInternalServerError)

	r, err := http.NewRequest(http.MethodGet, "", strings.NewReader("non valid json string"))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	var testBodyHandler TestBodyHandler
	t.Log("Expected unmarshal error message\n")
	if parseBody(w, r, &testBodyHandler) {
		t.Errorf("Wrong boolean return value")
	}
}
