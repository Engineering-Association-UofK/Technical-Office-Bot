package telegram

import (
	"log"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/locale"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot     *tgbotapi.BotAPI
	adminID int64
	locale  locale.TelegramLocale
}

func TelegramInit(token string, adminID int64) *TelegramBot {
	log.Println("Initializing Telegram bot...")

	// Initialize the bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Error initializing telegram bot: ", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	t := TelegramBot{
		bot:     bot,
		adminID: adminID,
	}

	// Load Default locale
	t.locale.Load()

	// Begin listening to incoming messages
	go t.Listen()

	return &t
}

func (t *TelegramBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, t.locale.WelcomeMessage)
				msg.ParseMode = "Markdown"
				t.bot.Send(msg)

			default:
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				t.bot.Send(msg)
			}
		}
	}
}

func (t *TelegramBot) NotifyAdmin(message string) {
	msg := tgbotapi.NewMessage(t.adminID, message)
	t.bot.Send(msg)
}
