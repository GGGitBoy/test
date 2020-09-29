package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	err := ChromedpPrintPdf("https://docs.rancher.cn/docs/rancher2/releases/v2.4.8/", "file.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ChromedpPrintPdf(url string, to string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	var html string
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		// chromedp.Click(`.pagination-nav__item--nex t > a`, chromedp.NodeVisible),
		chromedp.OuterHTML(".pagination-nav__item--next > a", &html, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().WithDisplayHeaderFooter(true).WithHeaderTemplate(html).
				Do(ctx)
			return err
		}),
	})

	fmt.Printf("%v", html)

	if err != nil {
		return fmt.Errorf("chromedp Run failed,err:%+v", err)
	}

	if err := ioutil.WriteFile(to, buf, 0644); err != nil {
		return fmt.Errorf("write to file failed,err:%+v", err)
	}

	return nil
}
