package env

import (
	"os"
	"time"
)

var (
	LISTEN_ADDRESS string
	CLIENT_ID      string
	CLIENT_SECRET  string
	USER_AGENT     string
	CLIENT_TIMEOUT time.Duration
)

func Init() {
	CLIENT_ID = os.Getenv("CLIENT_ID")
	if CLIENT_ID == "" {
		panic(newEnvMissingError("CLIENT_ID"))
	}

	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
	if CLIENT_SECRET == "" {
		panic(newEnvMissingError("CLIENT_SECRET"))
	}

	USER_AGENT = os.Getenv("USER_AGENT")
	if USER_AGENT == "" {
		panic(newEnvMissingError("USER_AGENT"))
	}

	LISTEN_ADDRESS = os.Getenv("SERVE_ADDRESS")
	if LISTEN_ADDRESS == "" {
		panic(newEnvMissingError("SERVE_ADDRESS"))
	}

	var ctstr string
	var err error
	ctstr = os.Getenv("CLIENT_TIMEOUT")
	if ctstr == "" {
		panic(newEnvMissingError("CLIENT_TIMEOUT"))
	}
	CLIENT_TIMEOUT, err = time.ParseDuration(ctstr)
	if err != nil {
		panic(newEnvInvalidError[time.Duration](
			"CLIENT_TIMEOUT",
			ctstr,
		))
	}
}
