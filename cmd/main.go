package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jkittell/array"
	"github.com/jkittell/mediastreamvalidator"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type db struct {
	Database *array.Array[*mediastreamvalidator.StreamValidator]
}

// contents array to store contents data.
var contents db

func getInstance() db {
	if contents.Database == nil {
		lock.Lock()
		defer lock.Unlock()
		if contents.Database == nil {
			log.Println("creating new contents database")
			contents.Database = array.New[*mediastreamvalidator.StreamValidator]()
		}
	}
	return contents
}

func main() {
	router := gin.Default()
	router.GET("/contents", getContents)
	router.GET("/contents/:id", getContentByID)
	router.POST("/contents", postContents)

	err := router.Run(":3001")
	if err != nil {
		log.Println(err)
	}
}

// getContents responds with the list of all contents as JSON.
func getContents(c *gin.Context) {
	getInstance()
	c.JSON(http.StatusOK, contents.Database)
}

// postContents adds content from JSON received in the request body.
func postContents(c *gin.Context) {
	newContent := &mediastreamvalidator.StreamValidator{
		Id:         uuid.New().String(),
		URL:        "",
		Validation: mediastreamvalidator.ValidationInfo{},
		Status:     "queued",
		StartTime:  time.Now().UTC(),
		EndTime:    time.Time{},
	}

	// Call BindJSON to bind the received JSON to
	// newContent.
	if err := c.BindJSON(&newContent); err != nil {
		log.Println(err)
		return
	}

	// Add the new content to the array.
	getInstance()
	contents.Database.Push(newContent)
	c.JSON(http.StatusCreated, newContent)
	go validate(newContent)
}

// getContentByID locates the content whose ID value matches the id
// parameter sent by the client, then returns that content as a response.
func getContentByID(c *gin.Context) {
	getInstance()
	id := c.Param("id")

	// Loop through the list of contents, looking for
	// content whose ID value matches the parameter.
	for i := 0; i < contents.Database.Length(); i++ {
		j := contents.Database.Lookup(i)
		if j.Id == id {
			c.JSON(http.StatusOK, j)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "content not found"})
}

func validate(content *mediastreamvalidator.StreamValidator) {
	if strings.HasSuffix(content.URL, ".mpd") {
		log.Println("skip mediastreamvalidator for dash", nil)
		content.Status = "skipped"
	} else {
		content.Status = "processing"
		// verify mediastreamvalidator available on server
		exe, err := exec.LookPath("mediastreamvalidator")
		if err != nil {
			log.Println("mediastreamvalidator is not available", err)
			content.Status = "error"
			return
		} else {
			// create tmp json file
			f, err := os.CreateTemp("/tmp", "msv-")
			if err != nil {
				log.Println("error creating temp file for mediastreamvalidator to use", err)
				content.Status = "error"
				return
			}
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					log.Println(err)
					content.Status = "error"
					return
				}
			}(f.Name())

			var arguments []string
			arguments = append(arguments, "-t")
			arguments = append(arguments, "30")
			arguments = append(arguments, fmt.Sprintf("--validation-data-path=%s", f.Name()))
			arguments = append(arguments, content.URL)

			// run mediastreamvalidator
			cmd := exec.Command(exe, arguments...)
			err = cmd.Run()
			if err != nil {
				log.Println("error running mediastreamvalidator", err)
				content.Status = "error"
				return
			}
			// parse the json file
			b, err := os.ReadFile(f.Name())
			if err != nil {
				log.Println(fmt.Sprintf("error reading json data file - %s", f.Name()), err)
				content.Status = "error"
				return
			}

			err = json.Unmarshal(b, &content.Validation)
			if err != nil {
				log.Println("error unmarshalling json data to msv struct", err)
				content.Status = "error"
				return
			} else {
				content.Status = "completed"
			}
		}
	}
}
