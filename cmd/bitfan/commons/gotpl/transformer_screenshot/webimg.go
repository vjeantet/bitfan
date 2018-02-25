package main

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Machiel/slugify"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

// var cm *chromedp.CDP
// var cs.ctx context.Context

type ChromeShot struct {
	cdp    *chromedp.CDP
	ctx    context.Context
	cancel context.CancelFunc
}

func NewChromeShot() (*ChromeShot, error) {

	// var
	// create context
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// create chrome instance
	// cm, err = chromedp.New(cs.ctx, chromedp.WithLogf(l.Printf))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cdp, err := chromedp.New(ctx, chromedp.WithRunnerOptions(
		// runner.Flag("headless", false),
		runner.Flag("disable-gpu", true),
		runner.Flag("no-first-run", true),
		runner.Flag("no-default-browser-check", true),
	), chromedp.WithLog(log.New(os.Stderr, "", 0).Printf))
	if err != nil {
		return nil, err
	}
	return &ChromeShot{
		cdp:    cdp,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (c *ChromeShot) urlForHtmlContent(content string) (string, error) {
	dir, _ := ioutil.TempDir("", "example")
	tmpFile := filepath.Join(dir, "index.html")
	err := ioutil.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		return "", err
	}
	return `file:///` + tmpFile, nil
}

func (cs *ChromeShot) EmbedImageForDomElements(htmlContent string, ids []string) (string, error) {
	urlStr, err := cs.urlForHtmlContent(htmlContent)
	if err != nil {
		return "", err
	}
	l.Println("urlStr-->", urlStr)

	cs.cdp.Run(cs.ctx, chromedp.Navigate(urlStr))
	// cs.cdp.Run(cs.ctx, chromedp.Sleep(2*time.Second))

	// cs.cdp.Run(cs.ctx, chromedp.Sleep(2*time.Second))
	datas := map[string][]byte{}
	l.Println("ids-->", ids)
	for k, id := range ids {
		l.Println("k,id-->", k, id)
		var buf []byte

		// cctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		if err := cs.cdp.Run(cs.ctx, chromedp.WaitVisible(id, chromedp.ByID)); err != nil {
			l.Println("err-->", err)
			continue
		}

		//cs.cdp.Run(context.Background(), chromedp.Screenshot(id, &buf, chromedp.NodeVisible, chromedp.ByID))
		cs.cdp.Run(cs.ctx, chromedp.CaptureScreenshot(&buf))
		l.Println(buf)
		datas[id] = buf
	}
	l.Println("--> END")

	// get loaded page source
	var source string
	if err := cs.cdp.Run(cs.ctx, chromedp.OuterHTML(`html`, &source)); err != nil {
		l.Println("err-->", err)
		return "", err
	}

	// Replace elements with images
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		return "", err
	}
	for i, v := range datas {
		// Write To Disk
		ioutil.WriteFile(slugify.Slugify(i)+".png", v, 0644)

		// Replace HTML
		buf64 := base64.StdEncoding.EncodeToString(datas[i])
		sel := doc.Find(i)
		for k := range sel.Nodes {
			single := sel.Eq(k)
			single.ReplaceWithHtml(`<img src="data:image/png;base64,` + buf64 + `" />`)
		}
	}

	return doc.Html()
}

func (cs *ChromeShot) Stop() error {
	// cs.cancel()
	// shutdown chrome
	err := cs.cdp.Shutdown(cs.ctx)
	if err != nil {
		return err
	}

	// wait for chrome to finish
	return cs.cdp.Wait()
}
