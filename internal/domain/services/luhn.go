package services

const (
	asciiZero = 48
	asciiTen  = 57
)

type LuhnDomainService struct {
}

func (l *LuhnDomainService) Validate(number string) bool {
	p := len(number) % 2
	sum, ok := l.calculateLuhnSum(number, p)
	// If the total modulo 10 is not equal to 0, then the number is invalid.
	return ok && sum%10 == 0
}

func (*LuhnDomainService) calculateLuhnSum(number string, parity int) (int64, bool) {
	var sum int64
	for i, d := range number {
		if d < asciiZero || d > asciiTen {
			return 0, false
		}

		d = d - asciiZero
		// Double the value of every second digit.
		if i%2 == parity {
			d *= 2
			// If the result of this doubling operation is greater than 9.
			if d > 9 {
				// The same final result can be found by subtracting 9 from that result.
				d -= 9
			}
		}

		// Take the sum of all the digits.
		sum += int64(d)
	}

	return sum, true
}
