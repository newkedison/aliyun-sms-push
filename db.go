package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbClient      *mongo.Client
	db            *mongo.Database
	colSendRecord *mongo.Collection
	colSignName   *mongo.Collection
	colPhoneList  *mongo.Collection
	colIpList     *mongo.Collection
)

type SendRecord struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	Domain          string
	SignName        string
	PhoneNumber     string
	TemplateCode    string
	TemplateParam   string
	ExtraInfo       string
	RequestId       string
	BizId           string
	ResponseCode    string
	ResponseMessage string
	SendStatus      string
	ErrCode         string
	Content         string
	SendDate        time.Time
	ReceiveDate     time.Time
	CreatedAt       time.Time
}

type allowedPhone struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Phone     string
	Remark    string
	CreatedAt time.Time
}

type allowedIp struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Ip        string
	Remark    string
	CreatedAt time.Time
}

func createUniqueIndex(ctx context.Context) {
	indexModel := mongo.IndexModel{
		bson.D{
			{"phone", 1},
		},
		options.Index().SetName("unique_phone").SetUnique(true),
	}
	_, err := colPhoneList.Indexes().CreateOne(
		ctx, indexModel, options.CreateIndexes())
	if err != nil {
		printError("CreateIndex fail", err)
	}

	indexModel = mongo.IndexModel{
		bson.D{
			{"ip", 1},
		},
		options.Index().SetName("unique_ip").SetUnique(true),
	}
	_, err = colIpList.Indexes().CreateOne(
		ctx, indexModel, options.CreateIndexes())
	if err != nil {
		printError("CreateIndex fail", err)
	}
}

func connectToDB(uri string) error {
	ctx := getContextWithTimeout(3000)
	option := options.Client().ApplyURI(uri)
	if globalConfig.MongoDB.User != "" {
		print("reset user name")
		option.Auth.Username = globalConfig.MongoDB.User
	}
	if globalConfig.MongoDB.Password != "" {
		print("reset password")
		option.Auth.Password = globalConfig.MongoDB.Password
	}
	client, err := mongo.NewClient(option)
	if err != nil {
		printError("NewClient fail", err)
		return err
	}
	err = client.Connect(ctx)
	if err != nil {
		printError("Connect fail", err)
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		printError("Ping fail", err)
		client.Disconnect(ctx)
		return err
	}
	_, err = client.ListDatabases(ctx, bson.D{})
	if err != nil {
		printError("ListDatabases fail", err)
		client.Disconnect(ctx)
		return err
	}
	dbClient = client
	db = client.Database("sms_push")
	colSendRecord = db.Collection("sms_record")
	colSignName = db.Collection("sms_sign_name")
	colPhoneList = db.Collection("sms_phone_list")
	colIpList = db.Collection("sms_ip_list")

	createUniqueIndex(ctx)

	return nil
}
