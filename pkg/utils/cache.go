package utils

// hotUsers simulates active users whose data is likely cached.
var hotUsers = map[int64]bool{1: true, 2: true, 3: true}

// UseCache returns true if the user should use the cache layer.
func UseCache(userID int64) bool {
	return hotUsers[userID]
}
