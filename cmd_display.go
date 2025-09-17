package krap

import (
	"fmt"

	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/io"
)

func DisplayError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func DisplayData[T any](data *T, rq *Request, err error) {
	if rq != nil {
		fmt.Println(rq.Output())
	}
	if err == nil {
		if data == nil {
			fmt.Println("Data: null")
		} else {
			output, err := io.StringifyIndentedJSON(data)
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

func DisplayOutput(rq *Request, err error) {
	if rq != nil {
		fmt.Println(rq.Output())
	}
	if err == nil {
		fmt.Println(okMessage)
	} else {
		DisplayError(err)
	}
}

func DisplayList[T any](list *ds.List[T], rq *Request, err error) {
	if rq != nil {
		fmt.Println(rq.Output())
	}
	if err == nil {
		for i, item := range list.Items {
			fmt.Printf("%d: %v\n", i+1, item)
		}
		fmt.Println("Count:", list.Count)
	} else {
		DisplayError(err)
	}
}
