package database

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"fmt"
	"strings"
)

const apiURL = "https://groupietrackers.herokuapp.com/api"

type IndexLocation struct {
	Index []Location `json:"index"`
}

type IndexDate struct {
	Index []Date `json:"index"`
}

type IndexRelation struct {
	Index []Relation `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID int `json:"id"`
	Relations map[string][]string `json:"datesLocations"`
}

type GroupPersonalPageInformation struct {
	ID int `json:"id"`
	Image string `json:"image"`
	Name string `json:"name"`
	Members []string `json:"members"`
	CreationDate int `json:"creationDate"`
	FirstAlbum string `json:"firstAlbum"`
	Locations      []string            `json:"locations"`
	ConcertDates   []string            `json:"concertDates"`
	Relations map[string][]string `json:"datesLocations"`
}

type GroupHomePageInformation struct {
	ID int `json:"id"`
	Image string `json:"image"`
	Name string `json:"name"`
	Members []string `json:"members"`
	Locations []string `json:"locations"`
	CreationDate int `json:"creationDate"`
	FirstAlbum string `json:"firstAlbum"`
}



var FilterArr []GroupPersonalPageInformation
var FilterTmpl GroupPersonalPageInformation
var LocationsData IndexLocation
var DatesData IndexDate
var RelationsData IndexRelation
var HomePageInformation []GroupHomePageInformation
var PersonalPageInformation []GroupPersonalPageInformation
var SearchArr []GroupPersonalPageInformation
var SearchTmpl GroupPersonalPageInformation

func Search(str string) {
	GetPersonalPageData()
	SearchArr = nil
	for i, v := range PersonalPageInformation {
		if str == v.Name {
			SearchTmpl = PersonalPageInformation[i]
			SearchArr = append(SearchArr, SearchTmpl)
			return
		}
		MembersCount := 0
		for _, value := range v.Members {
			if str == value {
				MembersCount++
			}
		}
		if MembersCount > 0 {
			SearchTmpl = PersonalPageInformation[i]
			SearchArr = append(SearchArr, SearchTmpl)
			return
		}
		LocationsCount := 0
		for _, value := range v.Locations {
			if str == value {
				LocationsCount++
			}
		}
		if LocationsCount > 0 {
			SearchTmpl = PersonalPageInformation[i]
			SearchArr = append(SearchArr, SearchTmpl)
		}
		CreationDate := strconv.Itoa(v.CreationDate)
		if str == CreationDate {
			SearchTmpl = PersonalPageInformation[i]
			SearchArr = append(SearchArr, SearchTmpl)
		}
		if str == v.FirstAlbum {
			SearchTmpl = PersonalPageInformation[i]
			SearchArr = append(SearchArr, SearchTmpl)
		}
	}
	fmt.Println(SearchArr)
}

func GetLocationData() error {
	res, err := http.Get(apiURL + "/locations")
	if err != nil {
		return errors.New("Error by get /locations")
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("Error by ReadAll /locations")
	}
	json.Unmarshal(bytes, &LocationsData)
	return nil
}

func GetDatesData() error {
	res, err := http.Get(apiURL + "/dates")
	if err != nil {
		return errors.New("Error by get /dates")
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("Error by ReadAll /dates")
	}
	json.Unmarshal(bytes, &DatesData)
	return nil
}

func GetRelationsData() error {
	res, err := http.Get(apiURL + "/relation")
	if err != nil {
		return errors.New("Error by get /relation")
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("Error by ReadAll /relation")
	}
	json.Unmarshal(bytes, &RelationsData)
	return nil
}

func GetHomePageData() ([]GroupHomePageInformation, error) {
	res, err := http.Get(apiURL + "/artists")
	if err != nil {
		return nil, errors.New("GetData error")
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("ReadAll error")
	}
	json.Unmarshal(bytes, &HomePageInformation)
	return HomePageInformation, nil
}


func GetPersonalPageData() error {

	if len(PersonalPageInformation) != 0 {
		return nil
	}

	_, errHomePageData := GetHomePageData()
	if errHomePageData != nil {
		return errHomePageData
	}

	errLocationsData := GetLocationData()
	if errLocationsData != nil {
		return errHomePageData
	}

	errDatesData := GetDatesData()
	if errDatesData != nil {
		return errHomePageData
	}

	errRelationsData := GetRelationsData()
	if errRelationsData != nil {
		return errHomePageData
	}

	for i := range HomePageInformation {
		var tmpl GroupPersonalPageInformation
		tmpl.ID = i + 1
		tmpl.Image = HomePageInformation[i].Image
		tmpl.Name = HomePageInformation[i].Name
		tmpl.Members = HomePageInformation[i].Members
		tmpl.CreationDate = HomePageInformation[i].CreationDate
		tmpl.FirstAlbum = HomePageInformation[i].FirstAlbum
		tmpl.Locations = LocationsData.Index[i].Locations
		tmpl.ConcertDates = DatesData.Index[i].Dates
		tmpl.Relations = RelationsData.Index[i].Relations
		PersonalPageInformation = append(PersonalPageInformation, tmpl)
	}
	return nil
}

func GetFilterInformation(startCD, endCD, startFA, endFA, location string, membersCount []int) {
	GetPersonalPageData()
	count := 0
	FilterArr = nil
	for i, v := range PersonalPageInformation {
		count = 0
		fromCD, _ := strconv.Atoi(startCD)
		tillCD, _ := strconv.Atoi(endCD)
		if v.CreationDate >= fromCD && v.CreationDate <= tillCD {
			count++
		}
		if GetDaysForCompareDate(v.FirstAlbum, 0) >= GetDaysForCompareDate(startFA, 1) && GetDaysForCompareDate(v.FirstAlbum, 0) <= GetDaysForCompareDate(endFA, 1) {
			count++
		}
		if location != "" {
			for _, k := range v.Locations {
				if location == k {
					count++
					break
				}
			}
		}
		if location == "" {
			count++
		}

		for _, j := range membersCount {
			if len(v.Members) == j {
				count++
				break
			}
		}
		if len(membersCount) == 0 {
			count++
		}
		if count == 4 {
			FilterArr = append(FilterArr, PersonalPageInformation[i])
		}
	}
}

func GetDaysForCompareDate(date string, flag int) int {
	dateSplit := strings.Split(date, "-")
	daysCount := 0
	days, _ := strconv.Atoi(dateSplit[0])
	months, _ := strconv.Atoi(dateSplit[1])
	years, _ := strconv.Atoi(dateSplit[2])
	if flag == 1 {
		days, _ := strconv.Atoi(dateSplit[2])
		months, _ := strconv.Atoi(dateSplit[1])
		years, _ := strconv.Atoi(dateSplit[0])
		daysCount += days
		daysCount += months * 30
		daysCount += (years - 1900) * 365
		return daysCount
	}
	daysCount += days
	daysCount += months * 30
	daysCount += (years - 1900) * 365
	return daysCount
}
