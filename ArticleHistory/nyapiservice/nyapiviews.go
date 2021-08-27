package nyapiservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	configurations "../configurations"
	schema "../schema"
)

// GetMeta - Function to get meta data
func GetMeta(year string, month string) (interface{}, int) {
	structArchive := schema.Archive{}

	schema.DBConn.Model(&structArchive).Where("status = ?", true).Where("is_cleaned = ?", true).Where("year = ?", year).Where("month = ?", month).First(&structArchive)
	if structArchive.ID == 0 {
		return "Data not found for the year " + year + " and month " + month + ".", 400
	}

	pathconfiguration := configurations.GetPathConfigurations()

	cleanedFilePath := pathconfiguration.CleanedCorpusRoot + "/" + pathconfiguration.ArchivePath + "/" + year
	cleanedFilename := cleanedFilePath + "/" + month + ".json"

	var jsonStruct []interface{}

	rawBytes, _ := ioutil.ReadFile(cleanedFilename)

	err := json.Unmarshal(rawBytes, &jsonStruct)

	if err != nil {

		return "Something went wrong.", 400

	}

	return jsonStruct, 200
}

// GetArchiveMetaData - Returns Archive API details
func GetArchiveMetaData(c *gin.Context) {
	_, statusID := c.Get("User")

	year := c.Param("year")
	month := c.Param("month")

	if statusID == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return

	}

	data, code := GetMeta(year, month)

	c.JSON(code, gin.H{"data": data})
}
