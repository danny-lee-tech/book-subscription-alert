package owlcrate

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const ScrapeUrl string = "https://www.owlcrate.com/blogs/oc"

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
		chromedp.WaitVisible(".article h3 a", chromedp.ByQuery),
		chromedp.Text(".article h3 a", &title, chromedp.ByQuery),
		chromedp.AttributeValue(".article h3 a", "href", &href, &attributeFound, chromedp.ByQuery),
	)
	if err != nil {
		return "", "", err
	}

	if !strings.Contains(strings.ToLower(title), "limited editions") {
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
	var result string
	err := chromedp.Run(*c,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#bloggy--article", chromedp.ByQuery),
		chromedp.Text("#bloggy--article", &result, chromedp.ByQuery),
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
