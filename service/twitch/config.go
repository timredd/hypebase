package twitch

var (
	// AppClientID is the client ID for Hypebase
	AppClientID = "3domgfkm0dgtw3yzp334x6514nu2vy"

	// RedirectURI is the URL users are taken to after logging in
	RedirectURI = "http://localhost"

	// DefaultScopes are the default permissions granted to app for user access
	DefaultScopes = []string{"user:read:email", "channel:manage:broadcast"}
)
