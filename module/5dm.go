/*
 *
 */
package module

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tagDong/mvcrawler"
	"github.com/tagDong/mvcrawler/dhttp"
	"github.com/tagDong/mvcrawler/util"
	"net/url"
	"strings"
)

type Dm5 struct {
	name    string
	baseUrl string
	logger  *util.Logger
}

func (this *Dm5) GetName() string {
	return this.name
}

func (this *Dm5) GetUrl() string {
	return this.baseUrl
}

func (this *Dm5) Search(txt string) []*mvcrawler.Message {
	ret := []*mvcrawler.Message{}
	data := url.Values{
		"s": {txt},
	}

	resp, err := dhttp.PostUrlencoded("https://www.5dm.tv", data, 0)
	if err != nil {
		this.logger.Errorln(err)
		return ret
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		this.logger.Errorln(err)
		return ret
	}
	_ = resp.Body.Close()

	doc.Find(".video-listing-content .video-item").Each(func(i int, selection *goquery.Selection) {
		var title, img, url string
		var ok bool
		title = selection.Find(".item-head h3 a").Text()

		if img, ok = selection.Find(".item-thumbnail a img").Attr("data-original"); !ok {
			return
		}
		if url, ok = selection.Find(".item-thumbnail a").Attr("href"); !ok {
			return
		}

		ret = append(ret, &mvcrawler.Message{
			Title: title,
			From:  this.GetName(),
			Img:   img,
			Url:   url,
		})
	})

	return ret
}

func (this *Dm5) Update() [][]*mvcrawler.Message {
	ret := [][]*mvcrawler.Message{}

	resp, err := dhttp.Get("https://www.5dm.tv/timeline", 0)
	if err != nil {
		this.logger.Errorln(err)
		return ret
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		this.logger.Errorln(err)
		return ret
	}
	_ = resp.Body.Close()

	doc.Find(".is-carousel").Each(func(i int, sele1 *goquery.Selection) {
		msgs := []*mvcrawler.Message{}
		sele1.Find(".video-item").Each(func(_ int, sele2 *goquery.Selection) {
			var title, img, url string
			var ok bool
			title = sele2.Find(".item-head h3 a").Text()

			if img, ok = sele2.Find(".item-thumbnail a img").Attr("data-original"); !ok {
				return
			}
			if url, ok = sele2.Find(".item-thumbnail a").Attr("href"); !ok {
				return
			}

			//更新状态
			var status string
			title, status = _5dmTitleStatus(title)

			msgs = append(msgs, &mvcrawler.Message{
				Title:  title,
				From:   this.GetName(),
				Img:    img,
				Url:    url,
				Status: status,
			})
		})
		ret = append(ret, msgs)
	})
	return ret
}

//拆分标题与更新状态
func _5dmTitleStatus(s string) (title, status string) {
	s1 := strings.Split(s, "【")
	if len(s1) != 2 {
		return s, "更新中..."
	}
	title = s1[0]

	s2 := strings.Split(s1[1], "】")
	if len(s2) != 2 {
		return title, "更新中..."
	}
	status = s2[0]
	return
}

func init() {
	mvcrawler.Register(mvcrawler.Dm5, func(l *util.Logger) mvcrawler.Module {
		return &Dm5{
			name:    "5dm",
			baseUrl: "https://www.5dm.tv",
			logger:  l,
		}
	})
}
