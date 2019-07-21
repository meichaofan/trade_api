package common

import "truxing/commons/log"

func ErrorHandler(err error) {
	if err != nil {
		log.Debugf("err: %s", err)
	}
}

func MessageHandler(what, str string) {
	log.Debugf("%s , this is a message I what to tell you : %s", what, str)
}
