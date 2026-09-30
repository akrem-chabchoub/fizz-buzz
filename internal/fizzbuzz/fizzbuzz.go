package fizzbuzz

import "strconv"

// Generate returns one string for each number from 1 to limit.
// A multiple of both numbers becomes the two words stuck together.
// When the two numbers are equal, every multiple takes that path.
func Generate(int1, int2, limit int, str1, str2 string) []string {
	out := make([]string, limit)
	both := str1 + str2
	for i := 1; i <= limit; i++ {
		switch {
		case i%int1 == 0 && i%int2 == 0:
			out[i-1] = both
		case i%int1 == 0:
			out[i-1] = str1
		case i%int2 == 0:
			out[i-1] = str2
		default:
			out[i-1] = strconv.Itoa(i)
		}
	}
	return out
}
