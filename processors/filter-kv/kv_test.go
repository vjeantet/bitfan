package kv

import (
	"fmt"
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
	conf := map[string]interface{}{
		"source": fmt.Errorf("junk"),
	}

	Convey("When source is missing", t, func() {
		_, err := testutils.NewProcessor(New, conf)
		Convey("and error happen", func() {
			So(err, ShouldBeError)
		})
	})

}

func TestInvalidSource(t *testing.T) {
	Convey("When source field is a slice of string", t, func() {
		Convey("Then each slice value is processed", func() {
			conf := map[string]interface{}{
				"source": "foo",
			}
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("blababla", map[string]interface{}{
				"foo": []string{
					"hello=world foo=bar",
					"hello2=world2 foo2=bar2",
					"hello3=world3 foo3=bar3",
				},
			})
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"hello":  "world",
				"foo":    "bar",
				"hello2": "world2",
				"foo2":   "bar2",
				"hello3": "world3",
				"foo3":   "bar3",
			}
			for expectedPath, expectedValue := range expected {
				So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
			}

		})
	})
	Convey("When source field is a number", t, func() {
		Convey("Then nothing happens to the event", func() {
			conf := map[string]interface{}{
				"source": "foo",
			}
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("blababla", map[string]interface{}{
				"foo": 43,
			})
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]
			So(em.Fields().Old(), ShouldHaveLength, 3)
		})
	})
}

func TestAllowDuplicateValues(t *testing.T) {
	Convey("Removing duplicate key/value pairs", t, func() {

		Convey("Then allow_duplicate_values false", func() {
			conf := map[string]interface{}{
				"field_split":            "&",
				"source":                 "message",
				"allow_duplicate_values": false,
			}
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("foo=yeah&foo=yeah&foo=bar", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]
			//pp.Println("em-->", em.Fields().Old())
			expected := map[string]interface{}{
				"foo": []string{"yeah", "bar"},
			}
			for expectedPath, expectedValue := range expected {
				ev, _ := em.Fields().ValuesForPath(expectedPath)
				So(ev[0], ShouldResemble, expectedValue)
			}
		})
		Convey("Then allow_duplicate_values false 2", func() {
			conf := map[string]interface{}{
				"field_split":            "&",
				"source":                 "message",
				"allow_duplicate_values": false,
			}
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("foo=bar&foo=yeah&foo=yeah", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]
			//pp.Println("em-->", em.Fields().Old())
			expected := map[string]interface{}{
				"foo": []string{"bar", "yeah"},
			}
			for expectedPath, expectedValue := range expected {
				ev, _ := em.Fields().ValuesForPath(expectedPath)
				So(ev[0], ShouldResemble, expectedValue)
			}
		})

		Convey("Then allow_duplicate_values true", func() {
			conf := map[string]interface{}{
				"field_split":            "&",
				"source":                 "message",
				"allow_duplicate_values": true,
			}
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("foo=yeah&foo=yeah&foo=yeah&foo=bar", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]
			//pp.Println("em-->", em.Fields().Old())
			expected := map[string]interface{}{
				"foo": []string{"yeah", "yeah", "yeah", "bar"},
			}
			for expectedPath, expectedValue := range expected {
				ev, _ := em.Fields().ValuesForPath(expectedPath)
				So(ev[0], ShouldResemble, expectedValue)
			}
		})
	})
}

func TestDefaultKeys(t *testing.T) {
	Convey("When using default_keys", t, func() {
		conf := map[string]interface{}{
			"default_keys": map[string]interface{}{
				"foo": "xxx",
				"goo": "yyy",
			},
		}

		Convey("Then ...", func() {
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("hello=world foo=bar baz=fizz doublequoted=\"hello world\" singlequoted='hello world'", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"hello":        "world",
				"foo":          "bar",
				"goo":          "yyy",
				"baz":          "fizz",
				"doublequoted": "hello world",
				"singlequoted": "hello world",
			}
			for expectedPath, expectedValue := range expected {
				So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
			}
		})

		Convey("Then with a specific target", func() {
			conf["target"] = "kv"
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("hello=world foo=bar baz=fizz doublequoted=\"hello world\" singlequoted='hello world'", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"kv.hello":        "world",
				"kv.foo":          "bar",
				"kv.goo":          "yyy",
				"kv.baz":          "fizz",
				"kv.doublequoted": "hello world",
				"kv.singlequoted": "hello world",
			}
			for expectedPath, expectedValue := range expected {
				So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
			}
		})
	})
}

