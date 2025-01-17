package errors

import "log"

func Must(err error, context string) {
	if err != nil {
		log.Fatal("%s: %v", context, err)
	}
}

func Normal(err error, context string) {
	if err != nil {
		log.Printf("%s: %v", context, err)
	}
}
