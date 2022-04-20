package integration

import (
	"fmt"

	"github.com/jelliflix/jelliflix/infrastructure/config"
	"github.com/jelliflix/jelliflix/infrastructure/logger"
	"gopkg.in/tucnak/telebot.v2"
)

type Integration struct {
	user *telebot.User
	bot  *telebot.Bot
}

func NewIntegration() *Integration {
	user := &telebot.User{ID: config.Cfg.Integrations.Telegram.User}
	bot, err := telebot.NewBot(telebot.Settings{
		Token: config.Cfg.Integrations.Telegram.Token,
	})
	if err != nil {
		logger.Log.Error(err)
	}

	return &Integration{bot: bot, user: user}
}

func (i *Integration) Queued(name, size, quality string) {
	x := `
**%s** added to queue.
Size: %s ▫️ Quality: %s
`
	_, err := i.bot.Send(i.user, fmt.Sprintf(x, name, size, quality), telebot.NoPreview)
	if err != nil {
		logger.Log.Warn(err)
	}
}

func (i *Integration) Finished(name string) {
	x := `
**%s** downloaded.
Enjoy watching 🙌!
`
	_, err := i.bot.Send(i.user, fmt.Sprintf(x, name), telebot.NoPreview)
	if err != nil {
		logger.Log.Warn(err)
	}
}
