package service_test

import (
	"fmt"

	"testing"

	"github.com/sea-auca/auc-auth/user/service"
	"github.com/stretchr/testify/assert"
)

func TestPermissionsNew(t *testing.T) {
	tt := []struct {
		in  []service.Permissions
		out service.Permissions
	}{
		{in: []service.Permissions{}, out: service.None},
		{in: []service.Permissions{service.AUCA_student}, out: service.AUCA_student},
		{in: []service.Permissions{service.SEA_development, service.AUCA_student}, out: service.SEA_development | service.AUCA_student},
	}

	for i, test := range tt {
		perm := service.NewPermissions(test.in...)
		assert.Equalf(t, test.out, *perm, "Expected %v, got %v, in test case #%v", test.out, *perm, i)
	}
}

func TestAssignNewPermission(t *testing.T) {
	tt := []struct {
		in  service.Permissions
		arg service.Permissions
		out service.Permissions
	}{
		{
			in:  *service.NewPermissions(),
			arg: service.AUCA_student,
			out: *service.NewPermissions(service.AUCA_student),
		},
		{
			in:  *service.NewPermissions(service.AUCA_student),
			arg: service.SEA_development,
			out: *service.NewPermissions(service.AUCA_student, service.SEA_development),
		},
	}

	for i, test := range tt {
		test.in.Assing(test.arg)
		assert.Equalf(t, test.out, test.in, "Expected %v, got %v, in test case #%v", test.out, test.in, i)
	}
}

func TestRevoke(t *testing.T) {
	tt := []struct {
		in  service.Permissions
		arg service.Permissions
		out service.Permissions
	}{
		{in: *service.NewPermissions(service.AUCA_student), arg: service.AUCA_student, out: service.None},
		{in: service.None, arg: service.AUCA_student, out: service.None},
		{in: *service.NewPermissions(service.SEA_development), arg: service.AUCA_student, out: service.SEA_development},
	}

	for i, test := range tt {
		test.in.Revoke(test.arg)
		assert.Equalf(t, test.out, test.in, "Expected %v, got %v, in test case #%v", test.out, test.in, i)
	}
}

func ExampleNewPermissions() {
	perm := service.NewPermissions()
	perm2 := service.NewPermissions(service.AUCA_student)
	fmt.Println(perm.Sprint(), perm2.Sprint())
	//Output:
	//
	//Auca student
}
