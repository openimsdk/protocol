package datautil

import "unicode"

// IsLegalUserID check if str is legal userID(Only contain letter/number/_).
func IsLegalUserID(str string) bool {
	for _, r := range str {
		if !IsAlphanumeric(r) && r != '_' {
			return false
		}
	}
	return true
}

// IsAlphanumeric check if b is a letter or number
func IsAlphanumeric(b rune) bool {
	if !unicode.IsLetter(b) && !unicode.IsDigit(b) {
		return false
	}
	return true
}
