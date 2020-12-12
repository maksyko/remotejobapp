package broadcasts

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ievgen-ma/remotejobapp/models"
	"log"
	"os"
	"strconv"
	"strings"
)

type broadcastHandler struct {
	posts  *models.PostsHandler
	clint  *tgbotapi.BotAPI
	chatID int64
}

func NewBroadcastHandler() *broadcastHandler {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	chatID, err := strconv.ParseInt(os.Getenv("TG_CHANNEL_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return &broadcastHandler{
		clint:  bot,
		chatID: chatID,
		posts:  models.NewPostsHandler(),
	}
}

func (h *broadcastHandler) getTemplate(title, salary, apply, company, position string) string {
	return fmt.Sprintf(
		"*%s* %s\n"+
			"👉 [%s](%s) %s - %s\n",
		title,
		salary,
		"APPLY",
		apply,
		company,
		position,
	)
}

func (h *broadcastHandler) send(text string) string {
	msg := tgbotapi.NewMessage(h.chatID, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true

	m, err := h.clint.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	return m.Text
}

func (h *broadcastHandler) SendPosts(limit int) error {
	ps, err := h.posts.FindPosts(limit)
	if err != nil {
		log.Fatal(err)
	}

	var postsText []string
	for _, p := range ps {
		postsText = append(postsText, h.getTemplate(p.Title, p.Salary, p.Apply, p.Company, p.Position))
	}
	text := strings.Join(postsText, "")

	h.send(text)

	return h.posts.Processed(ps)

}
