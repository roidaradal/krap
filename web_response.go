package krap

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roidaradal/fn/str"
)

const okMessage = "OK"

type actionResponse struct {
	Success bool
	Message string
}

type dataResponse struct {
	Data    any
	Message string
}

type ResponseType struct {
	SendErrorFn func(*gin.Context, *Request, error)
}

var WebData = &ResponseType{
	SendErrorFn: SendDataError,
}

var WebAction = &ResponseType{
	SendErrorFn: SendActionError,
}

func SendActionResponse(c *gin.Context, rq *Request, err error) {
	// TODO: APILog (output)
	output := getOutput(rq, err)
	fmt.Println("Output:\n" + output) // Temporary

	if err == nil {
		// url := c.Request.URL.Path
		// TODO: APILog(url, 200)
		c.JSON(http.StatusOK, actionResponse{
			Success: true,
			Message: okMessage,
		})
	} else {
		sendActionError(c, err)
	}
}

func SendActionError(c *gin.Context, rq *Request, err error) {
	// TODO: APILog (output)
	output := getOutput(rq, err)
	fmt.Println("Output:\n" + output) // Temporary

	sendActionError(c, err)
}

func SendDataResponse[T any](c *gin.Context, data *T, rq *Request, err error) {
	// TODO: APILog (output)
	output := getOutput(rq, err)
	fmt.Println("Output:\n" + output) // Temporary

	if err == nil {
		// url := c.Request.URL.Path
		if data != nil {
			// TODO: APILog (url, 200)
			c.JSON(http.StatusOK, dataResponse{
				Data:    data,
				Message: okMessage,
			})
		} else {
			// TODO: APILog (url, 200, Data: "null")
			c.JSON(http.StatusOK, dataResponse{
				Data:    nil,
				Message: okMessage,
			})
		}
	} else {
		sendDataError(c, err)
	}
}

func SendDataError(c *gin.Context, rq *Request, err error) {
	// TODO: APILog(output)
	output := getOutput(rq, err)
	fmt.Println("Output:\n" + output) // Temporary

	sendDataError(c, err)
}

func sendActionError(c *gin.Context, err error) {
	// url := c.Request.URL.Path
	message, ok := publicErrorMessage(err)
	if ok {
		// TODO: APILog(err, url, 400)
		c.JSON(http.StatusBadRequest, actionResponse{
			Success: false,
			Message: message,
		})
	} else {
		// TODO: APILog(err, url, 500)
		c.JSON(http.StatusInternalServerError, nil)
	}

}

func sendDataError(c *gin.Context, err error) {
	// url := c.Request.URL.Path
	message, ok := publicErrorMessage(err)
	if ok {
		// TODO: APILog(err, url, 400)
		c.JSON(http.StatusBadRequest, dataResponse{
			Data:    nil,
			Message: message,
		})
	} else {
		// TODO: APILog (err, url, 500)
		c.JSON(http.StatusInternalServerError, nil)
	}
}

func getOutput(rq *Request, err error) string {
	output := make([]string, 0)
	if rq != nil {
		output = append(output, rq.Output())
	}
	if err != nil {
		output = append(output, fmt.Sprintf("Error: %s", err.Error()))
	}
	return strings.Join(output, "\n")
}

func publicErrorMessage(err error) (string, bool) {
	msg := err.Error()
	if strings.HasPrefix(msg, "public: ") {
		parts := str.CleanSplit(msg, ":")
		return parts[1], true
	} else {
		return "Unexpected error during request", false
	}
}
