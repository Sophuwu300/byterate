package byterate

import (
	"errors"
	"strings"
	"time"
)

type Size uint64

const Prefix = "bkmgtpezy"

// parseFloat parses a float from a string and return the ending string
func parseFloat(s string) (float64, string, error) {
	if len(s) < 1 {
		return 0, s, errors.New("empty string")
	}
	var (
		f    float64 = 0
		sign bool    = false
		c    rune
		i    int
	)
	if s[0] == '-' {
		sign = true
		s = s[1:]
	}
	for i, c = range s {
		if c >= '0' && c <= '9' {
			f = f*10 + float64(c-'0')
		} else {
			break
		}
	}
	if s[i] == '.' {
		s = s[i+1:]
		var m float64 = 10
		for i, c = range s {
			if c >= '0' && c <= '9' {
				f += float64(c-'0') / m
				m *= 10
			} else {
				break
			}
		}
	}
	if sign {
		f = -f
	}
	return f, s[i:], nil
}

// ParseSize parses a string into a Size
func ParseSize(s string) (Size, error) {
	if len(s) < 1 {
		return 0, errors.New("empty string")
	}
	var (
		num float64
		err error
	)
	num, s, err = parseFloat(s)
	if err != nil {
		return 0, err
	}
	s = strings.ToLower(s)
	if strings.HasSuffix(s, "bit") {
		s = strings.TrimSuffix(s, "bit")
	} else {
		num *= 8
	}
	if len(s) > 0 {
		i := strings.Index(Prefix, s[0:1])
		if i > 0 {
			if strings.Contains(s, "i") {
				i = 1 << (10 * i)
				num *= float64(i)
			} else {
				for i > 0 {
					num *= 1e3
					i--
				}
			}
		}
	}
	return Size(num), nil
}

/*
Time returns the expected end time and duration of a transfer given the size and rate.
rate is size per second
*/
func Time(size, rate Size) (endTime time.Time, duration time.Duration, err error) {
	if rate == 0 {
		err = errors.New("rate must be greater than 0")
		return
	}
	if rate > size {
		err = errors.New("rate must be less than Size")
		return
	}
	duration = time.Duration(size/rate) * time.Second
	endTime = time.Now().Add(duration)
	return
}
