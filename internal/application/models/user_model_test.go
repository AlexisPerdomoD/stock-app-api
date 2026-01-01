package models

import "testing"

func TestUserLoginDTO_GetPasswordBytesAndClean(t *testing.T) {
	dto := &UserLoginDTO{
		Username: "test@test.com",
		Password: "secret123",
	}

	// first call
	pwd, err := dto.GetPasswordBytesAndClean()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(pwd) != "secret123" {
		t.Fatalf("unexpected password bytes: %s", pwd)
	}
	if dto.Password != "" {
		t.Fatalf("password was not cleaned")
	}

	// second call must fail
	_, err = dto.GetPasswordBytesAndClean()
	if err == nil {
		t.Fatal("expected error on second password consumption")
	}
}

func TestRegisterUserDTO_GetPasswordBytesAndClean(t *testing.T) {
	dto := &RegisterUserDTO{
		Username:  "test@test.com",
		Firstname: "John",
		Lastname:  "Doe",
		Password:  "secret123",
	}

	// first call
	pwd, err := dto.GetPasswordBytesAndClean()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(pwd) != "secret123" {
		t.Fatalf("unexpected password bytes: %s", pwd)
	}
	if dto.Password != "" {
		t.Fatalf("password was not cleaned")
	}

	// second call must fail
	_, err = dto.GetPasswordBytesAndClean()
	if err == nil {
		t.Fatal("expected error on second password consumption")
	}
}
