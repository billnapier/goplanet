package commands

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	// Alias this import
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Feeds []string
	Port  int
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch feeds",
	Long:  "Goplanet will fetch all feeds listed in the config file",
	Run:   fetchRun,
}

func init() {
	fetchCmd.Flags().Int("rsstimeout", 600, "Timeout (in seconds) for RSS retrieval.")
	// Binds the flag value to the config value.  So value could also appear in config file.
	viper.BindPFlag("rsstimeout", fetchCmd.Flags().Lookup("rsstimout"))
}

func fetchRun(cmd *cobra.Command, args []string) {
	Fetcher()

	sigChan := make(chan os.Signal, 1)
	// When you receive an interupt signal, put it on this channel
	signal.Notify(sigChan, os.Interrupt)
	// This will block until something comes in the channel
	<-sigChan
}

func Fetcher() {
	var config Config

	if err := viper.Marshal(&config); err != nil {
		fmt.Println(err)
		return
	}

	for _, feed := range config.Feeds {
		go PollFeed(feed)
	}
}

func PollFeed(uri string) {
	timeout := viper.GetInt("rsstimout")
	if timeout < 60 {
		timeout = 60
	}
	feed := rss.New(timeout/60, true, chanHandler, itemHandler)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s\n", uri, err)
			return
		}

		fmt.Printf("Sleeping for %d seconds on %s\n", feed.SecondsTillUpdate(), uri)
		time.Sleep(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	fmt.Printf("%d new channels in %s\n", len(newchannels), feed.Url)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	fmt.Printf("%d new items in %s\n", len(newitems), feed.Url)
}
