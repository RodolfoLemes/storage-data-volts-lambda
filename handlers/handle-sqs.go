package handlers

import (
	"context"
	"log"
	"os"
	"storage-data-volts-lambda/datavolts"
	"storage-data-volts-lambda/signals"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	AWS_ACCESS_KEY_ID     = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
)

var sess *session.Session
var svc *dynamodb.DynamoDB

func init() {
	/* sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})) */

	sess = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	}))

	svc = dynamodb.New(sess)
}

func HandleSQS(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {

		deviceTimestamp := record.MessageAttributes["timestamp"].StringValue
		messageId := record.MessageId
		queueTimestampString := record.Attributes["SentTimestamp"]

		dataVolts := datavolts.New(*deviceTimestamp, queueTimestampString, messageId)

		rTensions := strings.Split(*record.MessageAttributes["rTensions"].StringValue, ",")
		sTensions := strings.Split(*record.MessageAttributes["sTensions"].StringValue, ",")
		tTensions := strings.Split(*record.MessageAttributes["tTensions"].StringValue, ",")
		rCurrents := strings.Split(*record.MessageAttributes["rCurrents"].StringValue, ",")
		sCurrents := strings.Split(*record.MessageAttributes["sCurrents"].StringValue, ",")
		tCurrents := strings.Split(*record.MessageAttributes["tCurrents"].StringValue, ",")

		dataVolts.AddTensions(rTensions, "r")
		dataVolts.AddTensions(sTensions, "s")
		dataVolts.AddTensions(tTensions, "t")
		dataVolts.AddCurrents(rCurrents, "r")
		dataVolts.AddCurrents(sCurrents, "s")
		dataVolts.AddCurrents(tCurrents, "t")

		dynamoItem, err := dynamodbattribute.MarshalMap(dataVolts)

		if err != nil {
			log.Println("dynamo item err: ", err)
			continue
		}

		input := &dynamodb.PutItemInput{
			Item:      dynamoItem,
			TableName: aws.String(datavolts.DataVoltsTableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			log.Println("put item err: ", err)
			continue
		}

		log.Println("Successfully with the designed data as: ", dataVolts)
	}

	return nil
}

func HandleManually() {
	dataVolts := datavolts.DataVolts{
		Timestamp: time.Now().UTC(),
		RTensions: signals.BuildSin(5, 60, 0),
		STensions: signals.BuildSin(5, 60, 120),
		TTensions: signals.BuildSin(5, 60, 240),
		RCurrents: signals.BuildSin(1, 60, 0),
		SCurrents: signals.BuildSin(1, 60, 120),
		TCurrents: signals.BuildSin(1, 60, 240),
	}

	dynamoItem, err := dynamodbattribute.MarshalMap(dataVolts)

	if err != nil {
		log.Println("dynamo item err: ", err)
		return
	}

	input := &dynamodb.PutItemInput{
		Item:      dynamoItem,
		TableName: aws.String(datavolts.DataVoltsTableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Println("put item err: ", err)
		return
	}

	log.Println("Successfully with the designed data as: ", dataVolts)
}
