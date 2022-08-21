package quality

import "github.com/accrescent/apkstat"

func RunReviewTests(apk *apk.APK) []string {
	// We don't want this slice to ever be nil because returning it as JSON would result in a
	// null value in JavaScript which is an unnecessary extra case we would need to handle since
	// null isn't iterable. See https://github.com/gin-gonic/gin/issues/125.
	errors := []string{}

	permissions := apk.Manifest().UsesPermissions

	if permissions != nil {
		for _, permission := range *permissions {
			switch permission.Name {
			case "android.permission.ACCESS_BACKGROUND_LOCATION":
				fallthrough
			case "android.permission.ACCESS_COARSE_LOCATION":
				fallthrough
			case "android.permission.ACCESS_FINE_LOCATION":
				fallthrough
			case "android.permission.BLUETOOTH_SCAN":
				fallthrough
			case "android.permission.MANAGE_EXTERNAL_STORAGE":
				fallthrough
			case "android.permission.NEARBY_WIFI_DEVICES":
				fallthrough
			case "android.permission.PROCESS_OUTGOING_CALLS":
				fallthrough
			case "android.permission.QUERY_ALL_PACKAGES":
				fallthrough
			case "android.permission.READ_CALL_LOG":
				fallthrough
			case "android.permission.READ_EXTERNAL_STORAGE":
				fallthrough
			case "android.permission.READ_PHONE_STATE":
				fallthrough
			case "android.permission.READ_MEDIA_AUDIO":
				fallthrough
			case "android.permission.READ_MEDIA_IMAGES":
				fallthrough
			case "android.permission.READ_MEDIA_VIDEO":
				fallthrough
			case "android.permission.READ_SMS":
				fallthrough
			case "android.permission.RECEIVE_MMS":
				fallthrough
			case "android.permission.RECEIVE_SMS":
				fallthrough
			case "android.permission.RECEIVE_WAP_PUSH":
				fallthrough
			case "android.permission.REQUEST_INSTALL_PACKAGES":
				fallthrough
			case "android.permission.SEND_SMS":
				fallthrough
			case "android.permission.WRITE_CALL_LOG":
				errors = append(errors, permission.Name)
			}
		}
	}

	return errors
}
