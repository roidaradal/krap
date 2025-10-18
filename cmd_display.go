package krap

import (
	"fmt"

	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/rdb/ze"
)

// Prints error
func DisplayError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// Prints data, request logs, and error
func DisplayData[T any](data *T, rq *ze.Request, err error) {
	if rq != nil {
		fmt.Println(rq.Output())
	}
	if err == nil {
		if data == nil {
			fmt.Println("Data: null")
		} else {
			output, err := str.IndentedJSON(data)
			if err == nil {
				fmt.Println(output)
			} else {
				DisplayError(err)
			}
		}
	} else {
		DisplayError(err)
	}
}

// Prints list items, request logs, and error
func DisplayList[T any](list *ds.List[*T], rq *ze.Request, err error) {
	if rq != nil {
		fmt.Println(rq.Output())
	}
	if err == nil {
		for i, item := range list.Items {
			if item == nil {
				fmt.Printf("%d: nil\n", i+1)
			} else {
				fmt.Printf("%d: %v\n", i+1, *item)
			}
		}
		fmt.Println("Count:", list.Count)
	} else {
		DisplayError(err)
	}
}

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
