package main

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Call struct {
	ChatId  int64
	Action  func() tgbotapi.MessageConfig
	Payload map[string]interface{}
	Update  tgbotapi.Update
}

func (call *Call) DefaultKeyb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Status", "status%%{}"),
			tgbotapi.NewInlineKeyboardButtonData("Features", "features%%{}"),
			tgbotapi.NewInlineKeyboardButtonData("Jobs queue", "queue_list%%{}"),
		),
	)
}

func NewCall(update tgbotapi.Update) *Call {

	i := new(Call)
	i.Update = update

	if update.CallbackQuery != nil {
		i.ChatId = update.CallbackQuery.Message.Chat.ID
	} else if update.Message != nil {
		i.ChatId = update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {

		parts := strings.Split(update.CallbackQuery.Data, "%%")

		action, payloadJson := parts[0], parts[1]

		json.Unmarshal([]byte(payloadJson), &i.Payload)

		validActions := map[string]func() tgbotapi.MessageConfig{
			"status":     i.ActionStatus,
			"features":   i.ActionFeatures,
			"queue_list": i.ActionQueueList,
		}

		i.Action = validActions[action]
	} else {
		i.Action = i.ActionDefault
	}

	return i
}

func (call *Call) ActionDefault() tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		call.ChatId,
		"Welcome to korchasa bot. What do you want to know?",
	)
	keyb := call.DefaultKeyb()
	msg.ReplyMarkup = &keyb
	return msg
}

func (call *Call) ActionStatus() tgbotapi.MessageConfig {
	text := getAsText("status")
	msg := tgbotapi.NewMessage(call.ChatId, text)
	keyb := call.DefaultKeyb()
	msg.ReplyMarkup = &keyb
	return msg
}

func (call *Call) ActionFeatures() tgbotapi.MessageConfig {
	text := getAsText("features")
	msg := tgbotapi.NewMessage(call.ChatId, text)
	keyb := call.DefaultKeyb()
	msg.ReplyMarkup = &keyb
	return msg
}

func (call *Call) ActionQueueList() tgbotapi.MessageConfig {
	text := getAsText("jobs_queue")
	msg := tgbotapi.NewMessage(call.ChatId, text)
	keyb := call.DefaultKeyb()
	msg.ReplyMarkup = &keyb
	return msg
}

func getAsText(url string) string {
	resp, _ := http.Get("http://korchasa.host/api/v1/" + url)

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	var rawResp map[string]interface{}

	if err := json.Unmarshal(bytes, &rawResp); err != nil {
		return string(err.Error())
	}

	return tree2plain(rawResp["data"], "")
}

func tree2plain(m interface{}, prefix string) string {

	var text string

	switch vv := m.(type) {

	case string:
		text += prefix + vv + "\n"

	case int:
	case int64:
		text += prefix + strconv.FormatInt(vv, 10) + "\n"

	case float64:
		text += prefix + strconv.FormatFloat(vv, 'f', -1, 32) + "\n"

	case []interface{}:
		for _, u := range vv {
			text += prefix + tree2plain(u, prefix+"    ")
		}

	case map[string]interface{}:
		mk := make([]string, len(vv))
		i := 0
		for k, _ := range vv {
			mk[i] = k
			i++
		}
		sort.Strings(mk)

		text += "\n"
		for _, k := range mk {
			v := vv[k]
			text += prefix + "*" + k + "*:  " + tree2plain(v, prefix+"    ")
		}
	default:
		log.Printf("Type (%T) I don't know how to handle", vv)
	}

	return text
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		call := NewCall(update)
		msg := call.Action()
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
}
