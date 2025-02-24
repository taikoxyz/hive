package params

import "os"

func initJWT() {
	err := os.WriteFile("/tmp/jwt.hex", []byte("c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7"), 0644)
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("JWT_SECRET", "/tmp/jwt.hex")
}
