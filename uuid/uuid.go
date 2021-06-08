package uuid

import guuid "github.com/google/uuid"

// New will return a new UUID.v4 represented as a string
// or the empty string when generation fails
func NewString() string {
	u, err := guuid.NewRandom()
	if err != nil {
		return ""
	}

	return u.String()
}
