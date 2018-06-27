/*
Package dynamodbquery queries objects from Amazon DynamoDB


To be able to test this package you'll need to have access to Amazon DynamoDB. A sample table could be a table with the
name *data* and with *itemtype* as the partition key, and *itemid* as the sort key (both could be strings). Some sample
data (which can be generated with Mockaroo) can be

```
{
  "firstname": "John",
  "itemid": "57a98d98e4b00679b4a830af",
  "itemtype": "user",
  "lastname": "Doe",
  "password": "fec51acb3365747fc61247da5e249674cf8463c2",
  "username": "Jon_Doe"
}

{
  "firstname": "User",
  "itemid": "57a98d98e4b00679b4a830b2",
  "itemtype": "user",
  "lastname": "Name",
  "password": "e2de7202bb2201842d041f6de201b10438369fb8",
  "username": "user"
}

{
  "firstname": "Admin",
  "itemid": "57a98d98e4b00679b4a830b5",
  "itemtype": "admin",
  "lastname": "Name1",
  "password": "8f31df4dcc25694aeb0c212118ae37bbd6e47bcd",
  "username": "admin"
}
```

With this data you can test the KeyConditionExpression *itemtype = user* and a more complex
KeyConditionExpression *itemtype = user and itemid = 57a98d98e4b00679b4a830af*. The former will
return 2 objects, the latter only the John Doe record.
*/
package dynamodbquery

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

// Update these variables before testing to match your own AWS account
const (
	awsAccessKeyID     = ""
	awsSecretAccessKey = ""
	awsDefaultRegion   = "eu-west-1"
	dynamoDBTableName  = "Device"
)

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

// Test for a single condition string with no filtering
func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	// Set required attributes
	tc.SetInput("awsAccessKeyID", awsAccessKeyID)
	tc.SetInput("awsSecretAccessKey", awsSecretAccessKey)
	tc.SetInput("awsRegion", awsDefaultRegion)
	tc.SetInput("dynamoDBTableName", dynamoDBTableName)
	tc.SetInput("dynamoDBKeyConditionExpression", "ID = :id")

	// Prepare the Key Condition Expression as Name/Value pairs
	var a map[string]interface{}
	a = make(map[string]interface{})
	a[":id"] = "Test66"

	// Execute the activity
	tc.SetInput("dynamoDBExpressionAttributes", a)
	act.Eval(tc)

	// Check the result
	printOutput(tc)
}

func TestEvalMultiFields(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	// Set required attributes
	tc.SetInput("awsAccessKeyID", awsAccessKeyID)
	tc.SetInput("awsSecretAccessKey", awsSecretAccessKey)
	tc.SetInput("awsRegion", awsDefaultRegion)
	tc.SetInput("dynamoDBTableName", dynamoDBTableName)
	tc.SetInput("dynamoDBKeyConditionExpression", "Customer = :cust and Site = :site")
	tc.SetInput("dynamoDBIndexName", "Customer-Site-index")

	// Prepare the Key Condition Expression as Name/Value pairs
	var a map[string]interface{}
	a = make(map[string]interface{})
	a[":cust"] = "Acme Corp"
	a[":site"] = "Bordeaux"

	// Execute the activity
	tc.SetInput("dynamoDBExpressionAttributes", a)
	act.Eval(tc)

	// Check the result
	printOutput(tc)
}

func printOutput(tc *test.TestActivityContext) {
	result := tc.GetOutput("result")
	scannedCount := tc.GetOutput("scannedCount")
	consumedCapacity := tc.GetOutput("consumedCapacity")
	fmt.Printf("The ScannedCount of the query was: [%s]\n", scannedCount)
	fmt.Printf("The ConsumedCapacity of the query was: [%v]\n", consumedCapacity)
	fmt.Printf("The Result of the query was [%v]\n", result)

}
