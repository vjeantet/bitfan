package multilinecodec

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	javastacktrace = `Exception in thread "main" java.lang.NullPointerException
        at com.example.myproject.Book.getTitle(Book.java:16)
        at com.example.myproject.Author.getBookTitles(Author.java:25)
        at com.example.myproject.Bootstrap.main(Bootstrap.java:14)
Exception in thread "main" java.lang.NullPointerException
        at fr.example.myproject.Book.getTitle(Book.java:16)
        at fr.example.myproject.Author.getBookTitles(Author.java:25)
        at fr.example.myproject.Bootstrap.main(Bootstrap.java:14)
Exception in thread "main" java.lang.NullPointerException
        at de.example.myproject.Book.getTitle(Book.java:16)
        at de.example.myproject.Author.getBookTitles(Author.java:25)
        at de.example.myproject.Bootstrap.main(Bootstrap.java:14)
Exception in thread "main" java.lang.NullPointerException`

	nextcase = `HEADER 9200
LINE 1 2016-10-05 08:39:00 Some log data
LINE 2 2016-10-05 08:40:00 Some other log data
FOOTER
HEADER 9300
LINE 4 2016-11-05 08:39:00 Some log data in another log
LINE 5 2016-11-05 08:40:00 Some other log data in another log
FOOTER`

	customSeparator = `Nodality digital sprawl towards vehicle girl grenade. 
Tattoo systemic ablative face forwards girl Tokyo math-military-grade nodal point 

geodesic bomb.-ware marketing decay tattoo systema dead Chiba

spook wristwatch vinyl. cyber-tank-traps car jeans man render-farm knife. `

	codesource = `printf ("%10.10ld  \t %10.10ld \t %s\
%f", w, x, y, z );
printf ("Hello");
printf ("World");
printf ("%10.10ld  \t %10.10ld \t %s\
%f", w, x, y, z );
printf ("Not");`

	dataTimestamp = `[2015-08-24 11:49:14,389][INFO ][env                      ] [Letha] using [1] data paths, mounts [[/
(/dev/disk1)]], net usable_space [34.5gb], net total_space [118.9gb], types [hfs]
[2015-08-24 11:49:14,389][INFO ][env                      ] Some thing
[2015-08-24 11:49:14,389][INFO ][env                      ] [Letha] using [2] data paths, mounts [[/
(/dev/disk1)]], net usable_space [34.5gb], net total_space [118.9gb], types [hfs]`
)

func TestBuffer(t *testing.T) {
	d := NewDecoder(strings.NewReader(""))
	d.memory = "Hello\n"
	assert.Equal(t, []byte("Hello\n"), d.Buffer())
}

func TestSetOptionsError(t *testing.T) {
	d := NewDecoder(strings.NewReader("data"))
	conf := map[string]interface{}{
		"delimiter": 4,
	}
	err := d.SetOptions(conf, logrus.New(), "")
	assert.Error(t, err)
}

