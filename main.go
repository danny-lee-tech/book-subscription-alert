package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/danny-lee-tech/book-subscription-alert/internal/fairyloot"
	"github.com/danny-lee-tech/book-subscription-alert/internal/history"
	"github.com/danny-lee-tech/book-subscription-alert/internal/illumicrate"
	"github.com/danny-lee-tech/book-subscription-alert/internal/notifier"
	"github.com/danny-lee-tech/book-subscription-alert/internal/owlcrate"
	"github.com/danny-lee-tech/book-subscription-alert/internal/summarizer"
	"google.golang.org/genai"
)

func main() {
	geminiApiKey := os.Args[1]
	pushBulletApiKey := os.Args[2]

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(err)
	}

	var notif *notifier.Notifier
	if pushBulletApiKey != "" {
		notif = &notifier.Notifier{
			PushBullet: &notifier.PushBulleter{
				APIKey: pushBulletApiKey,
				Tag:    "book-subscription-alert",
				Title:  "New Book Subscription Alert",
			},
		}
	}

	err = checkOwlCrate(client, notif)
	if err != nil {
		log.Println("Error Checking OwlCrate:", err)
	}
	err = checkFairyLoot(client, notif)
	if err != nil {
		log.Println("Error Checking FairyLoot:", err)
	}
	err = checkIllumicrate(client, notif)
	if err != nil {
		log.Println("Error Checking FairyLoot:", err)
	}
}

func checkOwlCrate(client *genai.Client, notif *notifier.Notifier) error {
	post, postUrl, err := owlcrate.RetrieveLatestBlogPost()
	if err != nil {
		return err
	}

	owlCrateHistory := history.Init("owlcrate", 3)
	isDuplicate, err := owlCrateHistory.CheckIfExists(postUrl)
	if err != nil {
		return err
	}

	if isDuplicate {
		fmt.Println("Duplicate Owlcrate Post")
		return nil
	}

	summary, err := summarizer.SummarizeText(client, "OwlCrate", post, postUrl)
	if err != nil {
		return err
	}

	fmt.Println(summary)
	if notif != nil {
		notif.NotifyWithLink(summary, postUrl)
	}

	owlCrateHistory.RecordItemIfNotExist(postUrl)
	return nil
}

func checkFairyLoot(client *genai.Client, notif *notifier.Notifier) error {
	post, postUrl, err := fairyloot.RetrieveLatestBlogPost()
	if err != nil {
		return err
	}

	fairyLootHistory := history.Init("fairyloot", 3)
	isDuplicate, err := fairyLootHistory.CheckIfExists(postUrl)
	if err != nil {
		return err
	}

	if isDuplicate {
		fmt.Println("Duplicate FairyLoot Post")
		return nil
	}

	summary, err := summarizer.SummarizeText(client, "FairyLoot", post, postUrl)
	if err != nil {
		return err
	}

	fmt.Println(summary)
	if notif != nil {
		notif.NotifyWithLink(summary, postUrl)
	}

	fairyLootHistory.RecordItemIfNotExist(postUrl)
	return nil
}

func checkIllumicrate(client *genai.Client, notif *notifier.Notifier) error {
	post, postUrl, err := illumicrate.RetrieveLatestBlogPost()
	if err != nil {
		return err
	}

	illumicrateHistory := history.Init("illumicrate", 3)
	isDuplicate, err := illumicrateHistory.CheckIfExists(postUrl)
	if err != nil {
		return err
	}

	if isDuplicate {
		fmt.Println("Duplicate Illumicrate Post")
		return nil
	}

	summary, err := summarizer.SummarizeText(client, "Illumicrate", post, postUrl)
	if err != nil {
		return err
	}

	fmt.Println(summary)
	if notif != nil {
		notif.NotifyWithLink(summary, postUrl)
	}

	illumicrateHistory.RecordItemIfNotExist(postUrl)
	return nil
}
