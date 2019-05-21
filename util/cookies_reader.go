package util

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"sync"
)

var lock sync.RWMutex
var CookiesFile string

func CookiesSave(cookiesMap map[string][]*http.Cookie) error {
	lock.Lock()
	defer lock.Unlock()

	file, err := os.Create(CookiesFile)
	//fmt.Println("Save to file: ", cookiesMap)

	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(cookiesMap)
	}
	file.Close()

	return err
}

func CookiesLoad() (map[string][]*http.Cookie, error) {
	var cookiesMap map[string][]*http.Cookie
	var err error

	if _, err = os.Stat(CookiesFile); err == nil {
		// load file
		lock.RLock()
		defer lock.RUnlock()

		file, err := os.Open(CookiesFile)
		if err == nil {
			decoder := gob.NewDecoder(file)
			err = decoder.Decode(&cookiesMap)
		}
		//fmt.Println("Load from file: ", cookiesMap)

		file.Close()
		if cookiesMap == nil {
			//fmt.Println("map is nil returning empty map ")
			cookiesMap = make(map[string][]*http.Cookie)
		}

	} else {
		//fmt.Println("no file ", CookiesFile)

		// serialize an empty map of cookies to file
		cookiesMap = make(map[string][]*http.Cookie)
		err = CookiesSave(cookiesMap)
	}

	return cookiesMap, err
}

func CookiesAdd(response *http.Response) (string, []*http.Cookie, error) {
	cookiesMap, err := CookiesLoad()

	if err != nil {
		fmt.Println("failed to load cookies from file")
		return "", nil, err
	}

	//fmt.Println("Cookies cached:", cookiesMap)
	newCookies := response.Cookies()
	url := fmt.Sprintf("%s://%s", response.Request.URL.Scheme, response.Request.URL.Host)

	//fmt.Printf("adding: [%s] %v", url, newCookies)

	if len(newCookies) > 0 {
		if len(cookiesMap[url]) == 0 {
			cookiesMap[url] = newCookies
		} else {
			cookiesMap[url] = append(cookiesMap[url], newCookies...)
		}
	}

	err = CookiesSave(cookiesMap)
	if err != nil {
		fmt.Println("failed to save cookies")
		return "", nil, err
	}

	//fmt.Println("Cookies newly cached:", cookiesMap)
	return url, newCookies, nil
}

func CookiesClear() error {
	err := CookiesSave(nil)

	if err != nil {
		fmt.Println("failed to clear cookies cache:", err)
	}
	return err
}
