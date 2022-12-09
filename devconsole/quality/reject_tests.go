package quality

import (
	"errors"
	"fmt"

	"github.com/accrescent/apkstat"
	"golang.org/x/mod/semver"

	pb "github.com/accrescent/devconsole/pb"
)

func RunRejectTests(metadata *pb.BuildApksResult, apk *apk.APK, uploadType UploadType) error {
	// Bundletool version used to generate APK set
	bundletoolVersion := metadata.GetBundletool().GetVersion()
	if semver.Compare("v"+bundletoolVersion, "v"+MIN_BUNDLETOOL_VERSION) == -1 {
		return fmt.Errorf(
			"APK set generated with bundletool %s but mininum supported version is %s",
			bundletoolVersion,
			MIN_BUNDLETOOL_VERSION,
		)
	}

	manifest := apk.Manifest()
	targetSDK := manifest.UsesSDK.TargetSDKVersion

	// Target SDK
	switch {
	case targetSDK == nil:
		return errors.New("Required field 'targetSdk' not found")
	case uploadType == NewApp && *targetSDK < MIN_TARGET_SDK_NEW_APP:
		return fmt.Errorf(
			"App target SDK is %d but the minimum is %d",
			*targetSDK, MIN_TARGET_SDK_NEW_APP,
		)
	case uploadType == Update && *targetSDK < MIN_TARGET_SDK_UPDATE:
		return fmt.Errorf(
			"App target SDK is %d but the minimum is %d",
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
