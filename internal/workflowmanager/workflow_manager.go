package workflowmanager

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/go-github/v67/github"
	"sync"
)

var statusEmojiMap = map[string]string{
	"completed":   "✅",
	"success":     "✅",
	"cancelled":   "⛔️",
	"failure":     "🧨",
	"queued":      "🚚",
	"in_progress": "⏳",
	"requested":   "⏳",
	"waiting":     "⏳",
	"pending":     "⏳",
}

type WorkflowManager struct {
	progressMap map[int64]*models.Message
	mu          sync.Mutex
}

func NewWorkflowManager() *WorkflowManager {
	return &WorkflowManager{
		progressMap: make(map[int64]*models.Message),
	}
}

func (wm *WorkflowManager) HandleProgress(ctx context.Context, b *bot.Bot, workflowRun github.WorkflowRun, chatID any) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	status := getStatus(workflowRun)
	msgText := fmt.Sprintf("%v \\[%s\\] workflow run `%s`", statusEmojiFor(workflowRun), bot.EscapeMarkdown(status), *workflowRun.Name)
	msgText += bot.EscapeMarkdown("\n\n")
	msgText += fmt.Sprintf("[%s](%s)", bot.EscapeMarkdown(*workflowRun.Repository.FullName), bot.EscapeMarkdown(*workflowRun.Repository.HTMLURL))

	cancelButton := models.InlineKeyboardButton{
		Text:         "cancel",
		CallbackData: fmt.Sprintf("cancel_%v_%s", *workflowRun.ID, *workflowRun.Repository.FullName),
	}

	m, exists := wm.progressMap[*workflowRun.ID]
	if !exists {
		sendParams := &bot.SendMessageParams{
			ChatID:    chatID,
			ParseMode: models.ParseModeMarkdown,
			Text:      msgText,
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						cancelButton,
					},
				},
			},
		}

		m, err := b.SendMessage(ctx, sendParams)
		if err != nil {
			return err
		}

		wm.progressMap[*workflowRun.ID] = m
		return nil
	}

	editKeyboard := []models.InlineKeyboardButton{
		{Text: "View", URL: *workflowRun.HTMLURL},
	}

	if status == "queued" || status == "in_progress" {
		editKeyboard = append(editKeyboard, cancelButton)
	}

	editParams := &bot.EditMessageTextParams{
		ChatID:    chatID,
		ParseMode: models.ParseModeMarkdown,
		MessageID: m.ID,
		Text:      msgText,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				editKeyboard,
			},
		},
	}

	_, err := b.EditMessageText(ctx, editParams)
	if err != nil {
		return err
	}

	return nil
}

func statusEmojiFor(workflowRun github.WorkflowRun) string {
	status := getStatus(workflowRun)

	if emoji, exists := statusEmojiMap[status]; exists {
		return emoji
	}
	return ""
}

func getStatus(workflowRun github.WorkflowRun) string {
	if workflowRun.Conclusion != nil && *workflowRun.Conclusion != "" {
		return *workflowRun.Conclusion
	}

	return *workflowRun.Status
}
