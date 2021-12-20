package utils

import (
	"MaybeDB/server/database"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/*
 * http请求工具
 */

var HttpClient http.Client

func init() {
	HttpClient.Timeout = time.Second * 10
}

func Get(uri string, header map[string]string) (string, error) {

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		database.Loger.Println(err)
		return "", err
	}
	if len(header) > 0 {
		for key, value := range header {
			request.Header.Add(key, value)
		}
	}
	response, err := HttpClient.Do(request)
	if err != nil {
		database.Loger.Println(err)
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		database.Loger.Println(err)
		return "", err
	}
	return string(body), nil
}

func PostForm(uri string, header map[string]string, data map[string]string) (string, error) {

	formData := url.Values{}
	if len(data) > 0 {
		for key, value := range data {
			formData.Set(key, value)
		}
	}

	request, err := http.NewRequest("POST", uri, strings.NewReader(formData.Encode()))
	if err != nil {
		database.Loger.Println(err)
		return "", err
	}

	_, has := header["Content-Type"]
	if !has {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if len(header) > 0 {
		for key, value := range header {
			request.Header.Add(key, value)
		}
	}

	response, err := HttpClient.Do(request)
	if err != nil {
		database.Loger.Println(err)
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		database.Loger.Println(err)
		return "", err
	}

	return string(body), nil
}
