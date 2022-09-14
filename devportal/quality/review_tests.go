package quality

import (
	"github.com/accrescent/apkstat"
	"golang.org/x/exp/slices"
)

var permissionReviewBlacklist = []string{
	"android.permission.ACCESS_BACKGROUND_LOCATION",
	"android.permission.ACCESS_COARSE_LOCATION",
	"android.permission.ACCESS_FINE_LOCATION",
	"android.permission.BLUETOOTH_SCAN",
	"android.permission.MANAGE_EXTERNAL_STORAGE",
	"android.permission.NEARBY_WIFI_DEVICES",
	"android.permission.PROCESS_OUTGOING_CALLS",
	"android.permission.QUERY_ALL_PACKAGES",
	"android.permission.READ_CALL_LOG",
	"android.permission.READ_EXTERNAL_STORAGE",
	"android.permission.READ_PHONE_STATE",
	"android.permission.READ_MEDIA_AUDIO",
	"android.permission.READ_MEDIA_IMAGES",
	"android.permission.READ_MEDIA_VIDEO",
	"android.permission.READ_SMS",
	"android.permission.RECEIVE_MMS",
	"android.permission.RECEIVE_SMS",
	"android.permission.RECEIVE_WAP_PUSH",
	"android.permission.REQUEST_INSTALL_PACKAGES",
	"android.permission.SEND_SMS",
	"android.permission.WRITE_CALL_LOG",
}

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
