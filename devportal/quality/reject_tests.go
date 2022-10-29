package quality

import (
	"errors"
	"fmt"

	"github.com/accrescent/apkstat"
)

func RunRejectTests(apk *apk.APK, uploadType UploadType) error {
	manifest := apk.Manifest()
	targetSDK := manifest.UsesSDK.TargetSDKVersion

	// Target SDK
	switch {
	case targetSDK == nil:
		return errors.New("required field 'targetSdk' not found")
	case uploadType == NewApp && *targetSDK < MIN_TARGET_SDK_NEW_APP:
		return fmt.Errorf(
			"app target SDK is %d but the minimum is %d",
			*targetSDK, MIN_TARGET_SDK_NEW_APP,
		)
	case uploadType == Update && *targetSDK < MIN_TARGET_SDK_UPDATE:
		return fmt.Errorf(
			"app target SDK is %d but the minimum is %d",
			*targetSDK, MIN_TARGET_SDK_UPDATE,
		)
	}

	// android:debuggable
	if manifest.Application.Debuggable != nil && *manifest.Application.Debuggable {
		return errors.New("android:debuggable should not be set to true")
	}

	// android:usesCleartextTraffic
	usesCleartextTraffic := manifest.Application.UsesCleartextTraffic
	if usesCleartextTraffic != nil && *usesCleartextTraffic {
		return errors.New("android:usesCleartextTraffic should not be set to true")
	}

	return nil
}
