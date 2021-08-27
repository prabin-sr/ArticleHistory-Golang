package archivespider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	configurations "../../../configurations"
	schema "../../../schema"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// DataFetcher - Repeats data fetching.
func DataFetcher(structArchive schema.Archive, archiveURL string, apiKey string, year string, month string) {
	// Check month is current month and year is current year.
	// If so quit crawling. We dont need current month's data.
	yearNow, monthNowName, _ := time.Now().Date()
	monthNow := int(monthNowName)

	yearNowStr := strconv.Itoa(yearNow)
	monthNowStr := strconv.Itoa(monthNow)

	if year == yearNowStr && month == monthNowStr {
		fmt.Println("You have completed crawling until last month data. Crawler will check in next 24 hours.")
		// Sleep 24 Hours and Run again to check if day is on next month or not
		duration := time.Duration(24) * time.Hour
		time.Sleep(duration)

		// Closure function call.
		DataFetcher(structArchive, archiveURL, apiKey, year, month)
	}

	url := strings.ReplaceAll(archiveURL, "{year}", year)
	url = strings.ReplaceAll(url, "{month}", month)
	url = strings.ReplaceAll(url, "{yourkey}", apiKey)
	fmt.Println(url)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)

		duration := time.Duration(20) * time.Second
		time.Sleep(duration)

		// Closure function call.
		DataFetcher(structArchive, archiveURL, apiKey, year, month)
	}

	data, _ := ioutil.ReadAll(response.Body)

	if strings.Contains(string(data), `"detail":{"errorcode":"policies.ratelimit.QuotaViolation"}`) {

		fmt.Println("Repeating...")

		duration := time.Duration(5) * time.Second
		time.Sleep(duration)

		// Closure function call.
		DataFetcher(structArchive, archiveURL, apiKey, year, month)

	}

	pathConfig := configurations.GetPathConfigurations()

	archivePath := path.Join(pathConfig.RawCorpusRoot, pathConfig.ArchivePath, year)

	os.MkdirAll(archivePath, os.ModePerm)

	archiveFile := path.Join(archivePath, month+".json")

	fErr := ioutil.WriteFile(archiveFile, data, 0644)
	check(fErr)

	y, YErr := strconv.Atoi(year)
	check(YErr)

	m, MErr := strconv.Atoi(month)
	check(MErr)

	currentTime := time.Now()

	structArchive.ID = structArchive.ID + 1
	structArchive.Year = y
	structArchive.Month = m
	structArchive.Status = true
	if strings.Contains(string(data), `<?xml version="1.0" encoding="UTF-8"?>`) {
		structArchive.Status = false
	}
	structArchive.CreatedAt = &currentTime
	structArchive.UpdatedAt = &currentTime

	schema.DBConn.Create(&structArchive)

	if m < 12 {
		month = strconv.Itoa(m + 1)
	} else {
		year = strconv.Itoa(y + 1)
		month = "1"
	}

	fmt.Println("Done.")
	// Closure function call.
	DataFetcher(structArchive, archiveURL, apiKey, year, month)

}

// NYMetaDataGrabber - Extract data from NY API.
func NYMetaDataGrabber() {
	NYObj := configurations.GetNYTimesAPIConfig()
	credentials := NYObj.Credentials
	apiEndpoints := NYObj.APIEndpoints

	apiKey := credentials["api_key"]
	archiveURL := apiEndpoints["archive"]

	structArchive := schema.Archive{}

	schema.DBConn.Last(&structArchive)

	var year string
	var month string

	if structArchive.ID == 0 {

		year = "1851"
		month = "1"

	} else {

		if structArchive.Month < 12 {

			year = strconv.Itoa(structArchive.Year)
			month = strconv.Itoa(structArchive.Month + 1)

		} else {

			year = strconv.Itoa(structArchive.Year + 1)
			month = "1"

		}

	}

	DataFetcher(structArchive, archiveURL, apiKey, year, month)
}
