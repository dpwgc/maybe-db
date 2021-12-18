package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
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
	HttpClient.Timeout = time.Second * 3
}

func Get(uri string, header map[string]string) (string, error) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if len(header) > 0 {
		for key, value := range header {
			request.Header.Add(key, value)
		}
	}
	response, err := HttpClient.Do(request)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
		log.Println(err)
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(body), nil
}

func PostJson(uri string, header map[string]string, json string) (string, error) {

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(json)))
	if err != nil {
		log.Println(err)
		return "", err
	}

	_, has := header["Content-Type"]
	if !has {
		request.Header.Add("Content-Type", "application/json")
	}
	if len(header) > 0 {
		for key, value := range header {
			request.Header.Add(key, value)
		}
	}

	response, err := HttpClient.Do(request)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(body), nil
}

func HttpGet(url string, header map[string]string) (string, error) {

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if len(header) > 0 {
		for key, value := range header {
			request.Header.Add(key, value)
		}
	}

	response, err := HttpClient.Do(request)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
