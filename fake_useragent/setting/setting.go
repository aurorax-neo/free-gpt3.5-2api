package setting

import "time"

const (
	VERSION = "0.2.0"

	BROWSER_URL            = "https://developers.whatismybrowser.com/useragents/explore/%s/%s/%d"
	BROWSER_MAX_PAGE       = 5
	BROWSER_ALLOW_MAX_PAGE = 8

	CACHE_VERSION = "0.2.0"
	CACHE_URL     = "https://raw.githubusercontent.com/EDDYCJY/fake-useragent/v0.2.0/static/"

	HTTP_TIMEOUT         = 15 * time.Second
	HTTP_DELAY           = 100 * time.Millisecond
	HTTP_ALLOW_MIN_DELAY = 100 * time.Millisecond

	TEMP_FILE_NAME      = "fake_useragent_%s.json"
	TEMP_FILE_TEST_NAME = "fake_useragent_test_%s.json"
)

func GetMaxPage(maxPage int) int {
	if maxPage > BROWSER_ALLOW_MAX_PAGE || maxPage == 0 {
		maxPage = BROWSER_MAX_PAGE
	}

	return maxPage
}

func GetDelay(delay time.Duration) time.Duration {
	if delay < HTTP_ALLOW_MIN_DELAY {
		delay = HTTP_ALLOW_MIN_DELAY
	}

	return delay
}

func GetTimeout(timeout time.Duration) time.Duration {
	if timeout == 0 {
		timeout = HTTP_TIMEOUT
	}

	return timeout
}
