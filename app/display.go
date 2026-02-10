package app

import (
	"fmt"

	"github.com/roidaradal/rdb/ze"
)

// Prints request logs and error
func DisplayOutput(rq *ze.Request, err error) {
	if rq != nil {
		fmt.Println(rq.Output())
	}
	if err == nil {
		fmt.Println(okMessage)
	} else {
		DisplayError(err)
	}
}

// Prints error
func DisplayError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
