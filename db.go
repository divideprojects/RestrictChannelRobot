package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IgnoreList struct {
	ChatId          int64   `bson:"_id,omitempty" json:"_id,omitempty"`
	IgnoredChannels []int64 `bson:"ignored_channels" json:"ignored_channels"`
}

var (
	mongoClient *mongo.Client

	// Contexts
	tdCtx = context.TODO()
	bgCtx = context.Background()

	// define collections
	ignoreListCollection *mongo.Collection
)

func init() {
	mongoClient, err := mongo.NewClient(
		options.Client().ApplyURI(databaseUrl),
	)
	if err != nil {
		log.Errorf("[Database][Client]: %v", err)
	}

	ctx, cancel := context.WithTimeout(bgCtx, 10*time.Second)
	defer cancel()

	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Errorf("[Database][Connect]: %v", err)
	}

	// Open Connections to Collections
	log.Info("Opening Database Collections...")
	ignoreListCollection = mongoClient.Database(databaseName).Collection("ignore_list")
}

func updateOne(collecion *mongo.Collection, filter bson.M, data interface{}) (err error) {
	_, err = collecion.UpdateOne(tdCtx, filter, bson.M{"$set": data}, options.Update().SetUpsert(true))
	if err != nil {
		log.Errorf("[Database][updateOne]: %v", err)
	}
	return
}

func findOne(collecion *mongo.Collection, filter bson.M) (res *mongo.SingleResult) {
	res = collecion.FindOne(tdCtx, filter)
	return
}

func getIgnoreSettings(chatID int64) *IgnoreList {
	return _checkIgnoreSettings(chatID)
}

func _checkIgnoreSettings(chatID int64) (ignorerc *IgnoreList) {
	defaultIgnoreSettings := &IgnoreList{ChatId: chatID, IgnoredChannels: make([]int64, 0)}

	errS := findOne(ignoreListCollection, bson.M{"_id": chatID}).Decode(&ignorerc)
	if errS == mongo.ErrNoDocuments {
		ignorerc = defaultIgnoreSettings
		err := updateOne(ignoreListCollection, bson.M{"_id": chatID}, defaultIgnoreSettings)
		if err != nil {
			log.Errorf("[Database][_checkIgnoreSettings]: %v ", err)
		}
	} else if errS != nil {
		log.Errorf("[Database][_checkIgnoreSettings]: %v", errS)
		ignorerc = defaultIgnoreSettings
	}
	return ignorerc
}

func ignoreChat(chatID, ignoreChannelId int64) {
	ignorerc := _checkIgnoreSettings(chatID)
	ignorerc.IgnoredChannels = append(
		ignorerc.IgnoredChannels,
		ignoreChannelId,
	)
	err := updateOne(ignoreListCollection, bson.M{"_id": chatID}, ignorerc)
	if err != nil {
		log.Errorf("[Database] ignoreChat: %v", err)
	}
}

func unignoreChat(chatID, ignoreChannelId int64) {
	ignorerc := _checkIgnoreSettings(chatID)
	for i, v := range ignorerc.IgnoredChannels {
		if v == ignoreChannelId {
			ignorerc.IgnoredChannels = append(ignorerc.IgnoredChannels[:i], ignorerc.IgnoredChannels[i+1:]...)
			break
		}
	}
	err := updateOne(ignoreListCollection, bson.M{"_id": chatID}, ignorerc)
	if err != nil {
		log.Errorf("[Database] unignoreChat: %v", err)
	}
}
