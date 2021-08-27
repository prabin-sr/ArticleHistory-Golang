package archivewrangler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	configurations "../../../configurations"
	schema "../../../schema"
)

// CleanerResponse - Struct to Unmarshal json byte data
type CleanerResponse struct {
	Response map[string]interface{} `json:"response"`
}

// CleanedCorpus - struct type for cleaned corpus
type CleanedCorpus struct {
	Abstract      string `json:"abstract"`
	WebURL        string `json:"web_url"`
	PublishedDate string `json:"pub_date"`
	Headline      string `json:"headline"`
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// DataCrawlerSingle - get data for single URL
func DataCrawlerSingle(structArchive schema.Archive, archiveURL string, apiKey string, year string, month string) {
	url := strings.ReplaceAll(archiveURL, "{year}", year)
	url = strings.ReplaceAll(url, "{month}", month)
	url = strings.ReplaceAll(url, "{yourkey}", apiKey)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)

		duration := time.Duration(20) * time.Second
		time.Sleep(duration)

		// Closure function call.
		DataCrawlerSingle(structArchive, archiveURL, apiKey, year, month)
	}

	data, _ := ioutil.ReadAll(response.Body)

	if strings.Contains(string(data), `"detail":{"errorcode":"policies.ratelimit.QuotaViolation"}`) {

		fmt.Println("Repeating...")

		duration := time.Duration(5) * time.Second
		time.Sleep(duration)

		// Closure function call.
		DataCrawlerSingle(structArchive, archiveURL, apiKey, year, month)

	}

	pathConfig := configurations.GetPathConfigurations()

	archivePath := path.Join(pathConfig.RawCorpusRoot, pathConfig.ArchivePath, year)

	archiveFile := path.Join(archivePath, month+".json")

	fErr := ioutil.WriteFile(archiveFile, data, 0644)
	check(fErr)

	return
}

// ArchiveMetaCleaner - function to clean all corpus data
func ArchiveMetaCleaner() {
	archiveStructAll := []schema.Archive{}
	structArchive := schema.Archive{}

	pathconfiguration := configurations.GetPathConfigurations()

	// schema.DBConn.Where(&schema.Archive{Status: true, IsCleaned: false}).Find(&archiveStructAll)
	// Above code not working properly when applying multiple filters

	schema.DBConn.Model(&structArchive).Where("status = ?", true).Where("is_cleaned = ?", false).Find(&archiveStructAll)

	for _, row := range archiveStructAll {

		filename := pathconfiguration.RawCorpusRoot + "/" + pathconfiguration.ArchivePath + "/" + strconv.Itoa(row.Year) + "/" + strconv.Itoa(row.Month) + ".json"
		fmt.Println(filename)

		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		corpusByte, err := ioutil.ReadAll(file)

		cleanerResponse := &CleanerResponse{}

		// Decode JSON into our map
		jsonErr := json.Unmarshal(corpusByte, &cleanerResponse)
		if jsonErr != nil {
			fmt.Println(jsonErr)

			// Crawl data again
			NYObj := configurations.GetNYTimesAPIConfig()
			credentials := NYObj.Credentials
			apiEndpoints := NYObj.APIEndpoints

			apiKey := credentials["api_key"]
			archiveURL := apiEndpoints["archive"]

			go DataCrawlerSingle(row, archiveURL, apiKey, strconv.Itoa(row.Year), strconv.Itoa(row.Month))

			// Clean Crawled data after 1 hour / after all data cleaned. Continue required.
			continue

		}

		docData := cleanerResponse.Response["docs"]

		var monthData []interface{}

		for _, data := range docData.([]interface{}) {

			cleanedCorpus := CleanedCorpus{}

			for key, value := range data.(map[string]interface{}) {

				if key == "abstract" {

					if value != nil {
						cleanedCorpus.Abstract = value.(string)
					}

				} else if key == "web_url" {

					if value != nil {
						cleanedCorpus.WebURL = value.(string)
					}

				} else if key == "pub_date" {

					if value != nil {
						cleanedCorpus.PublishedDate = value.(string)
					}

				} else if key == "headline" {
					_, typeChecker := value.([]interface{})

					if typeChecker == true {

						cleanedCorpus.Headline = ""

					} else {

						for k, v := range value.(map[string]interface{}) {

							if k == "main" {
								if v != nil {
									cleanedCorpus.Headline = v.(string)
								}

							}

						}
					}

				}

			}

			monthData = append(monthData, cleanedCorpus)
			// Continuing next record in this month

		}

		cleanedFilePath := pathconfiguration.CleanedCorpusRoot + "/" + pathconfiguration.ArchivePath + "/" + strconv.Itoa(row.Year)
		os.MkdirAll(cleanedFilePath, os.ModePerm)
		cleanedFilename := cleanedFilePath + "/" + strconv.Itoa(row.Month) + ".json"

		// Write cleaned corpus into json file
		structJSON, _ := json.Marshal(monthData)
		err = ioutil.WriteFile(cleanedFilename, structJSON, 0644)

		// Changing the cleaned status to database entry.
		schema.DBConn.Model(&structArchive).Where("id = ?", row.ID).Update("is_cleaned", true)

		fmt.Println("Done.")

	}

	fmt.Println("You have completed Archive corpus cleaning. Next cleaning will start after 1 hour.")
	// Sleep 1 Hour and Run again to check if any new records found
	duration := time.Duration(1) * time.Hour
	time.Sleep(duration)

	// Closure function call to keep the function alive and checking for new records
	ArchiveMetaCleaner()

}
