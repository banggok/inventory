package utility

import log "github.com/sirupsen/logrus"

// LogError logs an error with specific action and context
func LogError(action, name string, err error) {
	log.WithFields(log.Fields{
		"error": err,
		"name":  name,
	}).Error(action)
}

// LogSuccess logs a success message with dynamic data
func LogSuccess(message string, data ...interface{}) {
	log.WithFields(log.Fields{
		"data": data,
	}).Info(message)
}
