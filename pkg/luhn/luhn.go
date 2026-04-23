// Package luhn реализует проверку номера по алгоритму Луна.
package luhn

// Valid проверяет, является ли строка корректным номером по алгоритму Луна.
// Возвращает false, если строка пустая или содержит нецифровые символы.
func Valid(number string) bool {
	if len(number) == 0 {
		return false
	}

	sum := 0
	nDigits := len(number)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		ch := number[i]
		if ch < '0' || ch > '9' {
			return false
		}
		digit := int(ch - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
