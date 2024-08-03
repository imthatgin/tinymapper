package tinymapper_test

import (
	"fmt"
	"github.com/imthatgin/tinymapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStructureUser struct {
	Id       uint
	Money    float64
	Username string

	// When mapping, these should get turned into a display name property
	Firstname string
	Lastname  string

	PinCode      uint
	PasswordHash string
}

type TestStructureUserDTO struct {
	Id          uint
	Money       int
	Username    string
	DisplayName string
}

type TestStructureWithSubType struct {
	Id         uint
	Owner      TestStructureUser
	Dependents []TestStructureUser
}

type TestStructureWithSubTypeDTO struct {
	StructId   uint
	Owner      TestStructureUserDTO
	Dependents []TestStructureUserDTO
}

func Test_MapSingleStruct(t *testing.T) {
	m := tinymapper.New()

	tinymapper.Register(m, func(user TestStructureUser, dto *TestStructureUserDTO) {
		dto.DisplayName = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
	})

	user := TestStructureUser{
		Id:           121,
		Money:        55.5,
		Username:     "test_user_121",
		Firstname:    "Test",
		Lastname:     "User",
		PinCode:      9939,
		PasswordHash: "aaabbbccc",
	}

	dto, err := tinymapper.To[TestStructureUserDTO](m, user)
	assert.NoError(t, err)
	assert.Equal(t, dto.DisplayName, fmt.Sprintf("%s %s", user.Firstname, user.Lastname))
	assert.Equal(t, dto.Username, user.Username)
	assert.Equal(t, 0, dto.Money)

}

func Test_MapListOfStructs(t *testing.T) {
	m := tinymapper.New()

	tinymapper.Register(m, func(user TestStructureUser, dto *TestStructureUserDTO) {
		dto.DisplayName = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
	})

	users := []TestStructureUser{
		{
			Id:           121,
			Money:        55.5,
			Username:     "test_user_121",
			Firstname:    "Test",
			Lastname:     "User",
			PinCode:      9939,
			PasswordHash: "aaabbbccc",
		},
	}

	mappedUsers, err := tinymapper.ArrayTo[TestStructureUserDTO](m, users)
	assert.NoError(t, err)
	for i, dto := range mappedUsers {
		assert.Equal(t, fmt.Sprintf("%s %s", users[i].Firstname, users[i].Lastname), dto.DisplayName)
		assert.Equal(t, users[i].Username, dto.Username)
		assert.Equal(t, 0, dto.Money)
	}
}

func Test_MapSingleStructWithSubType(t *testing.T) {
	m := tinymapper.New()

	tinymapper.Register(m, func(user TestStructureUser, dto *TestStructureUserDTO) {
		dto.DisplayName = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
	})

	tinymapper.Register(m, func(ts TestStructureWithSubType, dto *TestStructureWithSubTypeDTO) {
		dto.StructId = ts.Id
		dto.Owner, _ = tinymapper.To[TestStructureUserDTO](m, ts.Owner)
		dto.Dependents, _ = tinymapper.ArrayTo[TestStructureUserDTO](m, ts.Dependents)
	})

	original := TestStructureWithSubType{
		Owner: TestStructureUser{
			Id:           121,
			Username:     "test_user_121",
			Firstname:    "Test",
			Lastname:     "User",
			PinCode:      9939,
			PasswordHash: "aaabbbccc",
		},
		Dependents: []TestStructureUser{
			{
				Id:           555,
				Username:     "555",
				Firstname:    "SubTest",
				Lastname:     "SubUser",
				PinCode:      1111,
				PasswordHash: "001",
			},
		},
	}

	dto, err := tinymapper.To[TestStructureWithSubTypeDTO](m, original)
	assert.NoError(t, err)
	assert.Equal(t, original.Id, dto.StructId)

	assert.Equal(t, 0, dto.Owner.Money)
	assert.Equal(t, original.Owner.Id, dto.Owner.Id)
	assert.Equal(t, original.Owner.Username, dto.Owner.Username)
	assert.Equal(t, fmt.Sprintf("%s %s", original.Owner.Firstname, original.Owner.Lastname), dto.Owner.DisplayName)

	assert.Len(t, dto.Dependents, len(original.Dependents))

	for i, dependent := range original.Dependents {
		assert.Equal(t, fmt.Sprintf("%s %s", dependent.Firstname, dependent.Lastname), dto.Dependents[i].DisplayName)
	}
}
