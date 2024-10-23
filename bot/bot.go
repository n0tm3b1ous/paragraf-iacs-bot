package bot

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/n0tm3b1ous/paragraf-iacs-bot/api"
	"github.com/n0tm3b1ous/paragraf-iacs-bot/utils"

	try "github.com/dsnet/try"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	botApi              *tgbot.BotAPI
	updatesChannel      tgbot.UpdatesChannel
	currentUpdate       tgbot.Update
	paragrafApiInstance api.ParagrafApi
}

func (bot *TelegramBot) Init(token string, timeout int, apiInstance api.ParagrafApi) error {
	botApi, err := tgbot.NewBotAPI(token)
	if err != nil {
		return err
	}
	u := tgbot.NewUpdate(0)
	u.Timeout = timeout
	updates := botApi.GetUpdatesChan(u)
	bot.botApi = botApi
	bot.updatesChannel = updates
	bot.paragrafApiInstance = apiInstance
	return nil
}

func (bot TelegramBot) sendTextMsg(text string) {
	msg := tgbot.NewMessage(bot.currentUpdate.Message.Chat.ID, text)
	bot.botApi.Send(msg)
}

func (bot TelegramBot) sendErrMsg(apiErr error, log bool) {
	if log {
		utils.ErrorHandler(apiErr, bot.paragrafApiInstance.LogPath)
	}
	msg := tgbot.NewMessage(bot.currentUpdate.Message.Chat.ID, apiErr.Error())
	msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
	bot.botApi.Send(msg)
}

func (bot TelegramBot) sendKeyboard(buttonsPerRow int, data []string, msgText string) {
	var keyboard [][]tgbot.KeyboardButton
	for i := 0; i < len(data); i++ {
		if i%buttonsPerRow == 0 {
			keyboard = append(keyboard, tgbot.NewKeyboardButtonRow())
		}
		keyboard[i/buttonsPerRow] = append(keyboard[i/buttonsPerRow], tgbot.NewKeyboardButton(data[i]))
	}
	subjKeyboard := tgbot.NewReplyKeyboard(keyboard...)
	msg := tgbot.NewMessage(bot.currentUpdate.Message.Chat.ID, msgText)
	msg.ReplyMarkup = subjKeyboard
	bot.botApi.Send(msg)
}

func (bot TelegramBot) welcome(version string) {
	welcomeMsg := fmt.Sprintf("Paragraf AICS bot v%s\n\nCOMMANDS:\n/help - Вывести это сообщение\n/marks - Управление оценками\n/stats - Статистика бота\n/random - ???", version)
	bot.sendTextMsg(welcomeMsg)
}

func (bot TelegramBot) requestGrade() (api.Grade, error) {
	bot.sendKeyboard(4, []string{"8", "9", "10", "11"}, "Введите номер парралели")
	for nextUpdate := range bot.updatesChannel {
		if num, err := strconv.Atoi(nextUpdate.Message.Text); err == nil && num >= 8 && num <= 11 {
			grades, err := bot.paragrafApiInstance.GetMenu()
			if err != nil {
				return api.Grade{}, err
			}
			grade := grades[num-1]
			return grade, nil
		} else {
			bot.sendTextMsg("Ввод не распознан.")
		}
	}
	return api.Grade{}, errors.New("Не удалось запросить данные о параллели ")
}

func (bot TelegramBot) requestClass(grade api.Grade) (api.Class, error) {
	var subject api.Subject
	subjectNames := grade.GetSubjectsNames()
	bot.sendKeyboard(4, subjectNames, "Введите предмет")
	for nextUpdate := range bot.updatesChannel {
		if nextUpdate.Message != nil && slices.Contains(subjectNames, nextUpdate.Message.Text) {
			subject = grade.Subjects[slices.Index(subjectNames, nextUpdate.Message.Text)]
			break
		} else {
			bot.sendTextMsg("Ввод не распознан.")
		}
	}
	classesNames := subject.GetClassesNames()
	bot.sendKeyboard(4, classesNames, "Введите класс")
	for nextUpdate := range bot.updatesChannel {
		if nextUpdate.Message != nil && slices.Contains(classesNames, nextUpdate.Message.Text) {
			return subject.Classes[slices.Index(classesNames, nextUpdate.Message.Text)], nil
		} else {
			bot.sendTextMsg("Ввод не распознан.")
		}
	}
	return api.Class{}, errors.New("Не удалось запросить данные о классе")
}

func (bot TelegramBot) requestMarksInfo(class api.Class) ([]api.MarkDetails, error) {
	var marksInfo []api.MarkDetails
	journal, err := bot.paragrafApiInstance.GetJournal(class)
	if err != nil {
		return []api.MarkDetails{}, err
	}
	studentsNames := journal.GetStudentsNames()
	bot.sendKeyboard(5, studentsNames, "Введите имя студента")
	for nextUpdate := range bot.updatesChannel {
		if nextUpdate.Message != nil && slices.Contains(studentsNames, nextUpdate.Message.Text) {
			for _, mark := range journal.Marks {
				student := journal.Members[slices.Index(studentsNames, nextUpdate.Message.Text)]
				if mark.StudentId == student.Id {
					markInfo, err := bot.paragrafApiInstance.GetMarkDetails(mark)
					if err != nil {
						return []api.MarkDetails{}, err
					}
					marksInfo = append(marksInfo, markInfo)
				}
			}
			return marksInfo, nil
		} else {
			bot.sendTextMsg("Ввод не распознан.")
		}
	}
	return []api.MarkDetails{}, errors.New("Не удалось полчучить данные об оценках")
}

func (bot TelegramBot) markController() (err error) {
	defer try.HandleF(&err, func() { bot.sendErrMsg(err, true) })
	grade := try.E1(bot.requestGrade())
	class := try.E1(bot.requestClass(grade))
	marksInfo := try.E1(bot.requestMarksInfo(class))
	marksTable := ""
	for _, markInfo := range marksInfo {
		marksTable += fmt.Sprintf("%s - %s\n", markInfo.DateAdd, markInfo.Value)
	}
	if marksTable != "" {
		msg := tgbot.NewMessage(bot.currentUpdate.Message.Chat.ID, marksTable)
		msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		bot.botApi.Send(msg)
	} else {
		bot.sendErrMsg(errors.New("Список оценок пуст"), false)
	}
	return nil
}

func (bot *TelegramBot) StartTgBot() {
	ticker := time.NewTicker(30 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				bot.paragrafApiInstance.UpdateSession()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	for {
		for update := range bot.updatesChannel {
			bot.currentUpdate = update
			if update.Message == nil {
				continue
			}
			switch update.Message.Command() {
			case "help":
				bot.welcome(bot.paragrafApiInstance.Version)
			case "marks":
				bot.markController()
			case "stats":
				bot.sendTextMsg("В разработке")
			case "random":
				bot.sendTextMsg("В разработке")
			default:
				bot.sendTextMsg("Команда не распознана. Чтобы посмотреть список комманд, введите /help")
			}
		}
	}
}
