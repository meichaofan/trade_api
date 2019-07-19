package common

import "truxing/commons/log"

func ErrorHandler(err error) {
	if err != nil {
		log.Debugf("err: %s", err)
	}
}