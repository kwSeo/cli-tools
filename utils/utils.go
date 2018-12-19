package utils

import "log"

func MustNotErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
