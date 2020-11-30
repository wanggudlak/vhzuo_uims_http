package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateCSRFToken(t *testing.T) {
	token, err := GenerateCSRFToken()
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	if !VerifyCSRFToken(token) {
		t.Log(err)
		assert.NotNil(t, err)
	}
}
