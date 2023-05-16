package main

import (
	"fmt"
	"main/cache"
	"main/config"
	"main/controllers"
	"main/models"
	"net/http"
	"reflect"
	"time"
)

var currList map[int][]models.Ladle = make(map[int][]models.Ladle)
var bfjCookies []*http.Cookie
var mixCookies []*http.Cookie
var bfjIds map[int][]int = map[int][]int{}
var mixIds map[int][]int = map[int][]int{}
var data *cache.Data = cache.ReadYAMLFile(config.GlobalConfig.Path.CachePath)

func main() {
	controllers.AuthorizeBFJ(&bfjCookies)
	controllers.AuthorizeMix(&mixCookies)

	var nBF []int = controllers.GetListBf()
	var nMix []int = []int{1, 2, 3, 4}

	controllers.GetLastBFJJournalsData(nBF, &bfjIds)
	controllers.GetLastMIXJournalsData(nMix, &mixCookies, &mixIds)

	cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, bfjIds, data.Tappings)

	duration := time.Until(time.Now().Truncate(time.Minute).Add(time.Minute))
	time.Sleep(duration)

	fmt.Println(time.Now().String())

	for {
		service(nBF, &bfjIds, &bfjCookies, nMix, &mixIds, &mixCookies)

		ticker := time.NewTicker(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
		<-ticker.C
	}
}

func service(nBF []int, bfjIds *map[int][]int, bfjCookies *[]*http.Cookie, nMix []int, mixIds *map[int][]int, mixCookies *[]*http.Cookie) {
	now := time.Now().Truncate(time.Minute)
	if (now.Hour() == 8 && now.Minute() == 0) || (now.Hour() == 20 && now.Minute() == 0) {
		currList = make(map[int][]models.Ladle)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Shift started")

		controllers.GetLastBFJJournalsData(nBF, bfjIds)
		controllers.GetLastMIXJournalsData(nMix, mixCookies, mixIds)

		cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, *bfjIds, nil)

		minuteCheck(bfjIds, bfjCookies, mixIds, mixCookies)
	} else if now.Second() == 0 {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Minute check started")

		minuteCheck(bfjIds, bfjCookies, mixIds, mixCookies)
	}
}

func minuteCheck(bfjIds *map[int][]int, bfjCookies *[]*http.Cookie, mixIds *map[int][]int, mixCookies *[]*http.Cookie) {
	data = cache.ReadYAMLFile(config.GlobalConfig.Path.CachePath)
	for nBf, values := range *bfjIds {
		for _, idJournal := range values {
			tappings := controllers.GetBFJTappings(idJournal, bfjCookies)
			for _, tapping := range tappings {
				sendLadleMovements(nBf, data, tapping, mixIds, mixCookies)
			}
		}
	}
	cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, data.Tappings)
}

func sendLadleMovements(nBf int, tIds *cache.Data, tapping models.Tapping, mixIds *map[int][]int, mixCookies *[]*http.Cookie) {
	if !cache.TappingIdExists(tIds, tapping.ID) {
		ldlMvm := controllers.PostMixLadleMovement(nBf, tapping)
		controllers.PostMixListLadles(tapping.ListLaldes, ldlMvm, mixIds, mixCookies)
		controllers.PostMixChemical(tapping, mixCookies)
		tIds.Tappings = append(tIds.Tappings, map[int]int{tapping.ID: len(tapping.ListLaldes)})
	} else {
		tVal := cache.FindTappingIdValue(tIds, tapping.ID)
		if len(tapping.ListLaldes) != tVal {
			numMissingLadles := len(tapping.ListLaldes) - tVal
			missingLadles := tapping.ListLaldes[len(tapping.ListLaldes)-numMissingLadles:]

			ldlMvm := controllers.PostMixLadleMovement(nBf, tapping)
			controllers.PostMixListLadles(missingLadles, ldlMvm, mixIds, mixCookies)
			controllers.PostMixChemicalList(missingLadles, mixCookies)

			cache.UpdateTappingValue(tIds, tapping.ID, len(tapping.ListLaldes))
			cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, tIds.Tappings)
		}
	}

	// Обработка изменений в составе химикатов
	if currList[tapping.ID] == nil {
		currList[tapping.ID] = tapping.ListLaldes
	} else if !reflect.DeepEqual(currList[tapping.ID], tapping.ListLaldes) {
		oldLadles := currList[tapping.ID]
		currLadles := tapping.ListLaldes
		changedLadles := []models.Ladle{}

		for _, newLadle := range currLadles {
			for _, currLadle := range oldLadles {
				currLadle.Chemical.DtUpdate = newLadle.Chemical.DtUpdate
				if newLadle.Ladle == currLadle.Ladle && !reflect.DeepEqual(newLadle.Chemical, currLadle.Chemical) {
					changedLadles = append(changedLadles, newLadle)
				}
			}
		}

		if len(changedLadles) != 0 {
			fmt.Println(time.Now().Truncate(time.Minute).String(), "Chemical changed!")
			controllers.PostMixChemicalList(changedLadles, mixCookies)
		}

		currList[tapping.ID] = tapping.ListLaldes
	}
}
