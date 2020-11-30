package html_controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWeChatAppId(t *testing.T) {
	appId := getWeChatAppId(1)
	assert.NotEmpty(t, appId)

	appId = getWeChatAppId(0)
	assert.Equal(t, "", appId)
}
