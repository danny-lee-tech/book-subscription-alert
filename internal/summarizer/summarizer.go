package summarizer

import (
	"context"
	"embed"
	"fmt"
	"math"
	"strings"
	"time"

	"google.golang.org/genai"
)

const RetryMaximum int = 6

//go:embed summarizer-instructions.txt
var promptTemplate embed.FS

func SummarizeText(client *genai.Client, company string, text string, url string) (string, error) {
	retryCount := 0
	var result *genai.GenerateContentResponse

	for {
		fmt.Println("Summarizing Text for " + company)
		// Read the content of the embedded file.
		prompt, err := promptTemplate.ReadFile("summarizer-instructions.txt")
		if err != nil {
			return "", err
		}
		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
		result, err = client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(fmt.Sprintf(string(prompt), company, url, text)),
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
			fmt.Printf("Service unavailable: Retry #%d for %d seconds\n", int(retryCount), int(sleepTime))
			time.Sleep(time.Duration(sleepTime) * time.Second)
		} else {
			break
		}
	}
	return result.Text(), nil
}
