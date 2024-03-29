package quality

var permissionReviewBlacklist = []string{
	"android.permission.ACCESS_BACKGROUND_LOCATION",
	"android.permission.ACCESS_COARSE_LOCATION",
	"android.permission.ACCESS_FINE_LOCATION",
	"android.permission.BLUETOOTH_SCAN",
	"android.permission.CAMERA",
	"android.permission.MANAGE_EXTERNAL_STORAGE",
	"android.permission.NEARBY_WIFI_DEVICES",
	"android.permission.PROCESS_OUTGOING_CALLS",
	"android.permission.QUERY_ALL_PACKAGES",
	"android.permission.READ_CALL_LOG",
	"android.permission.READ_CONTACTS",
	"android.permission.READ_EXTERNAL_STORAGE",
	"android.permission.READ_MEDIA_AUDIO",
	"android.permission.READ_MEDIA_IMAGES",
	"android.permission.READ_MEDIA_VIDEO",
	"android.permission.READ_PHONE_NUMBERS",
	"android.permission.READ_PHONE_STATE",
	"android.permission.READ_SMS",
	"android.permission.RECEIVE_MMS",
	"android.permission.RECEIVE_SMS",
	"android.permission.RECEIVE_WAP_PUSH",
	"android.permission.RECORD_AUDIO",
	"android.permission.REQUEST_INSTALL_PACKAGES",
	"android.permission.SEND_SMS",
	"android.permission.WRITE_CALL_LOG",
	"android.permission.WRITE_CONTACTS",
	"android.permission.SYSTEM_ALERT_WINDOW",
}

var serviceIntentFilterActions = []string{
	"android.accessibilityservice.AccessibilityService",
	"android.net.VpnService",
	"android.view.InputMethod",
}
