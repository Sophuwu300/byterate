package main

import (
	"fmt"
	"os"
	"sophuwu.site/byterate"
	"strings"
	"time"
)

type inputs struct {
	nums []string
	help bool
	dur  bool
	time bool
}

func parseArgs(args []string) (i inputs, err error) {
	opts := ""
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			opts += arg[2:3]
		} else if strings.HasPrefix(arg, "-") {
			opts += arg[1:]
		} else {
			i.nums = append(i.nums, arg)
		}
	}
	for _, opt := range opts {
		switch opt {
		case 'h':
			i.help = true
		case 'd':
			i.dur = true
		case 't':
			i.time = true
		}
	}
	if !i.dur && !i.time {
		i.dur = true
	}
	if len(i.nums) != 2 {
		err = fmt.Errorf("expected 2 arguments, got %d", len(i.nums))
	}
	return
}

func tailS(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func main() {
	var (
		s, r byterate.Size
		t    time.Time
		d    time.Duration
		i    inputs
		err  error
		prnt string
	)

	i, err = parseArgs(os.Args[1:])

	if i.help {
		fmt.Println(strings.ReplaceAll(helpText, "{{ name }}", os.Args[0]))
		return
	}

	if err != nil {
		goto bad
	}

	s, err = byterate.ParseSize(i.nums[0])
	if err != nil {
		goto bad
	}

	r, err = byterate.ParseSize(i.nums[1])
	if err != nil {
		goto bad
	}

	t, d, err = byterate.Time(s, r)
	if err != nil {
		goto bad
	}

	if i.dur {
		if d.Hours() > 23 {
			days := int(d.Hours() / 24)
			if days > 365 {
				years := days / 365
				prnt = fmt.Sprintf("%d year%s, ", years, tailS(years))
				days %= 365
			}
			hrs := int(d.Hours()) % 24
			prnt += fmt.Sprintf("%d day%s, %02d hour%s, ", days, tailS(days), hrs, tailS(hrs))
		} else if d.Hours() >= 1 {
			prnt = fmt.Sprintf("%02d hour%s, ", int(d.Hours()), tailS(int(d.Hours())))
		}
		mins := int(d.Minutes()) % 60
		if prnt != "" || mins > 0 {
			prnt += fmt.Sprintf("%02d minute%s and ", mins, tailS(mins))
		}
		secs := int(d.Seconds()) % 60
		prnt += fmt.Sprintf("%02d second%s", secs, tailS(secs))
		fmt.Println(prnt)
	}
	if i.time {
		prnt = "Today "
		ndate := time.Now().Truncate(time.Hour * 24)
		tdate := t.Truncate(time.Hour * 24)
		if !tdate.Equal(ndate) {
			if tdate.Equal(ndate.AddDate(0, 0, 1)) {
				prnt = "Tomorrow "
			} else if t.Before(ndate.AddDate(0, 0, 6)) {
				prnt = t.Weekday().String() + " at " + prnt
			} else if t.Year() == ndate.Year() {
				prnt = "02 Jan "
			} else {
				prnt = "02 Jan 2006 "
			}
		}
		fmt.Println(t.Format(prnt + "15:04"))
	}
	return

bad:
	fmt.Printf("Usage: %s <size> <rate>\n", os.Args[0])
	fmt.Println("\t--help for more information")
}

const helpText = `{{ name }}:
	Print the time it takes to transfer a file at a given rate.

	Usage: {{ name }} [options] <size> <rate>

	Options:
		-h	--help
			Show this help message.
		-d	--duration
			Print the duration the transfer will take.
		-t	--time
			Print the time the transfer will finish if started now.

	if no options are given, duration will be printed.

	<size> is the size of the file to transfer.
	<rate> is the transfer rate as a size, always per second.

	SI prefixes are supported (e.g. 1k = 1000, 1ki = 1024)
	Supported units are: b for bytes, bit for bits.
	If no unit is given, bytes are assumed.

	Examples:
		10 GiB at 120 mbps : {{ name }} 10gib 120mbit
		16 MiB at 1.2 MiB/s: {{ name }} 16mib 1.2mib
		15 GB  at 1.5 MB/s : {{ name }} 1.5g 1.5m
`
