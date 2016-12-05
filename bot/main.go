package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"os"
)

type JobResponse struct {
	Jobs   []Job  `json:"data"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type Job struct {
	Author   string `json:"author"`
	Finish   string `json:"finish"`
	Kind     string `json:"kind"`
	Result   string `json:"result"`
	Start    string `json:"start"`
	Status   string `json:"status"`
	Callback string `json:"callback"`
	Params   string `json:"params"`
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

	defaultKeyb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Status", "Status"),
			tgbotapi.NewInlineKeyboardButtonData("Features", "Features"),
			tgbotapi.NewInlineKeyboardButtonData("Jobs queue", "Jobs queue"),
		),
	)

	for update := range updates {

		if update.CallbackQuery != nil {
			var msg tgbotapi.MessageConfig
			chatId := update.CallbackQuery.Message.Chat.ID
			// bot.Send(tgbotapi.NewChatAction(chatId, tgbotapi.ChatTyping))
			switch update.CallbackQuery.Data {
			case "Status":
				text := getJsonMarkdown("status")
				msg = tgbotapi.NewMessage(chatId, text)
				msg.ReplyMarkup = &defaultKeyb
			case "Features":
				text := getJsonMarkdown("features")
				msg = tgbotapi.NewMessage(chatId, text)
				msg.ReplyMarkup = &defaultKeyb
			case "Jobs queue":
				resp, _ := http.Get("http://korchasa.host/api/v1/jobs_queue")
				defer resp.Body.Close()

				data, err := ioutil.ReadAll(resp.Body)

				var jobResp JobResponse
				if err == nil && data != nil {
					json.Unmarshal(data, &jobResp)
				}

				text := "Jobs:\n"
				for _, job := range jobResp.Jobs {
					if job.Kind == "" || job.Author == "" {
						continue
					}
					if job.Status == "finished" {
						text += "*" + job.Kind + "* at _" + job.Author + "_\n"
						text += "*Start/Finish*: " + job.Start + " / " + job.Finish + "\n"
						text += "*Result*:\n" + job.Result + "\n"
						text += "*Status*: " + job.Status + "\n"
					} else {
						text += "*" + job.Kind + "* at _" + job.Author + "_\n"
						text += "*Params*: " + job.Params + "\n"
						text += "*Status*: " + job.Status + "\n"
					}
					text += "\n\n"
				}
				log.Println(text)
				msg = tgbotapi.NewMessage(chatId, text)
				msg.ReplyMarkup = &defaultKeyb
			}
			msg.ParseMode = "markdown"
			_, err = bot.Send(msg)
			log.Printf("%#v", err)
		} else if update.Message != nil {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Welcome to korchasa bot. What do you want to know?",
			)
			jsonStr, _ := json.MarshalIndent(msg, "\t", "\t")
			log.Println(string(jsonStr))
			msg.ReplyMarkup = &defaultKeyb
			_, err := bot.Send(msg)
			log.Println(err)
		}


	}
}

func tree2plain(m map[string]interface{}, prefix string) string {
	text := ""
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			text += fmt.Sprintf("**%s**: %s", k, vv)
		case int:
			log.Println(k, "is int", vv)
		case []interface{}:
			log.Println(k, "is an array:")
			for i, u := range vv {
				log.Println(i, u)
			}
		default:
			log.Println(k, "is of a type I don't know how to handle")
		}
	}
	return text
}

func getJsonMarkdown(url string) string {
	resp, _ := http.Get("http://korchasa.host/api/v1/" + url)

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	text := string(bytes)
	// replace with strings.NewReplacer
	text = strings.Replace(text, "\\n", "\n", -1)
	text = strings.Replace(text, "{", "", -1)
	text = strings.Replace(text, "},\n", "", -1)
	text = strings.Replace(text, "}", "", -1)
	text = strings.Replace(text, "\"", "", -1)
	text = strings.Replace(text, "[", "", -1)
	text = strings.Replace(text, "],\n", "", -1)
	text = strings.Replace(text, "]", "", -1)
	text = strings.Replace(text, "error: ,", "", -1)
	text = strings.Replace(text, "status: 200", "", -1)
	text = strings.Replace(text, "data:", "", -1)
	text = strings.Replace(text, ",\n", "\n", -1)
	text = strings.Replace(text, "\n        ", "\n", -1)
	text = strings.Replace(text, "    ", "  ", -1)

	return string(text)
}
