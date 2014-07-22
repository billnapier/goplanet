package commands

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	viper.BindPFlag("rsstimeout", fetchCmd.Flags().Lookup("rsstimeout"))
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
	timeout := viper.GetInt("rsstimeout")
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
	fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
	for _, ch := range newchannels {
		chnl := chnlify(ch)
		if err := Channels().Insert(chnl); err != nil {
			if !strings.Contains(err.Error(), "E11000") {
				fmt.Printf("Database error. Err: %v", err)
			}
		}
	}
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	fmt.Printf("%d new item(s) in %s\n", len(newitems), feed.Url)
	for _, item := range newitems {
		itm := itmify(item, ch)
		if err := Items().Insert(itm); err != nil {
			if !strings.Contains(err.Error(), "E11000") {
				fmt.Printf("Database error. Err: %v", err)
			}
		}
	}
}

type Itm struct {
	Date         time.Time
	Key          string
	ChannelKey   string
	Title        string
	FullContent  string
	Links        []*rss.Link
	Description  string
	Author       rss.Author
	Categories   []*rss.Category
	Comments     string
	Enclosures   []*rss.Enclosure
	Guid         *string `bson:",omitempty"`
	Source       *rss.Source
	PubDate      string
	Id           string `bson:",omitempty"`
	Generator    *rss.Generator
	Contributors []string
	Content      *rss.Content
	Extensions   map[string]map[string][]rss.Extension
}

type Chnl struct {
	Key            string
	Title          string
	Links          []rss.Link
	Description    string
	Language       string
	Copyright      string
	ManagingEditor string
	WebMaster      string
	PubDate        string
	LastBuildDate  string
	Docs           string
	Categories     []*rss.Category
	Generator      rss.Generator
	TTL            int
	Rating         string
	SkipHours      []int
	SkipDays       []int
	Image          rss.Image
	ItemKeys       []string
	Cloud          rss.Cloud
	TextInput      rss.Input
	Extensions     map[string]map[string][]rss.Extension
	Id             string
	Rights         string
	Author         rss.Author
	SubTitle       rss.SubTitle
}

func itmify(o *rss.Item, ch *rss.Channel) Itm {
	var x Itm
	x.Title = o.Title
	x.Links = o.Links
	x.ChannelKey = ch.Key()
	x.Description = o.Description
	x.Author = o.Author
	x.Categories = o.Categories
	x.Comments = o.Comments
	x.Enclosures = o.Enclosures
	x.Guid = o.Guid
	x.PubDate = o.PubDate
	x.Id = o.Id
	x.Key = o.Key()
	x.Generator = o.Generator
	x.Contributors = o.Contributors
	x.Content = o.Content
	x.Extensions = o.Extensions
	x.Date, _ = o.ParsedPubDate()

	if o.Content != nil && o.Content.Text != "" {
		x.FullContent = o.Content.Text
	} else {
		x.FullContent = o.Description
	}

	return x
}

func chnlify(o *rss.Channel) Chnl {
	var x Chnl
	x.Key = o.Key()
	x.Title = o.Title
	x.Links = o.Links
	x.Description = o.Description
	x.Language = o.Language
	x.Copyright = o.Copyright
	x.ManagingEditor = o.ManagingEditor
	x.WebMaster = o.WebMaster
	x.PubDate = o.PubDate
	x.LastBuildDate = o.LastBuildDate
	x.Docs = o.Docs
	x.Categories = o.Categories
	x.Generator = o.Generator
	x.TTL = o.TTL
	x.Rating = o.Rating
	x.SkipHours = o.SkipHours
	x.SkipDays = o.SkipDays
	x.Image = o.Image
	x.Cloud = o.Cloud
	x.TextInput = o.TextInput
	x.Extensions = o.Extensions
	x.Id = o.Id
	x.Rights = o.Rights
	x.Author = o.Author
	x.SubTitle = o.SubTitle

	var keys []string
	for _, y := range o.Items {
		keys = append(keys, y.Key())
	}
	x.ItemKeys = keys

	return x
}
