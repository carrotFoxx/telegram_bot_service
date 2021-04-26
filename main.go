package main

import (
	"telegram_bot_service/mongo_manager"
	"telegram_bot_service/tg_manager"
)

func main(){
	mongo_manager.MongoIni("mongodb://localhost:27017")
	tg_manager.Bot()
}
