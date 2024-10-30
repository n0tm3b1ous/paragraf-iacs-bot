package utils

import (
	"log"
	"os"
	"net/http"
)

func ErrorHandler(myErr error, logPath string) {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if (err != nil) {
		panic(err)
	}
	log.SetOutput(file)
	log.Println(myErr)
}

func ConnectivityCheck(path string) (int, error) {
	res, err := http.Get(path)
	if err != nil {
		return -1, err
	}
	return res.StatusCode, nil
}

func HttpHandler(path string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", path, nil)
	for header, value := range headers {
		req.Header.Add(header, value)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}