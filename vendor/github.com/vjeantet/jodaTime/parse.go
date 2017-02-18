package jodaTime

import (
	"fmt"
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
	for i := 0; i < len(formatRune); i++ {
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
				layout = fmt.Sprintf("%s3", layout)
			default:
				layout = fmt.Sprintf("%s03", layout)
			}

			i = i + j - 1
		case 'H':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			layout = fmt.Sprintf("%s15", layout)

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
				layout = fmt.Sprintf("%s4", layout)
			default:
				layout = fmt.Sprintf("%s04", layout)
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
				layout = fmt.Sprintf("%s5", layout)
			default:
				layout = fmt.Sprintf("%s05", layout)
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
				layout = fmt.Sprintf("%s2", layout)
			default:
				layout = fmt.Sprintf("%s02", layout)
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
				layout = fmt.Sprintf("%sMon", layout)
			default:
				layout = fmt.Sprintf("%sMonday", layout)
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
				layout = fmt.Sprintf("%s1", layout)
			case 2:
				layout = fmt.Sprintf("%s01", layout)
			case 3:
				layout = fmt.Sprintf("%sJan", layout)
			case 4:
				layout = fmt.Sprintf("%sJanuary", layout)

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
				layout = fmt.Sprintf("%s06", layout)
			default: // dd
				layout = fmt.Sprintf("%s2006", layout)
			}

			i = i + j - 1

		case 'S':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			layout = fmt.Sprintf("%s%s", layout, strings.Repeat("9", j))

			i = i + j - 1

		case 'a':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			layout = fmt.Sprintf("%sPM", layout)
		case 'Z':
			j := 1
			for ; i+j < lenFormat; j++ {
				if formatRune[i+j] != r {
					break
				}
			}

			switch j {
			case 1: // d
				layout = fmt.Sprintf("%s-0700", layout)
			case 2: // d
				layout = fmt.Sprintf("%s-07:00", layout)
			}

			i = i + j - 1
		default:
			layout = fmt.Sprintf("%s%s", layout, string(r))
		}
	}
	return layout

}
