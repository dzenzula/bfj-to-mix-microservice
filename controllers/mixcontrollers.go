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

func GetLastMIXJournalsData(nMIX []int, cookies []*http.Cookie) (ids map[int][]int) {
	var data models.MixJournals
	ids = map[int][]int{}

	for _, n := range nMIX {
		endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiGetLastJournals, strconv.Itoa(n))
		err := getMixApiResponse(endpoint, cookies, &data)
		if err != nil {
			return nil
		}

		if len(data.DataJournals) > 0 {
			for i := 0; i < 1; i++ {
				ids[n] = append(ids[n], data.DataJournals[i].ID)
			}
		}
	}

	return ids
}

func PostMixChemical(tapping models.Tapping, cookies []*http.Cookie) {
	listLadles := tapping.ListLaldes
	postMixChemical(listLadles, cookies)
}

func PostMixChemicalList(listLadles []models.Ladle, cookies []*http.Cookie) {
	postMixChemical(listLadles, cookies)
}

func postMixChemical(listLadles []models.Ladle, cookies []*http.Cookie) {
	endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiPostChemical)
	for nMix := 1; nMix < 5; nMix++ {
		for _, ladle := range listLadles {
			chem := models.ChemicalDTO{
				NMix:       nMix,
				Ladle:      ladle.Ladle,
				Proba:      int(ladle.Chemical.Proba),
				NumTaphole: ladle.Chemical.NumTaphole,
				DT:         checkChemDate(ladle),
				Si:         float64(ladle.Chemical.Si),
				Mn:         float64(ladle.Chemical.Mn),
				S:          float64(ladle.Chemical.S),
				P:          float64(ladle.Chemical.P),
				Belong:     "LadleMovement",
			}
			postMixApiRequest(endpoint, cookies, chem)
		}
	}
}

func checkChemDate(ldl models.Ladle) *string {
	if ldl.Chemical.Dt != "" {
		return &ldl.Chemical.Dt
	}
	return nil
}

func PostMixLadleMovement(nBf int, tapping models.Tapping) models.LadleMovement {
	var ldlMvm models.LadleMovement

	ldlMvm.Date = tapping.DtCloseTaphole
	ldlMvm.NumDp = nBf
	ldlMvm.NumTapping = tapping.NumTapping
	ldlMvm.DtCloseTaphole = tapping.DtCloseTaphole
	ldlMvm.TemperExhaustIron = int(tapping.Temper)

	return ldlMvm
}

func PostMixListLadles(listLadles []models.Ladle, ldlMvm models.LadleMovement, mixIds *map[int][]int, cookies []*http.Cookie) {
	for i := 0; i < len(listLadles); i++ {
		for _, keys := range *mixIds {
			ldlMvm.LadleTapping = listLadles[i].Ladle
			ldlMvm.MassCastIron = int(listLadles[i].Weight)

			for _, key := range keys {
				endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiPostLadleMovement, strconv.Itoa(key))
				postErr := postMixApiRequest(endpoint, cookies, ldlMvm)
				if postErr != nil {
					fmt.Println(postErr.Error())
				}

			}
		}
	}
}

func AuthorizeMix() []*http.Cookie {
	authJSON, _ := json.Marshal(config.GlobalConfig.Auth)

	for {
		success := true
		req, err := http.Post(config.GlobalConfig.MIXAPI.ApiPostAuthTest, "application/json", bytes.NewBuffer(authJSON))
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
			fmt.Println("Authorization MIX success!")
			return req.Cookies()
		}
	}
}

func postMixApiRequest(endpoint string, cookies []*http.Cookie, requestData interface{}) error {
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error encoding request JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}
	req.Header.Set("content-type", "application/json")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, errResp := client.Do(req)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println(string(body))

	if errResp != nil {
		fmt.Println("Error executing HTTP request:", err)
		return err
	} else if resp.StatusCode != http.StatusOK {
		newCookies := AuthorizeMix()

		fmt.Println("Respond MIX error: ", string(body))

		if string(body) == "NotAuthorized" {
			postMixApiRequest(endpoint, newCookies, requestData)
		} else {
			return nil
		}
	}

	return nil
}

func getMixApiResponse(endpoint string, cookies []*http.Cookie, data interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	} else if resp.StatusCode != http.StatusOK {
		newCookies := AuthorizeMix()
		getMixApiResponse(endpoint, newCookies, data)
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
