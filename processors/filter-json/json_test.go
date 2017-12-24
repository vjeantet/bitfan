package json

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vjeantet/bitfan/processors/testutils"
)

func TestCommonProcessorDetails(t *testing.T) {
	Convey("This processor... ", t, func() {
		p, _ := testutils.NewProcessor(New)

		Convey("is a valid bitfan processor", func() {
			_, ok := p.Processor.(*processor)
			So(ok, ShouldBeTrue)
		})

		Convey("does not have limit on concurent event processing", func() {
			So(p.MaxConcurent(), ShouldEqual, 0)
		})

		Convey("is self documented", func() {
			if p.Doc().Doc == "" {
				Println("Missing documentation for this processor")
			} else {
				So(true, ShouldBeTrue)
			}
		})
	})
}

func TestInvalidConfiguration(t *testing.T) {
	conf := map[string]interface{}{}

	Convey("When source is missing", t, func() {
		_, err := testutils.NewProcessor(New, conf)
		Convey("and error happen", func() {
			So(err, ShouldBeError)
		})
	})

}

func TestJsonFilter(t *testing.T) {

	Convey("Given an existing event source field contains JSON", t, func() {
		event := testutils.NewPacket("", map[string]interface{}{
			"thejson": `{ "hello": "world", "list": [ 1, 2, 3 ], "hash": { "k": "v" } }`,
		})
		conf := map[string]interface{}{
			"source": "thejson",
		}
		Convey("When the target is set to \"myjson\"", func() {

			conf["target"] = "myjson"

			Convey(`One event is produced`, func() {
				p, _ := testutils.NewProcessor(New, conf)
				p.Receive(event)

				So(p.SentPacketsCount(0), ShouldEqual, 1)

				Convey(`The "myjson.hello" field should contain "hello"`, func() {
					So(event.Fields().ValueOrEmptyForPathString("myjson.hello"), ShouldEqual, "world")
				})

				Convey("event tags field does not contains _jsonparsefailure", func() {
					tags, _ := event.Fields().ValueForPath("tags")
					So(tags, ShouldBeEmpty)
				})
			})

			Convey("When common options are set", func() {
				conf["add_field"] = map[string]interface{}{"addedfield": "addedfieldvalue"}
				p, _ := testutils.NewProcessor(New, conf)
				p.Receive(event)
				Convey("common options are applied", func() {
					So(event.Fields().ValueOrEmptyForPathString("addedfield"), ShouldEqual, "addedfieldvalue")
				})
			})

		})

		Convey("When the target is not set", func() {

			Convey("The parsed JSON is placed in the root", func() {
				p, _ := testutils.NewProcessor(New, conf)
				p.Receive(event)
				So(event.Fields().ValueOrEmptyForPathString("myjson.hello"), ShouldBeEmpty)
				So(event.Fields().ValueOrEmptyForPathString("hello"), ShouldEqual, "world")
			})

			Convey("event tags field does not contains _jsonparsefailure", func() {
				p, _ := testutils.NewProcessor(New, conf)
				p.Receive(event)
				tags, _ := event.Fields().ValueForPath("tags")
				So(tags, ShouldBeEmpty)
			})

			Convey("When common options are set", func() {
				Convey("common options are applied", func() {
					conf["add_field"] = map[string]interface{}{"addedfield": "addedfieldvalue"}
					p, _ := testutils.NewProcessor(New, conf)
					p.Receive(event)
					Convey("common options are applied", func() {
						So(event.Fields().ValueOrEmptyForPathString("addedfield"), ShouldEqual, "addedfieldvalue")
					})
				})
			})

		})

	})

}

func TestMissingSourceField(t *testing.T) {

	event := testutils.NewPacket("", map[string]interface{}{
		"thejson": `{ "hello": "world", "list": [ 1, 2, 3 ], "hash": { "k": "v" } }`,
	})
	conf := map[string]interface{}{
		"source": "unknow",
	}

	Convey("When SkipInvalidJson is not set", t, func() {
		p, _ := testutils.NewProcessor(New, conf)
		p.Receive(event)
		Convey("no event produced", func() {
			So(p.SentPacketsCount(0), ShouldEqual, 0)
		})
	})

	Convey("When SkipInvalidJson is true", t, func() {
		conf["skip_invalid_json"] = true
		p, _ := testutils.NewProcessor(New, conf)
		p.Receive(event)
		Convey("no event produced", func() {
			So(p.SentPacketsCount(0), ShouldEqual, 0)
		})
	})

}

func TestInvalidJsonData(t *testing.T) {
	event := testutils.NewPacket("", map[string]interface{}{
		"thejson": `, 3 ], "hash": { "k": "v" } }`,
	})
	conf := map[string]interface{}{
		"source": "thejson",
	}
	conf["add_field"] = map[string]interface{}{"addedfield": "addedfieldvalue"}
	Convey("When SkipInvalidJson is not set", t, func() {

		p, _ := testutils.NewProcessor(New, conf)
		p.Receive(event)
		Convey("one event produced", func() {
			So(p.SentPacketsCount(0), ShouldEqual, 1)
		})

		Convey("event tags field contains _jsonparsefailure", func() {
			tags, _ := event.Fields().ValueForPath("tags")
			So(tags, ShouldContain, "_jsonparsefailure")
		})

		Convey("common options are not applied", func() {
			So(event.Fields().Exists("addedfield"), ShouldBeFalse)
		})

	})

	Convey("When SkipInvalidJson is true", t, func() {
		conf["skip_on_invalid_json"] = true
		p, _ := testutils.NewProcessor(New, conf)
		p.Receive(event)

		Convey("no event produced", func() {
			So(p.SentPacketsCount(0), ShouldEqual, 0)
		})

	})

}
