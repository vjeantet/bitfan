// durationfmt provides a function to format durations according to a format
// string.
package durationfmt

import (
	"fmt"
	"time"
)

const Day = 24 * time.Hour
const Week = 7 * Day
const Year = 365 * Day

// durationUnit represets a possible duration unit. A durationUnit object
// contains the divisor that the duration unit uses as well as if that duration
// unit is present in the duration format.
type durationUnit struct {
	Present    bool
	DurDivisor time.Duration
}

// Format formats the given duration according to the given format string.
// %y - # of years
// %w - # of weeks
// %d - # of days
// %h - # of hours
// %m - # of minutes
// %s - # of seconds
// %% - print a percent sign
// You can place a 0 before the h, m, and s modifiers to zeropad those values to
// two digits. Zeropadding is undefined for the other modifiers.
func Format(dur time.Duration, fmtStr string) (string, error) {
	var durUnitSlice = []durationUnit{
		durationUnit{
			DurDivisor: Year,
		},
		durationUnit{
			DurDivisor: Week,
		},
		durationUnit{
			DurDivisor: Day,
		},
		durationUnit{
			DurDivisor: time.Hour,
		},
		durationUnit{
			DurDivisor: time.Minute,
		},
		durationUnit{
			DurDivisor: time.Second,
		},
	}
	var durUnitMap = map[string]*durationUnit{
		"y": &durUnitSlice[0],
		"w": &durUnitSlice[1],
		"d": &durUnitSlice[2],
		"h": &durUnitSlice[3],
		"m": &durUnitSlice[4],
		"s": &durUnitSlice[5],
	}

	sprintfFmt, durCount, err := parseFmtStr(fmtStr, durUnitMap)
	if err != nil {
		return "", err
	}

	durArray := make([]interface{}, durCount)
	calculateDurUnits(dur, durArray, durUnitSlice)

	return fmt.Sprintf(sprintfFmt, durArray...), nil
}

// calculateDurUnits takes a duration and breaks it up into its constituent
// duration unit values.
func calculateDurUnits(dur time.Duration, durArray []interface{}, durUnitSlice []durationUnit) {
	remainingDur := dur
	durCount := 0
	for _, d := range durUnitSlice {
		if d.Present {
			durArray[durCount] = remainingDur / d.DurDivisor
			remainingDur = remainingDur % d.DurDivisor
			durCount++
		}
	}
}

// parseFmtStr parses the given duration format string into its constituent
// units.
// parseFmtStr returns a format string that can be passed to fmt.Sprintf and a
// count of how many duration units are in the format string.
func parseFmtStr(fmtStr string, durUnitMap map[string]*durationUnit) (string, int, error) {
	modifier, zeropad := false, false
	sprintfFmt := ""
	durCount := 0
	for _, c := range fmtStr {
		fmtChar := string(c)
		if modifier == false {
			if fmtChar == "%" {
				modifier = true
			} else {
				sprintfFmt += fmtChar
			}
			continue
		}
		if _, ok := durUnitMap[fmtChar]; ok {
			durUnitMap[fmtChar].Present = true
			durCount++
			if zeropad {
				sprintfFmt += "%02d"
				zeropad = false
			} else {
				sprintfFmt += "%d"
			}
		} else {
			switch fmtChar {
			case "0":
				zeropad = true
				continue
			case "%":
				sprintfFmt += "%%"
			default:
				return "", durCount, fmt.Errorf("incorrect duration modifier")
			}
		}
		modifier = false
	}
	return sprintfFmt, durCount, nil
}
