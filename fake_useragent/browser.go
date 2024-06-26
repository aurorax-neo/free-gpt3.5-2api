package browser

import (
	"free-gpt3.5-2api/fake_useragent/setting"
	"free-gpt3.5-2api/fake_useragent/useragent"
)

func Random() string {
	return defaultBrowser.Random()
}

func Chrome() string {
	return defaultBrowser.Chrome()
}

func InternetExplorer() string {
	return defaultBrowser.InternetExplorer()
}

func Firefox() string {
	return defaultBrowser.Firefox()
}

func Safari() string {
	return defaultBrowser.Safari()
}

func Android() string {
	return defaultBrowser.Android()
}

func MacOSX() string {
	return defaultBrowser.MacOSX()
}

func IOS() string {
	return defaultBrowser.IOS()
}

func Linux() string {
	return defaultBrowser.Linux()
}

func IPhone() string {
	return defaultBrowser.IPhone()
}

func IPad() string {
	return defaultBrowser.IPad()
}

func Computer() string {
	return defaultBrowser.Computer()
}

func Mobile() string {
	return defaultBrowser.Mobile()
}

func (b *Browser) Random() string {
	return useragent.UA.GetAllRandom()
}

func (b *Browser) Chrome() string {
	return useragent.UA.GetRandom(setting.CHROME)
}

func (b *Browser) InternetExplorer() string {
	return useragent.UA.GetRandom(setting.INTERNET_EXPLORER)
}

func (b *Browser) Firefox() string {
	return useragent.UA.GetRandom(setting.FIREFOX)
}

func (b *Browser) Safari() string {
	return useragent.UA.GetRandom(setting.SAFARI)
}

func (b *Browser) Android() string {
	return useragent.UA.GetRandom(setting.ANDROID)
}

func (b *Browser) MacOSX() string {
	return useragent.UA.GetRandom(setting.MAC_OS_X)
}

func (b *Browser) IOS() string {
	return useragent.UA.GetRandom(setting.IOS)
}

func (b *Browser) Linux() string {
	return useragent.UA.GetRandom(setting.LINUX)
}

func (b *Browser) IPhone() string {
	return useragent.UA.GetRandom(setting.IPHONE)
}

func (b *Browser) IPad() string {
	return useragent.UA.GetRandom(setting.IPAD)
}

func (b *Browser) Computer() string {
	return useragent.UA.GetRandom(setting.COMPUTER)
}

func (b *Browser) Mobile() string {
	return useragent.UA.GetRandom(setting.MOBILE)
}
