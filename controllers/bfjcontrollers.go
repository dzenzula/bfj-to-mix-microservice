package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/config"
	"main/models"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{}

func GetListBf() (nBF []int) {
	var data models.ListBF
	url := config.GlobalConfig.BFJAPI.ApiGetListBF
	req, getListOfBFErr := http.Get(url)
	if getListOfBFErr != nil {
		fmt.Println(getListOfBFErr.Error())
		return nil
	}

	body, readingErr := io.ReadAll(req.Body)
	if readingErr != nil {
		fmt.Println(readingErr.Error())
		return nil
	}

	jsonError := json.Unmarshal(body, &data)
	if jsonError != nil {
		fmt.Println(jsonError.Error())
		return nil
	}

	return data.Name
}

func GetLastBFJJournalsData(nBF []int, ids *map[int][]int) {
	var data models.Journals

	if *ids == nil {
		*ids = make(map[int][]int)
	}

	for _, n := range nBF {
		endpoint := fmt.Sprintf(config.GlobalConfig.BFJAPI.ApiGetLastJournals, strconv.Itoa(n))
		err := getBfjApiResponse(endpoint, nil, &data)
		if err != nil {
			return
		}

		for i := 0; i < 1; i++ {
			(*ids)[n] = append((*ids)[n], data.DataJournals[i].ID)
		}
	}
}

func GetBFJTappings(journalId int, cookies *[]*http.Cookie) (tappingIds []models.Tapping) {
	var data []models.Tapping
	endpoint := fmt.Sprintf(config.GlobalConfig.BFJAPI.ApiGetTappings, strconv.Itoa(journalId))
	err := getBfjApiResponse(endpoint, cookies, &data)
	if err != nil {
		return nil
	}
	return data
}

func AuthorizeBFJ(cookies *[]*http.Cookie) {
	auth, _ := json.Marshal(config.GlobalConfig.Auth)

	for {
		success := true
		req, err := http.Post(config.GlobalConfig.BFJAPI.ApiPostAuthProd, "application/json", bytes.NewBuffer(auth))
		if err != nil {
			success = false
			fmt.Printf("Failed to send authorization request: %v", err)
		}
		defer req.Body.Close()

		if req.StatusCode != http.StatusOK {
			success = false
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				fmt.Printf("Failed to read authorization response body: %v\n", err)
			}
			fmt.Printf("Rejected authorization request: %s\n", bodyBytes)
			fmt.Println("Next try to authorize will be in a 5 minutes")
			time.Sleep(time.Minute * 5)
		}

		if success {
			*cookies = req.Cookies()
			return
		} else {
			fmt.Println("Next try to authorize will be in a 5 minutes")
			time.Sleep(time.Minute * 5)
		}
	}
}

func getBfjApiResponse(endpoint string, cookies *[]*http.Cookie, data interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	if cookies != nil {
		for _, cookie := range *cookies {
			req.AddCookie(cookie)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	} else if resp.StatusCode != http.StatusOK {
		AuthorizeBFJ(cookies)
		getBfjApiResponse(endpoint, cookies, data)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error decoding JSON string:", err)
		return err
	}

	return nil
}
