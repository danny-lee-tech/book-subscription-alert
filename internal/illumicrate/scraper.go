package illumicrate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const ScrapeUrl string = "https://us.illumicrate.com/blogs/news"
const RecentPostsLinkCssSelector string = ".article-card > a"
const PostArticleContentCssSelector string = ".article__content"

func RetrieveLatestBlogPost() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, _ := chromedp.NewContext(ctx)
	defer func() {
		if err := chromedp.Cancel(c); err != nil {
			panic("chromedp could not be cancelled")
		}
	}()

	var title string
	var href string
	var attributeFound bool
	err := chromedp.Run(c,
		chromedp.Navigate(ScrapeUrl),
		chromedp.WaitEnabled(RecentPostsLinkCssSelector, chromedp.ByQuery),
		chromedp.AttributeValue(RecentPostsLinkCssSelector, "aria-label", &title, &attributeFound, chromedp.ByQuery),
		chromedp.AttributeValue(RecentPostsLinkCssSelector, "href", &href, &attributeFound, chromedp.ByQuery),
	)
	if err != nil {
		return "", "", err
	}

	if !strings.Contains(strings.ToLower(title), "exclusive:") {
		return "", "", nil
	}

	href = "https://us.illumicrate.com" + href
	postContent, err := scrapeBlogPost(href, &c)
	if err != nil {
		return "", "", err
	}

	return postContent, href, nil
}

func scrapeBlogPost(url string, c *context.Context) (string, error) {
	fmt.Println("Scraping URL: " + url)
	var result string
	err := chromedp.Run(*c,
		chromedp.Navigate(url),
		chromedp.WaitEnabled(PostArticleContentCssSelector, chromedp.ByQuery),
		chromedp.Text(PostArticleContentCssSelector, &result, chromedp.ByQuery),
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
