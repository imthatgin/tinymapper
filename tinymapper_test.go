package tinymapper_test

import (
	"fmt"
	"github.com/imthatgin/tinymapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStructureUser struct {
	Id       uint
	Username string

	// When mapping, these should get turned into a display name property
	Firstname string
	Lastname  string

	PinCode      uint
	PasswordHash string
}

type TestStructureUserDTO struct {
	Id          uint
	Username    string
	DisplayName string
}

func Test_MapsSourcePairsToSingleDestinationField(t *testing.T) {
	m := tinymapper.New()

	// Define the mapping
	_ = tinymapper.Register(m, func(user TestStructureUser, dto *TestStructureUserDTO) {
		dto.DisplayName = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
	})

	// Test if the mapping can work
	user := TestStructureUser{
		Id:           121,
		Username:     "test_user_121",
		Firstname:    "Test",
		Lastname:     "User",
		PinCode:      9939,
		PasswordHash: "aaabbbccc",
	}

	dto := tinymapper.To[TestStructureUserDTO](m, user)
	assert.IsTypef(t, &TestStructureUserDTO{}, dto, "Expected returned value to be mapped type")
	assert.Equal(t, dto.DisplayName, fmt.Sprintf("%s %s", user.Firstname, user.Lastname))
}
