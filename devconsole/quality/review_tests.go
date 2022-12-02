package quality

import (
	"github.com/accrescent/apkstat"
	"golang.org/x/exp/slices"
)

func RunReviewTests(apk *apk.APK) (issues []string) {
	permissions := apk.Manifest().UsesPermissions

	if permissions != nil {
		for _, permission := range *permissions {
			if slices.Contains(permissionReviewBlacklist, permission.Name) {
				issues = append(issues, permission.Name)
			}
		}
	}

	return issues
}