func TestDefaultSettings(t *testing.T) {

	expectData := []string{
		`Exception in thread "main" java.lang.NullPointerException
        at com.example.myproject.Book.getTitle(Book.java:16)
        at com.example.myproject.Author.getBookTitles(Author.java:25)
        at com.example.myproject.Bootstrap.main(Bootstrap.java:14)`,
		`Exception in thread "main" java.lang.NullPointerException
        at fr.example.myproject.Book.getTitle(Book.java:16)
        at fr.example.myproject.Author.getBookTitles(Author.java:25)
        at fr.example.myproject.Bootstrap.main(Bootstrap.java:14)`,
		`Exception in thread "main" java.lang.NullPointerException
        at de.example.myproject.Book.getTitle(Book.java:16)
        at de.example.myproject.Author.getBookTitles(Author.java:25)
        at de.example.myproject.Bootstrap.main(Bootstrap.java:14)`,
	}

	d := NewDecoder(strings.NewReader(javastacktrace))
	conf := map[string]interface{}{
	//"pattern": `^\s`,
	// "what":    "previous",
	}

	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)

	var m interface{}

	for _, v := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, v, m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestSettingsNextNegate(t *testing.T) {

	expectData := []string{
		"HEADER 9200\nLINE 1 2016-10-05 08:39:00 Some log data\nLINE 2 2016-10-05 08:40:00 Some other log data\nFOOTER",
		"HEADER 9300\nLINE 4 2016-11-05 08:39:00 Some log data in another log\nLINE 5 2016-11-05 08:40:00 Some other log data in another log\nFOOTER",
	}

	d := NewDecoder(strings.NewReader(nextcase))
	conf := map[string]interface{}{
		"pattern": `^FOOTER$`,
		"what":    "next",
		"negate":  true,
	}

	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)

	var m interface{}

	for _, v := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, v, m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestSettingsPreviousDelimiter(t *testing.T) {

	expectData := []string{
		"Nodality digital sprawl towards vehicle girl grenade",
		" \nTattoo systemic ablative face forwards girl Tokyo math-military-grade nodal point \n\ngeodesic bomb",
		"-ware marketing decay tattoo systema dead Chiba\n\nspook wristwatch vinyl",
		" cyber-tank-traps car jeans man render-farm knife",
	}

	d := NewDecoder(strings.NewReader(customSeparator))
	conf := map[string]interface{}{
		"pattern":   `*`,
		"Delimiter": `.`,
		"what":      "previous",
	}

	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)

	var m interface{}

	for _, v := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, v, m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestSettingsNext(t *testing.T) {

	expectData := []string{
		"printf (\"%10.10ld  \\t %10.10ld \\t %s\\\n%f\", w, x, y, z );",
		`printf ("Hello");`,
		`printf ("World");`,
		"printf (\"%10.10ld  \\t %10.10ld \\t %s\\\n%f\", w, x, y, z );",
		`printf ("Not");`,
	}

	d := NewDecoder(strings.NewReader(codesource))
	conf := map[string]interface{}{
		"pattern": `\\$`,
		"what":    "next",
	}

	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)

	var m interface{}

	for _, v := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, v, m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestSettingsWhatPreviousNegate(t *testing.T) {

	expectData := []string{
		"[2015-08-24 11:49:14,389][INFO ][env                      ] [Letha] using [1] data paths, mounts [[/\n(/dev/disk1)]], net usable_space [34.5gb], net total_space [118.9gb], types [hfs]",
		"[2015-08-24 11:49:14,389][INFO ][env                      ] Some thing",
	}

	d := NewDecoder(strings.NewReader(dataTimestamp))
	conf := map[string]interface{}{
		"pattern": `^\[[^\]]+\]`,
		"what":    "previous",
		"negate":  true,
	}

	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)
	var m interface{}

	for _, v := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, v, m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestMore(t *testing.T) {

	expectData := []string{
		`Exception in thread "main" java.lang.NullPointerException
        at com.example.myproject.Book.getTitle(Book.java:16)
        at com.example.myproject.Author.getBookTitles(Author.java:25)
        at com.example.myproject.Bootstrap.main(Bootstrap.java:14)`,
		`Exception in thread "main" java.lang.NullPointerException
        at fr.example.myproject.Book.getTitle(Book.java:16)
        at fr.example.myproject.Author.getBookTitles(Author.java:25)
        at fr.example.myproject.Bootstrap.main(Bootstrap.java:14)`,
		`Exception in thread "main" java.lang.NullPointerException
        at de.example.myproject.Book.getTitle(Book.java:16)
        at de.example.myproject.Author.getBookTitles(Author.java:25)
        at de.example.myproject.Bootstrap.main(Bootstrap.java:14)`,
	}

	d := NewDecoder(strings.NewReader(javastacktrace))
	conf := map[string]interface{}{
	//"pattern": `^\s`,
	// "what":    "previous",
	}

	err := d.SetOptions(conf, logrus.New(), "")
	assert.NoError(t, err)

	var m interface{}
	var i = 0
	for d.More() {
		err := d.Decode(&m)
		if i+1 <= len(expectData) {
			assert.NoError(t, err)
			assert.Equal(t, expectData[i], m)
			i = i + 1
		} else {
			assert.Error(t, err)
		}

	}
	assert.Equal(t, 3, i)
}
