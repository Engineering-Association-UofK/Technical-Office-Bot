package telegram

import (
	"log"
	"strings"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/locale"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/models"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/repository"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type TelegramBot struct {
	bot       *tgbotapi.BotAPI
	adminID   int64
	lm        *locale.LocaleManager
	notify    <-chan string
	repo      *repository.TelegramRepo
	fbService *service.FeedbackService
}

func TelegramInit(token string, adminID int64, db *sqlx.DB, fbService *service.FeedbackService, notificationChannel <-chan string) *TelegramBot {
	log.Println("Initializing Telegram bot...")

	// Initialize the bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Error initializing telegram bot: ", err)
	}

	// bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	t := TelegramBot{
		bot:     bot,
		adminID: adminID,
		lm:      locale.NewLocaleManager(),
		notify:  notificationChannel,
		repo: &repository.TelegramRepo{
			Tuser: repository.BaseRepo[models.TelegramUser]{
				DB:        db,
				TableName: "telegram_users",
			},
			Tinteraction: repository.BaseRepo[models.TelegramInteraction]{
				DB:        db,
				TableName: "telegram_interactions",
			},
		},
		fbService: fbService,
	}

	// Begin listening to incoming messages
	go t.Listen()
	go t.NotifyAdmin()

	return &t
}

func (t *TelegramBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Save or update user on the database
		id, err := t.repo.InteractionSave(update.Message)
		if err != nil {
			log.Printf("DB Error: %v", err)
		}
		user, _ := t.repo.FindById(id)

		// Get user locale
		userLang := update.Message.From.LanguageCode
		locale := t.lm.Get(userLang)

		text := update.Message.Text

		// Check for feedback command
		if strings.HasPrefix(text, "/feedback") {

			// Remove the "/feedback" part and trim whitespace
			feedbackContent := strings.TrimSpace(strings.TrimPrefix(text, "/feedback"))

			if feedbackContent == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, locale.FeedbackEmpty)
				t.bot.Send(msg)
				return
			}

			_, err = t.fbService.TelegramFeedback(&user, feedbackContent)
			if err != nil {
				log.Println("Error Saving Telegram feedback: ", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, locale.FeedbackThanks)
			t.bot.Send(msg)
			return
		}

		switch text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, locale.WelcomeMessage)
			msg.ParseMode = "Markdown"
			t.bot.Send(msg)

		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, locale.HelpMessage)
			t.bot.Send(msg)

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			t.bot.Send(msg)
		}
	}
}

func (t *TelegramBot) NotifyAdmin() {
	for message := range t.notify {
		msg := tgbotapi.NewMessage(t.adminID, message)
		t.bot.Send(msg)
	}
}
