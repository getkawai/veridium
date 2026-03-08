package services

import "strings"

func buildMemoryScopeKey(userID, sessionID string) string {
	userID = strings.TrimSpace(userID)
	sessionID = strings.TrimSpace(sessionID)

	switch {
	case userID == "" && sessionID == "":
		return ""
	case userID == "":
		return "anon:" + sessionID
	case sessionID == "":
		return userID + ":default"
	default:
		return userID + ":" + sessionID
	}
}
