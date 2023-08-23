package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func content(url string) (string, error) {
	if url != "https://www.google.com" {
		return "", fmt.Errorf("invalid url: %s", url)
	}

	return "google", nil
}

func processURLs(ctx context.Context, urls []string) error {
	g, ctx := errgroup.WithContext(ctx)

	results := make([]string, len(urls))

	for i, url := range urls {
		a, b := i, url

		g.Go(func() error {
			c, err := content(b)
			if err != nil {
				return err
			}

			results[a] = c
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	for _, c := range results {
		fmt.Println(c)
	}

	return nil
}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.google.com",
		"https://www.google.com",
		"https://www.google.com",
		// "https://www.gogle.com",
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()

	if err := processURLs(ctx, urls); err != nil {
		fmt.Println(err)
	}

	fmt.Println("done")
}
