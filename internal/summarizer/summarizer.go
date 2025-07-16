package summarizer

import (
	"context"
	"time"

	"google.golang.org/genai"
)

func SummarizeText(client *genai.Client, company string, text string, url string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text("Output a very short summary with just the logistics (which includes the author name, dates, and any pricing) for the content added at the end of this. Also, generate me a google calendar link (text only) that uses the /calendar/render API to create an event for the early access sale with a description that includes the url, "+url+", and the same short summary from earlier. The title of the event should contain '"+company+"', the book name(s), the author, and the words 'Early Access Sale'. Here is the content: "+text),
		nil,
	)

	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
