package jodaTime

/*
jodaTime provides a date formatter using the yoda syntax.
http://joda-time.sourceforge.net/apidocs/org/joda/time/format/DateTimeFormat.html
*/

import (
	"fmt"
	"strconv"
	"time"
)

/*
 Symbol  Meaning                      Presentation  Examples
 ------  -------                      ------------  -------
 G       era                          text          AD
 C       century of era (>=0)         number        20
 Y       year of era (>=0)            year          1996

 x       weekyear                     year          1996
 w       week of weekyear             number        27
 e       day of week                  number        2
 E       day of week                  text          Tuesday; Tue

 y       year                         year          1996
 D       day of year                  number        189
 M       month of year                month         July; Jul; 07
 d       day of month                 number        10

 a       halfday of day               text          PM
 K       hour of halfday (0~11)       number        0
 h       clockhour of halfday (1~12)  number        12

 H       hour of day (0~23)           number        0
 k       clockhour of day (1~24)      number        24
 m       minute of hour               number        30
 s       second of minute             number        55
 S       fraction of second           number        978

 z       time zone                    text          Pacific Standard Time; PST
 Z       time zone offset/id          zone          -0800; -08:00; America/Los_Angeles

 '       escape for text              delimiter
 ''      single quote                 literal       '
*/

// Format formats a date based on joda conventions
func Format(format string, date time.Time) string {
	formatRune := []rune(format)
	lenFormat := len(formatRune)
	out := ""
	for i := 0; i < len(formatRune); i++ {
		switch r := formatRune[i]; r {
		case 'Y', 'y', 'x': // Y YYYY YY year

			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1

			switch j {
			case 1, 3, 4: // Y YYY YYY
				out += strconv.Itoa(date.Year())
			case 2: // YY
				out += strconv.Itoa(date.Year())[2:4]
			}

		case 'D': // D DD day of year
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1

			switch j {
			case 1: // D
				out = fmt.Sprintf("%s%d", out, date.YearDay())
			case 2: // DD
				out = fmt.Sprintf("%s%02d", out, date.YearDay())
			}

		case 'w': // w ww week of weekyear
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			_, w := date.ISOWeek()
			switch j {
			case 1: // w
				out = fmt.Sprintf("%s%d", out, w)
			case 2: // ww
				out = fmt.Sprintf("%s%02d", out, w)
			}

		case 'M': // M MM MMM MMMM month of year
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Month()
			switch j {
			case 1: // M
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // MM
				out = fmt.Sprintf("%s%02d", out, v)
			case 3: // MMM
				out = fmt.Sprintf("%s%s", out, v.String()[0:3])
			case 4: // MMMM
				out = fmt.Sprintf("%s%s", out, v.String())
			}

		case 'd': // d dd day of month
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Day()
			switch j {
			case 1: // d
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // dd
				out = fmt.Sprintf("%s%02d", out, v)
			}

		case 'e': // e ee day of week(number)
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Weekday()
			switch j {
			case 1: // e
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // ee
				out = fmt.Sprintf("%s%02d", out, v)
			}

		case 'E': // E EE
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Weekday()
			switch j {
			case 1, 2, 3: // E
				out = fmt.Sprintf("%s%s", out, v.String()[0:3])
			case 4: // EE
				out = fmt.Sprintf("%s%s", out, v.String())
			}
		case 'h': // h hh clockhour of halfday (1~12)
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Hour()
			if v > 12 {
				v = v - 11
			}
			switch j {
			case 1: // h
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // hh
				out = fmt.Sprintf("%s%02d", out, v)
			}

		case 'H': // H HH
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Hour()
			switch j {
			case 1: // H
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // HH
				out = fmt.Sprintf("%s%02d", out, v)
			}

		case 'a': // a
			if date.Hour() > 12 {
				out = fmt.Sprintf("%s%s", out, "PM")
			} else {
				out = fmt.Sprintf("%s%s", out, "AM")
			}

		case 'm': // m mm minute of hour
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Minute()
			switch j {
			case 1: // m
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // mm
				out = fmt.Sprintf("%s%02d", out, v)
			}
		case 's': // s ss
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Second()
			switch j {
			case 1: // s
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // ss
				out = fmt.Sprintf("%s%02d", out, v)
			}

		case 'S': // S SS SSS
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Nanosecond()
			switch j {
			case 1: // S
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // SS
				out = fmt.Sprintf("%s%02d", out, v)
			case 3: // SSS
				out = fmt.Sprintf("%s%03d", out, v)
			}

		case 'z': // z
			z, _ := date.Zone()
			out = fmt.Sprintf("%s%s", out, z)

		case 'Z': // Z ZZ
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			zs, z := date.Zone()
			sign := "+"
			if z < 0 {
				sign = "-"
				z = -z
			}
			switch j {
			case 1: // Z
				out = fmt.Sprintf("%s%s%02d00", out, sign, z/3600)
			case 2: // ZZ
				out = fmt.Sprintf("%s%s%02d:00", out, sign, z/3600)
			case 3: // ZZZ
				out = fmt.Sprintf("%s%s", out, timeZone[zs])
			}

		case 'G': //era                          text
			out = fmt.Sprintf("%sAD", out)

		case 'C': //century of era (>=0)         number
			out += strconv.Itoa(date.Year())[0:2]

		case 'K': // K KK hour of halfday (0~11)
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Hour()
			if v > 12 {
				v = v - 12
			}
			switch j {
			case 1: // K
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // KK
				out = fmt.Sprintf("%s%02d", out, v)
			}

		case 'k': // k kk clockhour of day (1~24)
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}

			}
			i = i + j - 1
			v := date.Hour() + 1
			switch j {
			case 1: // k
				out = fmt.Sprintf("%s%d", out, v)
			case 2: // kk
				out = fmt.Sprintf("%s%02d", out, v)
			}
		case '\'': // ' (text delimiter)  or '' (real quote)

			// real quote
			if formatRune[i+1] == r {
				out = fmt.Sprintf("%s'", out)
				i = i + 1
				continue
			}

			tmp := []rune{}
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					tmp = append(tmp, formatRune[i+j])
					continue
				}
				break
			}
			i = i + j

			out = fmt.Sprintf("%s%s", out, string(tmp))

		default:
			out = fmt.Sprintf("%s%s", out, string(r))
		}
	}
	return out

}

var timeZone = map[string]string{
	"GMT":     "Europe/London",
	"BST":     "Europe/London",
	"BSDT":    "Europe/London",
	"CET":     "Europe/Paris",
	"UTC":     "",
	"PST":     "America/Los_Angeles",
	"PDT":     "America/Los_Angeles",
	"LA":      "America/Los_Angeles",
	"LAX":     "America/Los_Angeles",
	"MST":     "America/Denver",
	"MDT":     "America/Denver",
	"CST":     "America/Chicago",
	"CDT":     "America/Chicago",
	"Chicago": "America/Chicago",
	"EST":     "America/New_York",
	"EDT":     "America/New_York",
	"NYC":     "America/New_York",
	"NY":      "America/New_York",
	"AEST":    "Australia/Sydney",
	"AEDT":    "Australia/Sydney",
	"AWST":    "Australia/Perth",
	"AWDT":    "Australia/Perth",
	"ACST":    "Australia/Adelaide",
	"ACDT":    "Australia/Adelaide",
}
