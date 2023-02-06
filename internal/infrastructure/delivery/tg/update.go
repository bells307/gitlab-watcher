package tg

import (
	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandler struct {
	updateService UpdateService
}

// Сервис обработки Update'ов телеграма
type UpdateService interface {
	AddedToChat(*tm.Update)
	RemovedFromChat(*tm.Update)
}

func NewUpdateHandler(updateService UpdateService) *UpdateHandler {
	return &UpdateHandler{updateService}
}

func (h *UpdateHandler) Register(api *tgbotapi.BotAPI, mux *tm.Mux) {
	mux.
		AddHandler(tm.NewHandler(
			func(u *tm.Update) bool {
				if mcm := u.MyChatMember; mcm != nil {
					ncm := mcm.NewChatMember
					if ncm.User != nil {
						if ncm.User.ID == api.Self.ID {
							if (ncm.Status != "left") && (ncm.Status != "kicked") {
								return true
							}
						}
					}
				}

				return false
			},
			func(u *tm.Update) {
				h.updateService.AddedToChat(u)
			},
		)).
		AddHandler(tm.NewHandler(
			func(u *tm.Update) bool {
				if mcm := u.MyChatMember; mcm != nil {
					ncm := mcm.NewChatMember
					if ncm.User != nil {
						if ncm.User.ID == api.Self.ID {
							if (ncm.Status == "left") || (ncm.Status == "kicked") {
								return true
							}
						}
					}
				}

				return false
			},
			func(u *tm.Update) {
				h.updateService.RemovedFromChat(u)
			},
		))
}
