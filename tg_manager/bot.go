package tg_manager

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"telegram_bot_service/common"
	"telegram_bot_service/mongo_manager"
	"time"
)

func Bot() {
	b, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		//URL: "http://195.129.111.17:8012",

		Token:  "1096211287:AAEq9jysmBjjbZRIJxK6zxQ2-j2pizgRNbk",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	var (
		// Universal markup builders.
		menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		selector = &tb.ReplyMarkup{}
	)
	var (
		btn1 = selector.Data("1", "1", "1")
		btn2 = selector.Data("2", "2", "2")
		btn3 = selector.Data("3", "3", "3")
		btn4 = selector.Data("4", "4", "4")
	)

	var (
		btnDevices  = menu.Text("Devices")
		btnSettings = menu.Text("⚙ Settings")
	)

	menu.Reply(
		menu.Row(btnDevices),
		menu.Row(btnSettings),
	)

	selector.Inline(
		selector.Row(btn1, btn2, btn3, btn4),
		//selector.Row(btn10, btn11),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}

		b.Send(m.Sender, "Hello! Write token pls)")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		_, ok := common.Tg_users_by_chat_id[b.Me.ID]
		if !ok {
			if mongo_manager.Get_user_by_id(m.Text, b.Me.ID) {
				b.Send(m.Sender, "welcome", menu)
			} else {
				b.Send(m.Sender, "error")
			}
			return
		}
		b.Send(m.Sender, "Hello!", menu)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		b.Send(m.Sender, m.Photo)
	})

	b.Handle(&btnDevices, func(m *tb.Message) {
		user, ok := common.Tg_users_by_chat_id[b.Me.ID]
		if ok {
			b.Send(m.Sender, get_string_device_data(user), selector)
		}
	})

	b.Handle(&btn1, func(c *tb.Callback) {
		b.Send(c.Sender, "Answer")
	})
	b.Handle(&btn2, func(c *tb.Callback) {
		b.Send(c.Sender, "Answer")
	})
	b.Handle(&btn3, func(c *tb.Callback) {
		b.Send(c.Sender, "Answer")
	})
	b.Handle(&btn4, func(c *tb.Callback) {
		b.Send(c.Sender, "Answer")
	})

	b.Start()
}

func get_string_device_data(user common.Tg_user) string {
	answ := "Devices:\r\n\r\n"
	for i, device := range user.Devices {
		dev_data := ""
		dev_data += strconv.Itoa(i+1) + "):\n\r" +
			"\tName: " + device.Name + "\n\r" +
			"\tType: " + device.Type + "\r\n" +
			"\t" + get_value(device.Id) + "\r\n" +
			"\r\n======================\r\n"
		answ += dev_data
	}
	return answ
}

func get_value(dev_id int) string {
	dev_data, _ := mongo_manager.Get_last_dev_data(dev_id)
	answ := "Last message time: " + time.Unix(dev_data.Rt/1000000000, 0).String() + "\r\n" +
		"\tData: " + fmt.Sprintf("%v", dev_data.Data)

	return answ
}
