package spiders

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"

	"free-gpt3.5-2api/fake_useragent/downloader"
	"free-gpt3.5-2api/fake_useragent/scheduler"
	"free-gpt3.5-2api/fake_useragent/setting"
	"free-gpt3.5-2api/fake_useragent/useragent"
)

type Spider struct {
	Attribute
	FullUrl string
}

type Attribute struct {
	Tag      string
	Category string
	Page     int
}

var urlAttributeResults = make(map[string]Attribute)

func NewBrowserSpider() *Spider {
	return &Spider{}
}

func (a *Attribute) GetSpider() *Spider {
	return &Spider{
		Attribute: Attribute{
			Tag:      a.Tag,
			Category: a.Category,
			Page:     a.Page,
		},
		FullUrl: fmt.Sprintf(setting.BROWSER_URL, a.Tag, a.Category, a.Page),
	}
}

func (s *Spider) AppendBrowser(maxPage int) {
	for tag, categorys := range setting.BrowserUserAgentMaps {
		for _, category := range categorys {
			for page := 1; page <= maxPage; page++ {
				attribute := Attribute{Tag: tag, Category: category, Page: page}
				urlAttributeResults[attribute.GetSpider().FullUrl] = attribute
				scheduler.AppendUrl(attribute.GetSpider().FullUrl)
			}
		}
	}
}

// TODO: use and limit goroutine.
func (s *Spider) StartBrowser(delay time.Duration, timeout time.Duration) {
	count := scheduler.CountUrl()
	for i := 0; i <= count; i++ {
		if url := scheduler.PopUrl(); url != "" {
			downloader := downloader.Download{Delay: delay, Timeout: timeout}
			resp, err := downloader.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return
			}

			doc.Find("td.useragent a").Each(func(i int, selection *goquery.Selection) {
				if value := selection.Text(); value != "" {
					useragent.UA.Set(urlAttributeResults[url].Category, value)
				}
			})
		}
	}
}
