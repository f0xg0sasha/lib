package rest

import log "github.com/sirupsen/logrus"

func logError(handlerName string, err error) {
	log.WithFields(log.Fields{
		"handler": handlerName,
		"error":   err,
	}).Error()
}