func TestTarget(t *testing.T) {
	Convey("When using target", t, func() {
		conf := map[string]interface{}{
			"target": "kv",
		}
		p, _ := testutils.NewProcessor(New, conf)

		Convey("Then ...", func() {
			event := testutils.NewPacketOld("hello=world foo=bar baz=fizz doublequoted=\"hello world\" singlequoted='hello world'", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"kv.hello":        "world",
				"kv.foo":          "bar",
				"kv.baz":          "fizz",
				"kv.doublequoted": "hello world",
				"kv.singlequoted": "hello world",
			}
			for expectedPath, expectedValue := range expected {
				So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
			}
		})
	})
}

func TestTrimKey(t *testing.T) {
	Convey("When using trim_key", t, func() {
		conf := map[string]interface{}{
			"field_split": "|",
			"value_split": "=",
			"trim_value":  " ",
			"trim_key":    " ",
		}
		p, _ := testutils.NewProcessor(New, conf)

		Convey("Then ...", func() {
			event := testutils.NewPacketOld("key1= value1 with spaces | key2 with spaces =value2", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"key1":             "value1 with spaces",
				"key2 with spaces": "value2",
			}
			for expectedPath, expectedValue := range expected {
				So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
			}
		})
	})
}

func TestExcludeKeys(t *testing.T) {
	Convey("When using exclude_keys", t, func() {
		conf := map[string]interface{}{
			"exclude_keys": []string{"foo", "singlequoted"},
		}
		p, _ := testutils.NewProcessor(New, conf)

		Convey("Then the specified keys are not valued into the event", func() {
			event := testutils.NewPacketOld("hello=world foo=bar baz=fizz doublequoted=\"hello world\" singlequoted='hello world'", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"hello":        "world",
				"foo":          nil,
				"baz":          "fizz",
				"doublequoted": "hello world",
				"singlequoted": nil,
			}
			for expectedPath, expectedValue := range expected {
				if expectedValue == nil {
					So(em.Fields().Exists(expectedPath), ShouldBeFalse)
				} else {
					So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
				}
			}
		})
	})
}

func TestIncludeKeys(t *testing.T) {
	Convey("When using include_keys", t, func() {
		conf := map[string]interface{}{
			"include_keys": []string{"foo", "singlequoted"},
		}
		p, _ := testutils.NewProcessor(New, conf)

		Convey("Then only then specified keys are valued into the event", func() {
			event := testutils.NewPacketOld("hello=world foo=bar baz=fizz doublequoted=\"hello world\" singlequoted='hello world'", nil)
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"hello":        nil,
				"foo":          "bar",
				"baz=":         nil,
				"doublequoted": nil,
				"singlequoted": "hello world",
				"brackets":     nil,
			}
			for expectedPath, expectedValue := range expected {
				if expectedValue == nil {
					So(em.Fields().Exists(expectedPath), ShouldBeFalse)
				} else {
					So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
				}
			}
		})
	})
}

func TestValueSplitUsingAlternateSplitter(t *testing.T) {
	Convey("Using a alternate value_split", t, func() {
		Convey("When value_split is :", func() {
			conf := map[string]interface{}{
				"value_split": ":",
			}
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("hello:=world foo:bar baz=:fizz doublequoted:\"hello world\" singlequoted:'hello world' brackets:(hello world)", nil)

			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"hello":        "=world",
				"foo":          "bar",
				"baz=":         "fizz",
				"doublequoted": "hello world",
				"singlequoted": "hello world",
				"brackets":     "hello world",
			}
			for expectedPath, expectedValue := range expected {
				So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
			}

		})
	})
}

