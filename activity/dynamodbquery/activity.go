// Package dynamodbquery queries objects from Amazon DynamoDB
package dynamodbquery

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// Constants used by the code to represent the input and outputs of the JSON structure
const (
	ivAwsAccessKeyID                 = "awsAccessKeyID"
	ivAwsSecretAccessKey             = "awsSecretAccessKey"
	ivAwsRegion                      = "awsRegion"
	ivDynamoDBTableName              = "dynamoDBTableName"
	ivDynamoDBKeyConditionExpression = "dynamoDBKeyConditionExpression"
	ivDynamoDBExpressionAttributes   = "dynamoDBExpressionAttributes"
	ivDynamoDBFilterExpression       = "dynamoDBFilterExpression"
	ivDynamoDBIndexName              = "dynamoDBIndexName"

	ovResult           = "result"
	ovScannedCount     = "scannedCount"
	ovConsumedCapacity = "consumedCapacity"
)

// log is the default package logger
var log = logger.GetLogger("activity-dynamodbquery")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// Get the inputs
	awsRegion := context.GetInput(ivAwsRegion).(string)
	dynamoDBTableName := context.GetInput(ivDynamoDBTableName).(string)
	dynamoDBKeyConditionExpression := context.GetInput(ivDynamoDBKeyConditionExpression).(string)
	dynamoDBExpressionAttributes := context.GetInput(ivDynamoDBExpressionAttributes).(map[string]interface{})
	dynamoDBFilterExpression := context.GetInput(ivDynamoDBFilterExpression).(string)
	dynamoDBIndexName := context.GetInput(ivDynamoDBIndexName).(string)

	// AWS Credentials, only if needed
	var awsAccessKeyID, awsSecretAccessKey = "", ""
	if context.GetInput(ivAwsAccessKeyID) != nil {
		awsAccessKeyID = context.GetInput(ivAwsAccessKeyID).(string)
	}
	if context.GetInput(ivAwsSecretAccessKey) != nil {
		awsSecretAccessKey = context.GetInput(ivAwsSecretAccessKey).(string)
	}

	// Create a session with Credentials only if they are set
	var awsSession *session.Session
	if awsAccessKeyID != "" && awsSecretAccessKey != "" {
		// Create new credentials using the accessKey and secretKey
		awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

		// Create a new session with AWS credentials
		awsSession = session.Must(session.NewSession(&aws.Config{
			Credentials: awsCredentials,
			Region:      aws.String(awsRegion),
		}))
	} else {
		// Create a new session without AWS credentials
		awsSession = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(awsRegion),
		}))
	}

	// Create a new login to the DynamoDB service
	dynamoService := dynamodb.New(awsSession)

	expressionAttributeMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range dynamoDBExpressionAttributes {
		dav, err := dynamodbattribute.Marshal(v)

		if err != nil {
			log.Errorf("DynamoDB Marshal Error [%v]", err)
			return false, err
		}
		expressionAttributeMap[k] = dav
	}

	// Construct the DynamoDB query
	var queryInput = &dynamodb.QueryInput{}
	if dynamoDBFilterExpression == "" {
		if dynamoDBIndexName == "" {
			queryInput = &dynamodb.QueryInput{
				TableName:                 aws.String(dynamoDBTableName),
				KeyConditionExpression:    aws.String(dynamoDBKeyConditionExpression),
				ExpressionAttributeValues: expressionAttributeMap,
				ReturnConsumedCapacity:    aws.String("TOTAL"),
			}
		} else {
			queryInput = &dynamodb.QueryInput{
				TableName:                 aws.String(dynamoDBTableName),
				IndexName:                 aws.String(dynamoDBIndexName),
				KeyConditionExpression:    aws.String(dynamoDBKeyConditionExpression),
				ExpressionAttributeValues: expressionAttributeMap,
				ReturnConsumedCapacity:    aws.String("TOTAL"),
			}
		}
	} else {
		if dynamoDBIndexName == "" {
			queryInput = &dynamodb.QueryInput{
				TableName:                 aws.String(dynamoDBTableName),
				KeyConditionExpression:    aws.String(dynamoDBKeyConditionExpression),
				ExpressionAttributeValues: expressionAttributeMap,
				FilterExpression:          aws.String(dynamoDBFilterExpression),
				ReturnConsumedCapacity:    aws.String("TOTAL"),
			}
		} else {
			queryInput = &dynamodb.QueryInput{
				TableName:                 aws.String(dynamoDBTableName),
				IndexName:                 aws.String(dynamoDBIndexName),
				KeyConditionExpression:    aws.String(dynamoDBKeyConditionExpression),
				ExpressionAttributeValues: expressionAttributeMap,
				FilterExpression:          aws.String(dynamoDBFilterExpression),
				ReturnConsumedCapacity:    aws.String("TOTAL"),
			}
		}
	}

	// Prepare and execute the DynamoDB query
	var queryOutput, err1 = dynamoService.Query(queryInput)
	if err1 != nil {
		log.Errorf("Error while executing query [%s]", err1)
		return false, err1
	}

	var result []map[string]interface{}
	err = dynamodbattribute.UnmarshalListOfMaps(queryOutput.Items, &result)
	if err != nil {
		log.Errorf("DYNAMO DB Unmarshall results error [%v]", err)
		return false, err
	}
	// Set the output value in the context
	context.SetOutput(ovResult, result)
	sc := *queryOutput.ScannedCount
	context.SetOutput(ovScannedCount, sc)
	cc := *queryOutput.ConsumedCapacity.CapacityUnits
	context.SetOutput(ovConsumedCapacity, cc)

	// Complete the activity
	return true, nil
}
