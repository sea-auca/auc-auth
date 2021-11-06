package user_test

import "testing"

// TODO: Write down these tests
// MAKE THE USEFUL
// these test will be used as a busniess logic
// requirements in future
// add test if you want

func TestRegisterUserSucessful(t *testing.T) {
	// init repos and tests

	//use simple auca email and register
	//expect a new user record and new verification link to be generated and sent
}

func TestRegisterUserIncorrectEmail(t *testing.T) {
	//init repos and tests

	//incorrect email will revert process before any insertion
}

// Registration may fail due to problems with email of verification link gen
// Do not know how to test it

func TestUpdateUserSuccesful(t *testing.T) {
	//expect a sucessful update of data
}

// Test update on a set of invalid data
// it should always fail and never perform any database activity
func TestUpdateUserInvalidData(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
func TestDeactivateAccountSuccessful(t *testing.T) {
	// Account should be set as innactive sucessfully
}

func TestVerifyUser(t *testing.T) {
	//Test should verify that if link code exists
	// users verify flag is set correctly
}

func TestSettingNewPasswordSuccess(t *testing.T) {
	// check the verification code and set new password
}

// check for different combinations of codes and passwords
func TestSettingNewPasswordError(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
