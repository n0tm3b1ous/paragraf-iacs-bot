package main

import (
	"github.com/n0tm3b1ous/paragraf-iacs-bot/api"
	"github.com/n0tm3b1ous/paragraf-iacs-bot/bot"
	"github.com/n0tm3b1ous/paragraf-iacs-bot/conf"
	"github.com/n0tm3b1ous/paragraf-iacs-bot/utils"
)

func main() {
	apiInstance := api.ParagrafApi{conf.DefaultConfig.Version, conf.DefaultConfig.ApiLogin, conf.DefaultConfig.ApiPassword, conf.DefaultConfig.ApiBasePath, conf.DefaultConfig.LogPath, ""}
	apiInstance.UpdateSession()
	ParagrafBot := bot.TelegramBot{}
	err := ParagrafBot.Init(conf.DefaultConfig.TgBotToken, 60, apiInstance)
	if err != nil {
		utils.ErrorHandler(err, conf.DefaultConfig.LogPath)
	}
	ParagrafBot.StartTgBot()
}
