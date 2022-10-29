package quality

import (
	"github.com/accrescent/apkstat"
	"golang.org/x/exp/slices"
)

func RunReviewTests(apk *apk.APK) []string {
	// We don't want this slice to ever be nil because returning it as JSON would result in a
	// null value in JavaScript which is an unnecessary extra case we would need to handle since
	// null isn't iterable. See https://github.com/gin-gonic/gin/issues/125.
	errors := []string{}

	permissions := apk.Manifest().UsesPermissions

	if permissions != nil {
		for _, permission := range *permissions {
			if slices.Contains(permissionReviewBlacklist, permission.Name) {
				errors = append(errors, permission.Name)
			}
		}
	}

	return errors
}
