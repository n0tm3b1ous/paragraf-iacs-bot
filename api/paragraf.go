package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/n0tm3b1ous/paragraf-iacs-bot/utils"
)

type ParagrafApi struct {
	Version, ApiLogin, ApiPassword, BasePath, LogPath, CurrentSession string
}

func (api *ParagrafApi) UpdateSession() error {
	args := url.Values{}
	args.Set("user-name", api.ApiLogin)
	args.Add("user-password", api.ApiPassword)

	res, err := http.PostForm(api.BasePath+"login", args)
	if err != nil {
		return err
	} else if res.StatusCode == 200 {
		api.CurrentSession = strings.Join(res.Header["Set-Cookie"][:], ",")
		return nil
	}
	return errors.New("Login failed (username/password mismatch)")
}

func (api ParagrafApi) GetJournal(class Class) (Journal, error) {
	var journal Journal
	res, err := utils.HttpHandler(api.BasePath+"webservice/app.cj/execute?action=getdata&id="+class.Id, map[string]string{"Cookie": api.CurrentSession})
	if err != nil {
		return Journal{}, err
	} else if res.StatusCode == 200 {
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		err := json.Unmarshal([]byte(body), &journal)
		if err != nil {
			return Journal{}, err
		}
		return journal, nil
	}
	return Journal{}, errors.New("Failed to get journal")
}

func (api ParagrafApi) GetMenu() ([]Grade, error) {
	var menu []Grade
	res, err := utils.HttpHandler(api.BasePath+"webservice/app.cj/execute?action=menu", map[string]string{"Cookie": api.CurrentSession})
	if err != nil {
		return []Grade{}, err
	} else if res.StatusCode == 200 {
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		err := json.Unmarshal([]byte(body), &menu)
		if err != nil {
			return []Grade{}, err
		}
		return menu, nil
	}
	return []Grade{}, errors.New("Failed to get menu")
}

func (api ParagrafApi) GetMarkDetails(mark Mark) (MarkDetails, error) {
	var markDetails MarkDetails
	res, err := utils.HttpHandler(api.BasePath+"webservice/app.cj/execute?action=mark_details&id="+mark.Id, map[string]string{"Cookie": api.CurrentSession})
	if err != nil {
		return MarkDetails{}, err
	} else if res.StatusCode == 200 {
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		err := json.Unmarshal([]byte(body), &markDetails)
		if err != nil {
			return MarkDetails{}, err
		}
		return markDetails, nil
	}
	return MarkDetails{}, errors.New("Failed to get mark details")
}
