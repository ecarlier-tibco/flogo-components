// Package dynamodbinsert inserts an object into Amazon DynamoDB
package dynamodbinsert

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	/*
		defer func() {
			if r := recover(); r != nil {
				t.Failed()
				t.Errorf("panic during execution: %v", r)
			}
		}()*/

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	// To test this example, you can create a dynamodb with table name Music
	// where the key is called Artist
	tc.SetInput("awsAccessKeyID", "")
	tc.SetInput("awsSecretAccessKey", "")
	tc.SetInput("awsRegion", "eu-west-1")
	tc.SetInput("DynamoDBTableName", "Device")

	var r map[string]interface{}
	r = make(map[string]interface{})
	r["ID"] = "Test66"
	r["Chain_Address"] = "mc_address_6666"
	r["Latitude"] = 42.76
	r["Longitude"] = 4.55
	r["MeasureTypes"] = [2]string{"temperature", "humidity"}

	// b, _ := json.Marshal(payload)

	// tc.SetInput("DynamoDBRecord", string(b))
	tc.SetInput("DynamoDBRecord", r)
	act.Eval(tc)

	//check result attr
	result := tc.GetOutput("result")
	fmt.Printf("The Result of the insert was:\n[%s]\n", result)
}
