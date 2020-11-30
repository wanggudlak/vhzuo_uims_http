package gjwt_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/pkg/gjwt"
	"uims/pkg/randc"
)

func TestJwt_New(t *testing.T) {
	j := gjwt.Jwt{}
	token, err := gjwt.CreateToken(&j)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	t.Log(token)
	j2 := gjwt.Jwt{}
	err = gjwt.Parse(token, &j2)
	assert.Nil(t, err)
	assert.Equal(t, j.Audience, j2.Audience)
}

func TestJwt_New2(t *testing.T) {
	type Jwt2 struct {
		gjwt.Jwt
		Test string
	}
	j := Jwt2{
		Test: "test",
	}
	j.SetIssue()
	j.SetTTL(0)
	j.SetAudience("test")
	token, err := gjwt.CreateToken(&j)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	t.Log(token)
	j2 := Jwt2{}
	err = gjwt.Parse(token, &j2)
	assert.Nil(t, err)
	assert.Equal(t, j.Audience, j2.Audience)
	t.Log(j2.Test)
}

func TestParse(t *testing.T) {
	type Jwt2 struct {
		gjwt.Jwt
		Test string
	}
	j := Jwt2{}
	err := gjwt.Parse("", &j)
	assert.NotNil(t, err)

	j2 := Jwt2{}
	err = gjwt.Parse(randc.UUID(), &j2)
	assert.NotNil(t, err)

	j3 := Jwt2{
		Test: "test",
	}

	token, err := gjwt.CreateToken(&j3)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	j4 := Jwt2{}
	err = gjwt.Parse(token, &j4)
	assert.Nil(t, err)
	assert.Equal(t, j3.Test, j4.Test)
}
