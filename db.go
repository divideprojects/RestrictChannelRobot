package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogSettings struct {
	ChatId       int64 `bson:"_id,omitempty" json:"_id,omitempty"`
	LogChannelID int64 `bson:"log_channel" json:"log_channel"`
}

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
	ignoreListCollection  *mongo.Collection
	logSettingsCollection *mongo.Collection
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
	logSettingsCollection = mongoClient.Database(databaseName).Collection("log_settings")
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

func countDocs(collecion *mongo.Collection, filter bson.M) (count int64, err error) {
	count, err = collecion.CountDocuments(tdCtx, filter)
	if err != nil {
		log.Errorf("[Database][countDocs]: %v", err)
	}
	return
}

func findAll(collecion *mongo.Collection, filter bson.M) (cur *mongo.Cursor) {
	cur, err := collecion.Find(tdCtx, filter)
	if err != nil {
		log.Errorf("[Database][findAll]: %v", err)
	}
	return
}

func deleteOne(collecion *mongo.Collection, filter bson.M) (err error) {
	_, err = collecion.DeleteOne(tdCtx, filter)
	if err != nil {
		log.Errorf("[Database][deleteOne]: %v", err)
	}
	return
}

func deleteMany(collecion *mongo.Collection, filter bson.M) (err error) {
	_, err = collecion.DeleteMany(tdCtx, filter)
	if err != nil {
		log.Errorf("[Database][deleteMany]: %v", err)
	}
	return
}

// GetChatSettings Get admin settings for a chat
func getLogSettings(chatID int64) *LogSettings {
	return _checkLogSetting(chatID)
}

// check Chat Admin Settings, used to get data before performing any operation
func _checkLogSetting(chatID int64) (adminSrc *LogSettings) {
	dLogSrc := &LogSettings{ChatId: chatID, LogChannelID: 0}

	err := findOne(logSettingsCollection, bson.M{"_id": chatID}).Decode(&adminSrc)
	if err == mongo.ErrNoDocuments {
		adminSrc = dLogSrc
		err := updateOne(logSettingsCollection, bson.M{"_id": chatID}, dLogSrc)
		if err != nil {
			log.Errorf("[Database][checkChatSetting]: %v ", err)
		}
	} else if err != nil {
		adminSrc = dLogSrc
		log.Errorf("[Database][checkChatSetting]: %v ", err)
	}
	return adminSrc
}

// SetAnonAdminMode Set anon admin mode for a chat
func setLogChannelID(chatId int64, logChannelId int64) {
	dLogSrc := _checkLogSetting(chatId)
	dLogSrc.LogChannelID = logChannelId

	err := updateOne(logSettingsCollection, bson.M{"_id": chatId}, dLogSrc)
	if err != nil {
		log.Errorf("[Database] SetLogChannelID: %v - %d", err, chatId)
	}
}

func getIgnoreSettings(chatID int64) *IgnoreList {
	return _checkIgnoreSettings(chatID)
}

// check Chat Approval Settings, used to get data before performing any operation
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
