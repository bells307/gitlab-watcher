package main

import (
	"flag"
	"fmt"
	"github.com/bells307/gitlab-watcher/config"
	"github.com/bells307/gitlab-watcher/internal/adapters/template"
	"log"

	"github.com/bells307/gitlab-watcher/internal/adapters/sender"

	tm "github.com/and3rson/telemux/v2"
	"github.com/bells307/gitlab-watcher/internal/application"
	"github.com/bells307/gitlab-watcher/internal/infrastructure/delivery/http"
	"github.com/bells307/gitlab-watcher/internal/infrastructure/delivery/tg"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfgPath := flag.String("config", "./config.yml", "config file path")
	flag.Parse()

	log.Printf("loading configuration ...")
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal("can't load configuration:", err)
	}

	var hookServices []http.HookService
	if cfg.Telegram != nil {
		hookServices = append(hookServices, createTelegramService(cfg.Telegram, cfg.SendTemplateErrors))
	}

	router := gin.Default()
	hookHandler := http.NewHookHandler(hookServices)
	hookHandler.Register(router)

	log.Printf("start listening on %s ...", cfg.Listen)
	if err = router.Run(cfg.Listen); err != nil {
		log.Fatalf("error starting http listener: %v", err)
	}
}

// TODO: вынести из main.go
func createTelegramService(cfg *config.TelegramConfig, sendTemplateErrors bool) http.HookService {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal("can't create telegram bot api:", err)
	}

	if cfg.Debug {
		api.Debug = true
	}

	telegramBot := sender.NewTelegramBot(api, cfg.Chats, cfg.ParseMode)
	updateHandler := tg.NewUpdateHandler(telegramBot)
	mux := tm.NewMux()
	updateHandler.Register(api, mux)

	telegramHookTemplater := template.NewHookTextTemplater(cfg.Templates)
	telegramHookTemplater.AddFunc("mapUser", func(gitlabUserName string) (string, error) {
		if user, ok := cfg.Users[gitlabUserName]; ok {
			return user, nil
		} else {
			return "", fmt.Errorf("can't map user %s", gitlabUserName)
		}
	})

	srv := application.NewHookService(telegramBot, telegramHookTemplater, sendTemplateErrors)

	go runTelegramBot(api, mux)

	return srv
}

func runTelegramBot(api *tgbotapi.BotAPI, mux *tm.Mux) {
	log.Println("starting telegram bot ...")

	updConfig := tgbotapi.NewUpdate(0)
	updConfig.Timeout = 60
	updChan := api.GetUpdatesChan(updConfig)

	for upd := range updChan {
		mux.Dispatch(api, upd)
	}
}
