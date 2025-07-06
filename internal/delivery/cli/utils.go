package cli

import "strconv"

func formatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}
