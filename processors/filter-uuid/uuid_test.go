package uuid

import (
	"testing"

	"bitfan/processors/testutils"
	. "github.com/smartystreets/goconvey/convey"
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

	Convey("When target is missing", t, func() {
		_, err := testutils.NewProcessor(New, conf)
		Convey("and error happen", func() {
			So(err, ShouldBeError)
		})
	})
}

func TestUuidFilterGeneration(t *testing.T) {
	Convey("When the target field name does not exist", t, func() {
		event := testutils.NewPacketOld("", map[string]interface{}{})
		conf := map[string]interface{}{
			"target": "name1",
		}
		p, _ := testutils.NewProcessor(New, conf)
		p.Receive(event)

		Convey("one event is produced", func() {
			So(p.SentPacketsCount(0), ShouldEqual, 1)
		})
		Convey("event contains a new field name1", func() {
			em := p.SentPackets(0)[0]
			So(em.Fields().Exists("name1"), ShouldBeTrue)
		})
	})

	Convey("When the target field exists", t, func() {
		event := testutils.NewPacketOld("", map[string]interface{}{
			"name1": "test",
		})

		Convey("When overwrite is true", func() {
			conf := map[string]interface{}{
				"target":    "name1",
				"overwrite": true,
			}
			p, _ := testutils.NewProcessor(New, conf)
			p.Receive(event)

			Convey("event field name1 is not modifier", func() {
				em := p.SentPackets(0)[0]
				So(em.Fields().ValueOrEmptyForPathString("name1"), ShouldNotEqual, "test")
			})

		})

		Convey("When overwrite is false", func() {
			conf := map[string]interface{}{
				"target":    "name1",
				"overwrite": false,
			}
			p, _ := testutils.NewProcessor(New, conf)
			p.Receive(event)

			Convey("event field name1 contains a new value", func() {
				em := p.SentPackets(0)[0]
				So(em.Fields().ValueOrEmptyForPathString("name1"), ShouldEqual, "test")
			})
		})

	})
}
