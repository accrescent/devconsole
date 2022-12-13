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

	services := apk.Manifest().Application.Services
	if services != nil {
		for _, service := range *services {
			filters := service.IntentFilters
			if filters != nil {
				for _, filter := range *filters {
					for _, action := range filter.Actions {
						if slices.Contains(serviceIntentFilterActions, action.Name) &&
							!slices.Contains(issues, action.Name) {
							issues = append(issues, action.Name)
						}
					}
				}
			}
		}
	}

	return issues
}
