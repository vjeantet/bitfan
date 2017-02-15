//go:generate bitfanDoc
package twitter

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt    *options
	stream *anaconda.Stream
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	Consumer_key       string
	Consumer_secret    string
	Oauth_token        string
	Oauth_token_secret string

	// Any keywords to track in the Twitter stream. For multiple keywords,
	// use the syntax ["foo", "bar"]. There’s a logical OR between each keyword
	// string listed and a logical AND between words separated by spaces per keyword string.
	// See https://dev.twitter.com/streaming/overview/request-parameters#track for more details.
	Keywords []string

	// A comma separated list of user IDs, indicating the users to return statuses for
	// in the Twitter stream.
	// See https://dev.twitter.com/streaming/overview/request-parameters#follow for more details.
	Follows []string

	// Record full tweet object as given to us by the Twitter Streaming API
	Full_tweet bool

	// Lets you ingore the retweets coming out of the Twitter API. Default false
	Ignore_retweets bool

	// A list of BCP 47 language identifiers corresponding to any of the languages
	// listed on Twitter’s advanced search page will only return tweets that have been
	// detected as being written in the specified languages
	Languages []string

	// A comma-separated list of longitude, latitude pairs specifying a set of bounding boxes
	// to filter tweets by. See
	// https://dev.twitter.com/streaming/overview/request-parameters#locations for more details
	Locations []string

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	Tags []string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt.Full_tweet = false
	p.opt.Ignore_retweets = false

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {
	anaconda.SetConsumerKey(p.opt.Consumer_key)
	anaconda.SetConsumerSecret(p.opt.Consumer_secret)
	api := anaconda.NewTwitterApi(p.opt.Oauth_token, p.opt.Oauth_token_secret)

	v := url.Values{}
	v.Set("track", strings.Join(p.opt.Keywords[:], " "))

	if len(p.opt.Languages) > 0 {
		v.Set("language", strings.Join(p.opt.Languages[:], ","))
	}

	if len(p.opt.Locations) > 0 {
		v.Set("locations", strings.Join(p.opt.Locations[:], ","))
	}

	if len(p.opt.Follows) > 0 {
		v.Set("follow", strings.Join(p.opt.Follows[:], ","))
	}

	p.stream = api.PublicStreamFilter(v)

	go p.doStream(p.stream, e, p.opt)

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.stream.Stop()
	return nil
}

func (p *processor) doStream(stream *anaconda.Stream, packet processors.IPacket, opt *options) {
	for streamObj := range stream.C {
		switch t := streamObj.(type) {
		case anaconda.Tweet:

			if opt.Ignore_retweets == true && t.Retweeted == true {
				continue
			}

			var r map[string]interface{}
			if opt.Full_tweet == true {
				msg, _ := json.Marshal(t)
				json.Unmarshal(msg, &r)
			} else {
				r = map[string]interface{}{}
				r["message"] = t.Text
				r["user"] = t.User.ScreenName
				r["user_name"] = t.User.Name
				r["client"] = t.Source
				r["retweeted"] = t.Retweeted
				r["source"] = fmt.Sprintf("http://twitter.com/%s/status/%d", t.User.ScreenName, t.Id)
				hashtags := []string{}
				for _, v := range t.Entities.Hashtags {
					hashtags = append(hashtags, v.Text)
				}
				r["hashtags"] = hashtags

				urls := []string{}
				for _, v := range t.Entities.Urls {
					urls = append(urls, v.Expanded_url)
				}
				r["urls"] = urls

				user_mentions := []map[string]string{}
				user_mention := map[string]string{}
				for _, v := range t.Entities.User_mentions {
					user_mention["screen_name"] = v.Screen_name
					user_mention["name"] = v.Name
					user_mentions = append(user_mentions, user_mention)
				}

				r["user_mentions"] = user_mentions
			}

			createdAtTime, err := t.CreatedAtTime()
			if err == nil {
				r["@timestamp"] = createdAtTime
			}

			e := p.NewPacket(t.Text, r)
			processors.AddFields(opt.Add_field, e.Fields())
			if len(opt.Tags) > 0 {
				processors.AddTags(opt.Tags, e.Fields())
			}
			p.Send(e, 0)
		default:
			// log.Warnf("unhandled type %T", t)
		}
	}
}
