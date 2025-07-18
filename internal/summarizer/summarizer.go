package summarizer

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"google.golang.org/genai"
)

const RetryMaximum int = 6

func SummarizeText(client *genai.Client, company string, text string, url string) (string, error) {
	retryCount := 0
	var result *genai.GenerateContentResponse
	var err error

	for true {
		fmt.Println("Summarizing Text for " + company)
		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
		result, err = client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text("Output a very short summary with just the logistics (which includes the author name, dates, and any pricing) for the content added at the end of this. Also, generate me a google calendar link (text only) that uses the /calendar/render API to create an event for the early access sale with a description that includes the url, "+url+", and the same short summary from earlier. The title of the event should contain '"+company+"', the book name(s), the author, and the words 'Early Access Sale'. The summary should focus on the US, not the UK. Here is the content: "+text),
			nil,
		)

		if err != nil {
			if !strings.Contains(err.Error(), "The model is overloaded") {
				return "", err
			}

			retryCount++
			if retryCount > RetryMaximum {
				return "", err
			}

			sleepTime := math.Pow(2, float64(retryCount-1)) * 10
			fmt.Println(fmt.Sprintf("Service unavailable: Retry #%d for %d seconds", int(retryCount), int(sleepTime)))
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}

		break
	}
	return result.Text(), nil
}
