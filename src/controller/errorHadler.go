package controller

import (
	"fmt"
	"github.com/lib/pq"
	"log"
)

func ShowError(err error) {
	switch e := err.(type) {
	case *pq.Error:
		switch e.Code {
		case "23502":
			// not-null constraint violation
			log.Printf(fmt.Sprint("Some required data was left out:\n", e.Message))
			return

		case "23503":
		// foreign key violation
		case "DELETE":
			log.Printf(fmt.Sprint("This record can’t be deleted because another record refers to it:\n", e.Message))
			return

		case "23505":
			// unique constraint violation
			log.Printf(fmt.Sprint("This record contains duplicated data that conflicts with what is already in the database:\n", e.Message))
			return

		case "23514":
			// check constraint violation
			log.Printf(fmt.Sprint("This record contains inconsistent or out-of-range data:\n", e.Message))
			return

		default:
			msg := e.Message
			if d := e.Detail; d != "" {
				msg += "\n\n" + d
			}
			if h := e.Hint; h != "" {
				msg += "\n\n" + h
			}
			log.Println(msg)
			return
		}
	}
}
