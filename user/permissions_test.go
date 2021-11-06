package user_test

import (
	"fmt"
	"sea/auth/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionsNew(t *testing.T) {
	tt := []struct {
		in  []user.Permissions
		out user.Permissions
	}{
		{in: []user.Permissions{}, out: user.None},
		{in: []user.Permissions{user.AUCA_student}, out: user.AUCA_student},
		{in: []user.Permissions{user.SEA_development, user.AUCA_student}, out: user.SEA_development | user.AUCA_student},
	}

	for i, test := range tt {
		perm := user.NewPermissions(test.in...)
		assert.Equalf(t, test.out, *perm, "Expected %v, got %v, in test case #%v", test.out, *perm, i)
	}
}

func TestAssignNewPermission(t *testing.T) {
	tt := []struct {
		in  user.Permissions
		arg user.Permissions
		out user.Permissions
	}{
		{
			in:  *user.NewPermissions(),
			arg: user.AUCA_student,
			out: *user.NewPermissions(user.AUCA_student),
		},
		{
			in:  *user.NewPermissions(user.AUCA_student),
			arg: user.SEA_development,
			out: *user.NewPermissions(user.AUCA_student, user.SEA_development),
		},
	}

	for i, test := range tt {
		test.in.Assing(test.arg)
		assert.Equalf(t, test.out, test.in, "Expected %v, got %v, in test case #%v", test.out, test.in, i)
	}
}

func TestRevoke(t *testing.T) {
	tt := []struct {
		in  user.Permissions
		arg user.Permissions
		out user.Permissions
	}{
		{in: *user.NewPermissions(user.AUCA_student), arg: user.AUCA_student, out: user.None},
		{in: user.None, arg: user.AUCA_student, out: user.None},
		{in: *user.NewPermissions(user.SEA_development), arg: user.AUCA_student, out: user.SEA_development},
	}

	for i, test := range tt {
		test.in.Revoke(test.arg)
		assert.Equalf(t, test.out, test.in, "Expected %v, got %v, in test case #%v", test.out, test.in, i)
	}
}

func ExampleNewPermissions() {
	perm := user.NewPermissions()
	perm2 := user.NewPermissions(user.AUCA_student)
	fmt.Println(perm.Sprint(), perm2.Sprint())
	//Output:
	//
	//Auca student
}
