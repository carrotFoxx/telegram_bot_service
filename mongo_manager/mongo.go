package mongo_manager

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"telegram_bot_service/common"
)

var mng *mongo.Client
var clientQerry *mongo.Client
var user_collection *mongo.Collection
var devs_collection *mongo.Collection

func MongoIni(uri string) {
	common.Tg_users_by_chat_id = make(map[int]common.Tg_user)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	client2, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	mng = client
	clientQerry = client2

	user_collection = clientQerry.Database("users").Collection("tg_users")
	devs_collection = clientQerry.Database("dev_data").Collection("150")
	fill_user_map_by_db()
}

func fill_user_map_by_db() {
	cur, err := user_collection.Find(context.TODO(), bson.D{{}}) //, findOptions)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var user common.Tg_user
		err := cur.Decode(&user)
		if err != nil {
			log.Println(err)
		}
		if user.Chat_id != 0 {
			common.Tg_users_by_chat_id[user.Chat_id] = user
		}
	}
}

func Set_user_by_token(token string) {
	//user_collection.Find()
}

func Get_user_by_id(id string, chat_id int) bool {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	doc := user_collection.FindOne(context.TODO(), bson.M{"_id": objectId})
	var user common.Tg_user
	if doc != nil {
		err := doc.Decode(&user)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	user.Chat_id = chat_id
	filter := bson.M{"_id": bson.M{"$eq": objectId}}
	update := bson.M{
		"$set": bson.M{
			"Chat_id": chat_id,
		},
	}
	_, err = user_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return false
	}
	fill_user_map_by_db()
	return true
}

func Get_last_dev_data(id int) (common.Dev_data, error) {
	filter := bson.M{"id": bson.M{"$eq": id}}
	//filter := bson.M{"id": bson.M{"$eq": id},"_id":"$max"}
	data := common.Dev_data{}
	opt := options.Find()
	opt.SetLimit(1)
	opt.SetSort(bson.M{"_id": -1})
	curs, err := devs_collection.Find(context.TODO(), filter, opt)
	if err != nil {
		log.Println(err)
		return data, err
	}
	for curs.Next(context.TODO()) {
		err = curs.Decode(&data)
	}
	if err == nil {
		return data, err
	}
	return data, nil
}
