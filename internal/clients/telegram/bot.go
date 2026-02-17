package telegram

import (
	"log/slog"
	"strings"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/locale"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/repository"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type TelegramBot struct {
	bot       *tgbotapi.BotAPI
	lm        *locale.LocaleManager
	notify    <-chan string
	repo      *repository.TelegramRepo
	fbService *service.FeedbackService
}

func TelegramInit(token string, db *sqlx.DB, fbService *service.FeedbackService, notificationChannel <-chan string) (*TelegramBot, error) {
	slog.Info("Initializing Telegram bot...")

	// Initialize the bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	// bot.Debug = true
	slog.Info("Authorized on account " + bot.Self.UserName)

	t := TelegramBot{
		bot:    bot,
		lm:     locale.NewLocaleManager(),
		notify: notificationChannel,
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
	go t.NotifyTechnicalAdmins()

	return &t, nil
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
		_, err := t.repo.InteractionSave(update.Message)
		if err != nil {
			slog.Error("DB Error: " + err.Error())
		}

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

			_, err = t.fbService.NotifyFeedback(update.Message.From.UserName, update.Message.From.ID, feedbackContent)
			if err != nil {
				slog.Error("Error Saving Telegram feedback: " + err.Error())
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

func (t *TelegramBot) NotifyTechnicalAdmins() {
	for message := range t.notify {
		admins, err := t.repo.FindTechnicalAdmins()
		if err != nil {
			slog.Error("Error getting admins to notify: " + err.Error())
		}
		for _, admin := range admins {
			msg := tgbotapi.NewMessage(admin.TelegramID, message)
			t.bot.Send(msg)
		}
	}
}
