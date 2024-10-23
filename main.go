package main

import (
	"github.com/n0tm3b1ous/paragraf-iacs-bot/api"
	"github.com/n0tm3b1ous/paragraf-iacs-bot/bot"
	"github.com/n0tm3b1ous/paragraf-iacs-bot/conf"
	"github.com/n0tm3b1ous/paragraf-iacs-bot/utils"
)

func main() {
	apiInstance := api.ParagrafApi{conf.MyConfig.Version, conf.MyConfig.ApiLogin, conf.MyConfig.ApiPassword, conf.MyConfig.ApiBasePath, conf.MyConfig.LogPath, ""}
	apiInstance.UpdateSession()
	ParagrafZla := bot.TelegramBot{}
	err := ParagrafZla.Init(conf.MyConfig.TgBotToken, 60, apiInstance)
	if err != nil {
		utils.ErrorHandler(err, conf.MyConfig.LogPath)
	}
	ParagrafZla.StartTgBot()
}
