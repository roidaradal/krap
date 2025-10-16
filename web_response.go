package krap

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roidaradal/fn"
	"github.com/roidaradal/rdb/ze"
)

// Response for action web requests
type actionResponse struct {
	Success bool
	Message string
}

// Response for data web requests
type dataResponse struct {
	Data    any
	Message string
}

type ResponseType struct {
	SendErrorFn func(*gin.Context, *ze.Request, error)
}

var (
	WebAction = &ResponseType{
		SendErrorFn: SendActionError,
	}
	WebData = &ResponseType{
		SendErrorFn: SendDataError,
	}
)

// Sends actionResponse
func SendActionResponse(c *gin.Context, rq *ze.Request, err error) {
	// TODO: APILog output
	fmt.Println("Output:\n" + getOutput(rq, err)) // temporary

	if err == nil {
		// url := c.Request.URL.Path
		// TODO: APILog url, status
		status := ze.OK200
		if rq != nil {
			status = fn.Ternary(rq.Status == 0, ze.OK200, rq.Status)
		}
		c.JSON(status, actionResponse{
			Success: true,
			Message: okMessage,
		})
	} else {
		sendActionError(c, rq, err)
	}
}

// Sends dataResponse
func SendDataResponse[T any](c *gin.Context, data *T, rq *ze.Request, err error) {
	// TODO: APILog output
	fmt.Println("Output:\n" + getOutput(rq, err)) // temporary

	if err == nil {
		// url := c.Request.URL.Path
		// TODO: APILog url, status, include data in logs?
		status := ze.OK200
		if rq != nil {
			status = fn.Ternary(rq.Status == 0, ze.OK200, rq.Status)
		}
		c.JSON(status, dataResponse{
			Data:    data,
			Message: okMessage,
		})

	} else {
		sendDataError(c, rq, err)
	}
}

// Sends actionResponse with given error
func SendActionError(c *gin.Context, rq *ze.Request, err error) {
	// TODO: APILog output
	fmt.Println("Output:\n" + getOutput(rq, err)) // temporary
	sendActionError(c, rq, err)
}

// Sends dataResponse with given error
func SendDataError(c *gin.Context, rq *ze.Request, err error) {
	// TODO: APILog output
	fmt.Println("Output:\n" + getOutput(rq, err)) // temporary
	sendDataError(c, rq, err)
}

// Common: sends an actionResponse with the given error
func sendActionError(c *gin.Context, rq *ze.Request, err error) {
	// url := c.Request.URL.Path
	status, message := getStatusMessage(rq, err)
	c.JSON(status, actionResponse{
		Success: false,
		Message: message,
	})
}

// Common: sends a dataResponse with the given error
func sendDataError(c *gin.Context, rq *ze.Request, err error) {
	// url := c.Request.URL.Path
	status, message := getStatusMessage(rq, err)
	c.JSON(status, dataResponse{
		Data:    nil,
		Message: message,
	})
}

// Combines the request output and the error message into one string, joined by newline
func getOutput(rq *ze.Request, err error) string {
	output := make([]string, 0)
	if rq != nil {
		output = append(output, rq.Output())
	}
	if err != nil {
		output = append(output, fmt.Sprintf("Error: %s", err.Error()))
	}
	return strings.Join(output, "\n")
}

// Get error message and status code
func getStatusMessage(rq *ze.Request, err error) (int, string) {
	message, ok := publicErrorMessage(err)
	status := fn.Ternary(ok, ze.Err400, ze.Err500)
	if rq != nil {
		status = fn.Ternary(rq.Status == 0, status, rq.Status)
	}
	return status, message
}
