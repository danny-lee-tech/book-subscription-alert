package owlcrate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const ScrapeUrl string = "https://www.owlcrate.com/blogs/oc"
const RecentPostLinkCssSelector string = ".featured__blog__container a.article__link"
const RecentPostTitleCssSelector string = ".featured__blog__container .featured__blog__sub__heading h1"
const PostArticleContentCssSelector string = "#bloggy--article"

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

	var title string
	var href string
	var attributeFound bool
	err := chromedp.Run(c,
		chromedp.Navigate(ScrapeUrl),
		chromedp.WaitVisible(RecentPostLinkCssSelector, chromedp.ByQuery),
		chromedp.Text(RecentPostTitleCssSelector, &title, chromedp.ByQuery),
		chromedp.AttributeValue(RecentPostLinkCssSelector, "href", &href, &attributeFound, chromedp.ByQuery),
	)
	if err != nil {
		return "", "", err
	}

	if !strings.Contains(strings.ToLower(title), "limited edition") {
		return "", "", nil
	}

	href = "https://www.owlcrate.com" + href
	postContent, err := scrapeBlogPost(href, &c)
	if err != nil {
		return "", "", err
	}

	return postContent, href, nil
}

func scrapeBlogPost(url string, c *context.Context) (string, error) {
	fmt.Println("Scraping: " + url)
	var result string
	err := chromedp.Run(*c,
		chromedp.Navigate(url),
		chromedp.WaitVisible(PostArticleContentCssSelector, chromedp.ByQuery),
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
