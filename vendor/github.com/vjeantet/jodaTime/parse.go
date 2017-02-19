package jodaTime

import (
	"strings"
	"time"
)

func ParseInLocation(format, value, timezone string) (time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(getLayout(format), value, location)

}

// Parse parses a value into a time.time
func Parse(format, value string) (time.Time, error) {
	return time.Parse(getLayout(format), value)
}

func getLayout(format string) string {
	//replace ? or for rune ?
	formatRune := []rune(format)
	lenFormat := len(formatRune)
	layout := ""
	for i := 0; i < lenFormat; i++ {
		switch r := formatRune[i]; r {
		case 'h':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}
			switch j {
			case 1: // d
				layout += "3"
			default:
				layout += "03"
			}

			i = i + j - 1
		case 'H':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			layout += "15"

			i = i + j - 1
		case 'm':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}
			switch j {
			case 1: // d
				layout += "4"
			default:
				layout += "04"
			}

			i = i + j - 1
		case 's':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}
			switch j {
			case 1: // d
				layout += "5"
			default:
				layout += "05"
			}

			i = i + j - 1
		case 'd':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}
			switch j {
			case 1: // d
				layout += "2"
			default:
				layout += "02"
			}
			i = i + j - 1
		case 'E':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}
			switch j {
			case 1, 2, 3: // d
				layout += "Mon"
			default:
				layout += "Monday"
			}
			i = i + j - 1
		case 'M':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			switch j {
			case 1: // d
				layout += "1"
			case 2:
				layout += "01"
			case 3:
				layout += "Jan"
			case 4:
				layout += "January"

			}
			i = i + j - 1

		case 'Y', 'y', 'x':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}
			switch j {
			case 2: // d
				layout += "06"
			default: // dd
				layout += "2006"
			}

			i = i + j - 1

		case 'S':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			layout += strings.Repeat("9", j)

			i = i + j - 1

		case 'a':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			layout += "PM"
		case 'Z':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			switch j {
			case 1: // d
				layout += "-0700"
			case 2: // d
				layout += "-07:00"
			}

			i = i + j - 1
		case '\'': // ' (text delimiter)  or '' (real quote)

			// real quote
			if formatRune[i+1] == r {
				layout += "'"
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

			layout += string(tmp)
		default:
			layout += string(r)
		}
	}
	return layout

}
