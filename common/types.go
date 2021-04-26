package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type Device struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Id int `json:"id"`
}

type Tg_user struct{
	Id primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username string
	Devices []Device
	Chat_id int
}
var Tg_users_by_chat_id map[int]Tg_user

type Dev_data struct{
	Rt int64 `bson:"_id"`
	Id int `json:"id"`
	Type int `json:"type"`
	Data interface{} `bson:"d"`
}

