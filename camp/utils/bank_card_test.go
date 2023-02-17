package utils

import (
	"fmt"
	"testing"
)

func TestBankCardReplace3(t *testing.T) {
	fmt.Println(BankCardReplace3("133"))
	fmt.Println(BankCardReplace3("15626446447"))
	fmt.Println(BankCardReplace3("156264464471233"))
}
