package fairyloot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const ScrapeUrl string = "https://community.fairyloot.com/category/book-announcements/"
const RecentPostsLinkCssSelector string = ".blogArchive-posts h3.uppercase a"

func RetrieveLatestBlogPost() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, _ := chromedp.NewContext(ctx)
	defer func() {
		if err := chromedp.Cancel(c); err != nil {
			panic("chromedp could not be cancelled")
		}
	}()

	fmt.Println("Scraping: " + ScrapeUrl)

	var href string
	var attributeFound bool
	err := chromedp.Run(c,
		chromedp.Navigate(ScrapeUrl),
		chromedp.WaitEnabled(RecentPostsLinkCssSelector, chromedp.ByQuery),
		chromedp.AttributeValue(RecentPostsLinkCssSelector, "href", &href, &attributeFound, chromedp.ByQuery),
	)
	if err != nil {
		return "", "", err
	}

	fmt.Println("Scraping: " + href)
	postContent, err := scrapeBlogPost(href, &c)
	if err != nil {
		return "", "", err
	}

	return postContent, href, nil
}

func scrapeBlogPost(url string, c *context.Context) (string, error) {
	var result string
	err := chromedp.Run(*c,
		chromedp.Navigate(url),
		chromedp.WaitVisible(".singleBlog-content .wysiwyg", chromedp.ByQuery),
		chromedp.Text(".singleBlog-content .wysiwyg", &result, chromedp.ByQuery),
	)
	if err != nil {
		return "", err
	}

	return minimizeText(result), nil
}

func minimizeText(text string) string {
	text = strings.ReplaceAll(text, "\n\n", "\n")
	text = strings.ReplaceAll(text, "\n\n", "\n")
	text = strings.ReplaceAll(text, "\n\n", "\n")
	return text
}
