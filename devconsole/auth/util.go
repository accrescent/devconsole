package auth

import "crypto/subtle"

func ConstantTimeEqInt64(x, y int64) int {
	lower := subtle.ConstantTimeEq(int32(x), int32(y))
	upper := subtle.ConstantTimeEq(int32(x>>32), int32(y>>32))

	return lower & upper
}
