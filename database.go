package database

import (
	"strconv"
	"strings"
)

var FilterArr []GroupPersonalPageInformation
var FilterTmpl GroupPersonalPageInformation

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
