package utils

import (
	"fmt"
	"testing"
)

func TestVerify(t *testing.T) {
	fmt.Println(VerifyPhoneFormat("1562237331"))
}

func TestReplace(t *testing.T) {
	phone := "15626446447"
	fmt.Println(phone)
	fmt.Println(PhoneReplace4(phone))
}
func TestReplace2(t *testing.T) {
	fmt.Println(PhoneReplace6("133"))
	fmt.Println(PhoneReplace6("15626446447"))
	fmt.Println(PhoneReplace6("156264464471233"))
}
