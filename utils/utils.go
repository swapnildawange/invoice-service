package utils

import (
	"invoice_service/security"
	"testing"

	"github.com/spf13/viper"
)

func InitPlatformJWT(t *testing.T, userId int, role int) string {
	var err error
	key := viper.GetString("JWTSECRET")

	platformJWT, err := security.GenerateJWT(key, userId, role)
	if err != nil {
		t.Errorf("Failed to encode platformJWT for test")
	}
	return platformJWT
}