func TestSpacesAttachedFields(t *testing.T) {
	Convey("When spaces are arround key pair value", t, func() {
		conf := map[string]interface{}{}
		p, _ := testutils.NewProcessor(New, conf)
		event := testutils.NewPacketOld("hello = world foo =bar baz= fizz doublequoted = \"hello world\" singlequoted= 'hello world' brackets =(hello world)", nil)
		p.Receive(event)

		So(p.SentPacketsCount(0), ShouldEqual, 1)
		em := p.SentPackets(0)[0]

		Convey("Then the produced event results in new fields/values for each keypair without the spaces", func() {
			expected := map[string]interface{}{
				"hello":        "world",
				"foo":          "bar",
				"baz":          "fizz",
				"doublequoted": "hello world",
				"singlequoted": "hello world",
				"brackets":     "hello world",
			}
			for expectedPath, expectedValue := range expected {
				So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
			}
		})
	})

	Convey("When there is escaped space in key or value", t, func() {
		conf := map[string]interface{}{
			"value_split": ":",
		}
		p, _ := testutils.NewProcessor(New, conf)
		event := testutils.NewPacketOld(`IKE:=Quick\ Mode\ completion IKE\ IDs:=subnet:\ x.x.x.x\ (mask=\ 255.255.255.254)\ and\ host:\ y.y.y.y`, nil)

		Convey("Then the produced event results in new fields/values for each keypair", func() {
			p.Receive(event)
			So(p.SentPacketsCount(0), ShouldEqual, 1)
			em := p.SentPackets(0)[0]

			expected := map[string]interface{}{
				"IKE":      `=Quick\ Mode\ completion`,
				`IKE\ IDs`: `=subnet:\ x.x.x.x\ (mask=\ 255.255.255.254)\ and\ host:\ y.y.y.y`,
			}
			for expectedPath, expectedValue := range expected {
				So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
			}
		})
	})

}

func TestDefaults(t *testing.T) {
	Convey("Given a processor with default configuration", t, func() {
		conf := map[string]interface{}{}
		Convey("When processor receive an event with a message containing key=value pairs", func() {
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("hello=world foo=bar baz=fizz doublequoted=\"hello world\" singlequoted='hello world' bracketsone=(hello world) bracketstwo=[hello world] bracketsthree=<hello world>", nil)
			p.Receive(event)

			Convey("Then the produced event results in new fields/values for each keypair", func() {
				So(p.SentPacketsCount(0), ShouldEqual, 1)
				em := p.SentPackets(0)[0]

				expected := map[string]interface{}{
					"hello":         "world",
					"foo":           "bar",
					"baz":           "fizz",
					"doublequoted":  "hello world",
					"singlequoted":  "hello world",
					"bracketsone":   "hello world",
					"bracketstwo":   "hello world",
					"bracketsthree": "hello world",
				}
				for expectedPath, expectedValue := range expected {
					So(em.Fields().ValueOrEmptyForPathString(expectedPath), ShouldEqual, expectedValue)
				}
			})

		})
		Convey("When processor receive an event with a invalid message", func() {
			p, _ := testutils.NewProcessor(New, conf)
			event := testutils.NewPacketOld("hree<hello world>", nil)
			eventCopy := event.Clone()
			p.Receive(event)

			Convey("one unmodified event is produced", func() {
				So(p.SentPacketsCount(0), ShouldEqual, 1)
				em := p.SentPackets(0)[0]
				So(eventCopy.Fields().Old(), ShouldHaveLength, len(em.Fields().Old()))
				So(em.Fields().ValueOrEmptyForPathString("message"), ShouldEqual, eventCopy.Fields().ValueOrEmptyForPathString("message"))
				So(em.Fields().ValueOrEmptyForPathString("@timestamp"), ShouldEqual, eventCopy.Fields().ValueOrEmptyForPathString("@timestamp"))
			})
		})
	})
}
