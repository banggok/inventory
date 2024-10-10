package utility

import log "github.com/sirupsen/logrus"

// LogError logs an error with specific action and context
func LogError(action, name string, err error) {
	log.WithFields(log.Fields{
		"error": err,
		"name":  name,
	}).Error(action)
}

// LogSuccess logs successful operations with action and context
func LogSuccess(action string, id uint, name string) {
	log.WithFields(log.Fields{
		"id":   id,
		"name": name,
	}).Info(action)
}
