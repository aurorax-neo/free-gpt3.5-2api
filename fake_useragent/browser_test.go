package browser

import (
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	if Random() == "" {
		t.Error("browser.Random is empty")
	}
}

func TestChrome(t *testing.T) {
	if Chrome() == "" {
		t.Error("browser.Chrome is empty")
	}
}

func TestInternetExplorer(t *testing.T) {
	if InternetExplorer() == "" {
		t.Error("browser.InternetExplorer is empty")
	}
}

func TestFirefox(t *testing.T) {
	if Firefox() == "" {
		t.Error("browser.Firefox is empty")
	}
}

func TestSafari(t *testing.T) {
	if Safari() == "" {
		t.Error("browser.Safari is empty")
	}
}

func TestAndroid(t *testing.T) {
	if Android() == "" {
		t.Error("browser.Android is empty")
	}
}

func TestMacOSX(t *testing.T) {
	if MacOSX() == "" {
		t.Error("browser.MacOSX is empty")
	}
}

func TestIOS(t *testing.T) {
	if IOS() == "" {
		t.Error("browser.IOS is empty")
	}
}

func TestLinux(t *testing.T) {
	if Linux() == "" {
		t.Error("browser.IOS is empty")
	}
}

func TestIPhone(t *testing.T) {
	if IPhone() == "" {
		t.Error("browser.IPhone is empty")
	}
}

func TestIPad(t *testing.T) {
	if IPad() == "" {
		t.Error("browser.IPad is empty")
	}
}

func TestComputer(t *testing.T) {
	if Computer() == "" {
		t.Error("browser.Computer is empty")
	}
}

func TestMobile(t *testing.T) {
	if Mobile() == "" {
		t.Error("browser.Mobile is empty")
	}
}

func TestBrowser_Random(t *testing.T) {
	b := NewBrowser(Client{
		MaxPage: 1,
		Delay:   200 * time.Millisecond,
		Timeout: 10 * time.Second,
	}, Cache{
		UpdateFile: true,
	})

	if b.Random() == "" {
		t.Error("NewBrowser.Random is empty")
	}
}

func TestBrowser_Chrome(t *testing.T) {
	b := NewBrowser(Client{
		MaxPage: 1,
		Delay:   250 * time.Millisecond,
		Timeout: 20 * time.Second,
	}, Cache{})

	if b.Chrome() == "" {
		t.Error("NewBrowser.Chrome is empty")
	}
}

func TestBrowser_IOS(t *testing.T) {
	b := NewBrowser(Client{
		MaxPage: 1,
		Delay:   350 * time.Millisecond,
		Timeout: 20 * time.Second,
	}, Cache{
		UpdateFile: true,
	})

	if b.IOS() == "" {
		t.Error("NewBrowser.IOS is empty")
	}
}
