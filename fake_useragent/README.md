# Fake Useragent ![image](https://api.travis-ci.org/EDDYCJY/fake-useragent.svg?branch=master)

A wide variety of random useragents

[简体中文](https://github.com/EDDYCJY/fake-useragent/blob/master/README_ZH.md)

## Support

- All User-Agent Random
- Chrome
- InternetExplorer (IE)
- Firefox
- Safari
- Android
- MacOSX
- IOS
- Linux
- IPhone
- IPad
- Computer
- Mobile

## Installation

```
$ go get github.com/EDDYCJY/fake-useragent
```

## Usage

``` go
package main

import (
	"log"

	"github.com/EDDYCJY/fake-useragent"
)

func main() {
	// recommend to use
	random := browser.Random()
	log.Printf("Random: %s", random)

	chrome := browser.Chrome()
	log.Printf("Chrome: %s", chrome)

	internetExplorer := browser.InternetExplorer()
	log.Printf("IE: %s", internetExplorer)

	firefox := browser.Firefox()
	log.Printf("Firefox: %s", firefox)

	safari := browser.Safari()
	log.Printf("Safari: %s", safari)

	android := browser.Android()
	log.Printf("Android: %s", android)

	macOSX := browser.MacOSX()
	log.Printf("MacOSX: %s", macOSX)

	ios := browser.IOS()
	log.Printf("IOS: %s", ios)

	linux := browser.Linux()
	log.Printf("Linux: %s", linux)

	iphone := browser.IPhone()
	log.Printf("IPhone: %s", iphone)

	ipad := browser.IPad()
	log.Printf("IPad: %s", ipad)

	computer := browser.Computer()
	log.Printf("Computer: %s", computer)

	mobile := browser.Mobile()
	log.Printf("Mobile: %s", mobile)
}
```

### Customize

You can adjust the maximum number of crawl pages and time intervals, maximum timeouts. If not, it is the default.

``` go
client := browser.Client{
	MaxPage: 3,
	Delay: 200 * time.Millisecond,
	Timeout: 10 * time.Second,
}
cache := browser.Cache{}
b := browser.NewBrowser(client, cache)

random := b.Random()
```

Update the browser temporary file cache (re-get the source to get the latest data).

``` go
client := browser.Client{}
cache := browser.Cache{
	UpdateFile: true,
}
b := browser.NewBrowser(client, cache)
```

Finally, I recommend the regular usage.

## Output

``` sh
Random: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36

Chrome: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36

IE: Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)

Firefox: Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0

Safari: Mozilla/5.0 (iPhone; CPU iPhone OS 11_2_5 like Mac OS X) AppleWebKit/604.5.6 (KHTML, like Gecko) Version/11.0 Mobile/15D60 Safari/604.1

Android: Mozilla/5.0 (Linux; Android 6.0; MYA-L22 Build/HUAWEIMYA-L22) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.84 Mobile Safari/537.36

MacOSX: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/602.2.14 (KHTML, like Gecko) Version/10.0.1 Safari/602.2.14

IOS: Mozilla/5.0 (iPhone; CPU iPhone OS 10_1 like Mac OS X) AppleWebKit/602.2.14 (KHTML, like Gecko) Version/10.0 Mobile/14B72 Safari/602.1

Linux: Mozilla/5.0 (X11; Linux x86_64; rv:42.0) Gecko/20100101 Firefox/42.0

IPhone: Mozilla/5.0 (iPhone; CPU iPhone OS 10_2 like Mac OS X) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0 Mobile/14C92 Safari/602.1

IPad: Mozilla/5.0 (iPad; CPU OS 5_0_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A405 Safari/7534.48.3

Computer: Mozilla/5.0 (Windows NT 10.0; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0

Mobile: Mozilla/5.0 (Linux; Android 7.0; Redmi Note 4 Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.111 Mobile Safari/537.36
```
