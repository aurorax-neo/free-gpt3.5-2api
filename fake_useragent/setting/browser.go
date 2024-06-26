package setting

const (
	CHROME            = "chrome"
	INTERNET_EXPLORER = "internet-explorer"
	FIREFOX           = "firefox"
	SAFARI            = "safari"

	ANDROID  = "android"
	MAC_OS_X = "mac-os-x"
	IOS      = "ios"
	LINUX    = "linux"

	IPHONE = "iphone"
	IPAD   = "ipad"

	COMPUTER = "computer"
	MOBILE   = "mobile"
)

var (
	BrowserUserAgentMaps = map[string][]string{
		"software_name": {
			CHROME,
			INTERNET_EXPLORER,
			FIREFOX,
			SAFARI,
		},
		"operating_system_name": {
			ANDROID,
			MAC_OS_X,
			IOS,
			LINUX,
		},
		"operating_platform": {
			IPHONE,
			IPAD,
		},
		"hardware_type_specific": {
			COMPUTER,
			MOBILE,
		},
	}
)
