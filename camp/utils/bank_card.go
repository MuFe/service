package utils

// BankCardReplace3 银行卡号屏蔽，前三后四
func BankCardReplace3(phone string) string {
	if phone == "" {
		return ""
	}
	l := len(phone)
	ph := "********"
	if l < 4 {
		return ph
	}
	if l < 12 {
		return phone[:4] + ph
	}
	return phone[:4] + ph + phone[12:]
}
