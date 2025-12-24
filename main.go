package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/danny-lee-tech/book-subscription-alert/internal/config"
	"github.com/danny-lee-tech/book-subscription-alert/internal/fairyloot"
	"github.com/danny-lee-tech/book-subscription-alert/internal/history"
	"github.com/danny-lee-tech/book-subscription-alert/internal/illumicrate"
	"github.com/danny-lee-tech/book-subscription-alert/internal/notifier"
	"github.com/danny-lee-tech/book-subscription-alert/internal/owlcrate"
	"github.com/danny-lee-tech/book-subscription-alert/internal/summarizer"
	"google.golang.org/genai"
	"gopkg.in/yaml.v2"
)

var DefaultConfigLocation = "configs/config.yml"

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(err)
	}

	notif := &notifier.Notifier{
		PushBullet: &notifier.PushBulleter{
			APIKey: config.PushBulletApiKey,
			Tag:    "book-subscription-alert",
			Title:  "New Book Subscription Alert",
		},
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

	if postUrl == "" {
		fmt.Println("No recent special edition post found")
		return nil
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
		err = notif.NotifyWithLink(summary, postUrl)
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
		}
	}

	owlCrateHistory.RecordItemIfNotExist(postUrl)
	return nil
}

func checkFairyLoot(client *genai.Client, notif *notifier.Notifier) error {
	post, postUrl, err := fairyloot.RetrieveLatestBlogPost()
	if err != nil {
		return err
	}

	if postUrl == "" {
		fmt.Println("No recent special edition post found")
		return nil
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

	if postUrl == "" {
		fmt.Println("No recent special edition post found")
		return nil
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

func getConfig() (config.Config, error) {
	configLocation := getConfigLocation()
	cfgBytes, err := os.ReadFile(configLocation)
	if err != nil {
		return config.Config{}, err
	}

	var cfg config.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return config.Config{}, err
	}

	err = validateConfig(&cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func getConfigLocation() string {
	configLocation := os.Getenv("CONFIG_LOCATION")
	if configLocation != "" {
		return configLocation
	}
	return DefaultConfigLocation
}

func validateConfig(cfg *config.Config) error {
	if cfg.GeminiApiKey == "" {
		return errors.New("missing required field: gemini_api_key")
	}

	if cfg.PushBulletApiKey == "" {
		return errors.New("missing required field: pushbullet_api_key")
	}

	return nil
}
