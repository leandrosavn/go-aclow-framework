package config

type Config struct {
	Context map[string]map[string]string
	Event   map[string]map[string]string
	Rest    map[string]map[string]string
}

func ConfigRoutes() *Config {
	config := Config{}
	config.Context = map[string]map[string]string{
		"public": {
			"login-event":         "",
			"check-token-event":   "",
			"refresh-token-event": "",
		},
		"protected": {
			"login-event":         "auth@user-login",
			"check-token-event":   "auth@user-validating-access",
			"refresh-token-event": "auth@user-refreshing-access",
		},
	}

	config.Event = map[string]map[string]string{
		// api
		"api@ping": {
			"auth": "public",
		},

		"api@authenticating": {
			"auth": "public",
		},

		//auth
		"auth@user-login": {
			"auth": "public",
		},
		"auth@user-refreshing-access": {
			"auth": "public",
		},
	}

	config.Rest = map[string]map[string]string{
		"routes": {
			"GET    /ping	  ": "api@ping",
			"POST   /userLogin":            "api@authenticating",
			"POST   /userRefreshingAccess": "auth@user-refreshing-access",
		},
	}

	return &config
}
