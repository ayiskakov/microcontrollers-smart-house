package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	OpenDoor          = "Open door"
	CloseDoor         = "Close door"
	TurnOnLed         = "Turn on LED"
	TurnOffLed        = "Turn off LED"
	TurnOnSecureMode  = "Turn on secure mode"
	TurnOffSecureMode = "Turn off secure mode"
	GetStats          = "Get statistical data"
	StartMessage      = "Hey!\n\nIt is Home Automation Bot.\n\nI know next commands:\n\n/health - to check health of api\n/client_id - to get your id"
	ApiUrl            = "http://localhost:8000/api"
)

type Home struct {
	ID           string `json:"id,omitempty"`
	ClientId     string `json:"client_id"`
	Temperature  string `json:"temperature"`
	IsGateOpened bool   `json:"is_gate_opened"`
	IsRobbery    bool   `json:"is_robbery"`
	IsLedTurned  bool   `json:"is_led_turned"`
	SecureMode   bool   `json:"secure_mode"`
}
type UpdateHomeCommandInput struct {
	NewClientId *string `json:"new_client_id,omitempty""`
	SecureMode  *bool   `json:"secure_mode,omitempty"`
	OpenGate    *bool   `json:"open_gate,omitempty"`
	LedTurn     *bool   `json:"turn_led,omitempty"`
}

func renderKeyboard(homeState *HomeState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(homeState.DoorState),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(homeState.LedState),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(homeState.SecurityState),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(GetStats),
		),
	)
}

type HomeState struct {
	LedState      string
	SecurityState string
	DoorState     string
}

type DB struct {
	homeStates sync.Map
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	db := &DB{}
	bot, err := tgbotapi.NewBotAPI("YOUR BOT API TOKEN")

	if err != nil {
		logrus.Fatalf("error occured while connecting to bot: %s", err.Error())
	}

	bot.Debug = false

	logrus.Print("Bot connected: @", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go UpdateHandler(updates, db, bot)

	logrus.Print("Home Automation Bot Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Home Automation Bot Shutting Down")
}

func UpdateHandler(updates tgbotapi.UpdatesChannel, db *DB, bot *tgbotapi.BotAPI) {
	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		_, ok := db.homeStates.Load(msg.ChatID)

		if !ok {
			db.homeStates.Store(msg.ChatID, &HomeState{
				DoorState:     OpenDoor,
				LedState:      TurnOnLed,
				SecurityState: TurnOnSecureMode,
			})
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = StartMessage
			case "health":
				if !checkHealth() {
					msg.Text = "Service is not available"
				} else {
					msg.Text = "Service is healthy and running"
				}
			case "client_id":
				msg.Text = "Your client ID: " + strconv.FormatInt(msg.ChatID, 10)
			default:
				msg.Text = "Sorry, I don't know that command."
			}
			v, _ := db.homeStates.Load(msg.ChatID)
			var keyboard = renderKeyboard(v.(*HomeState))
			msg.ReplyMarkup = keyboard

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			continue
		}
		v, _ := db.homeStates.Load(msg.ChatID)

		switch update.Message.Text {
		case OpenDoor:
			u := UpdateHomeCommandInput{
				OpenGate: boolPointer(true),
			}
			updateHome(u, strconv.FormatInt(msg.ChatID, 10))
			msg.Text = "Door opened successfully"

			v.(*HomeState).DoorState = CloseDoor
		case CloseDoor:
			u := UpdateHomeCommandInput{
				OpenGate: boolPointer(false),
			}
			updateHome(u, strconv.FormatInt(msg.ChatID, 10))
			msg.Text = "Door closed successfully"

			v.(*HomeState).DoorState = OpenDoor
		case TurnOnLed:
			u := UpdateHomeCommandInput{
				LedTurn: boolPointer(true),
			}
			updateHome(u, strconv.FormatInt(msg.ChatID, 10))
			msg.Text = "LED turned ON successfully"
			v.(*HomeState).LedState = TurnOffLed
		case TurnOffLed:
			u := UpdateHomeCommandInput{
				LedTurn: boolPointer(false),
			}
			updateHome(u, strconv.FormatInt(msg.ChatID, 10))
			msg.Text = "LED turned OFF successfully"
			v.(*HomeState).LedState = TurnOnLed
		case TurnOnSecureMode:
			u := UpdateHomeCommandInput{
				SecureMode: boolPointer(true),
			}
			updateHome(u, strconv.FormatInt(msg.ChatID, 10))
			msg.Text = "Secure mode turned ON successfully"
			v.(*HomeState).SecurityState = TurnOffSecureMode
		case TurnOffSecureMode:
			u := UpdateHomeCommandInput{
				SecureMode: boolPointer(false),
			}
			updateHome(u, strconv.FormatInt(msg.ChatID, 10))
			msg.Text = "Secure mode turned OFF successfully"
			v.(*HomeState).SecurityState = TurnOnSecureMode
		case GetStats:
			msg.Text = getStats(strconv.FormatInt(msg.ChatID, 10))
		}
		msg.ReplyMarkup = renderKeyboard(v.(*HomeState))

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

		continue
	}
}

func getStats(chatId string) string {
	api := fmt.Sprintf("%s/home/telegram/%s", ApiUrl, chatId)
	resp, err := http.Get(api)
	defer resp.Body.Close()
	var h Home
	if err != nil {
		logrus.Fatal(err.Error())
	}
	err = json.NewDecoder(resp.Body).Decode(&h)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	data := fmt.Sprintf("Temperature: %s\n", h.Temperature)
	return data
}

func updateHome(u UpdateHomeCommandInput, chatId string) {
	api := fmt.Sprintf("%s/home/telegram/%s", ApiUrl, chatId)
	ct := "application/json"
	b, err := json.Marshal(u)
	fmt.Println(string(b[:]))
	if err != nil {
		logrus.Fatal(err.Error())
	}
	resp, err := http.Post(api, ct, bytes.NewBuffer(b))
	defer resp.Body.Close()
	logrus.Info("resp update: ", resp.StatusCode)
	if err != nil {
		logrus.Fatal(err)
	}
}

func checkHealth() bool {
	api := fmt.Sprintf("%s/healthz", ApiUrl)
	resp, err := http.Get(api)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func boolPointer(b bool) *bool {
	return &b
}
