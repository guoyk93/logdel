package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	optToday    string
	optConfDir  string
	optNoDelete bool
)

func exit(err *error) {
	if *err != nil {
		log.Printf("exited with error: %s", (*err).Error())
		os.Exit(1)
	}
}

func main() {
	var err error
	defer exit(&err)

	// flags
	flag.StringVar(&optConfDir, "conf-dir", "/etc/logdel.d", "configuration directory")
	flag.BoolVar(&optNoDelete, "no-delete", false, "don't actual delete files")
	flag.StringVar(&optToday, "today", "", "set today date (in format of YYYY-MM-DD), for debugging")
	flag.Parse()

	// base date
	var today time.Time
	if len(optToday) != 0 {
		if today, err = time.Parse("2006-01-02", optToday); err != nil {
			return
		}
		today = dateMidnight(today)
	} else {
		today = dateMidnight(time.Now())
	}
	log.Printf("today: %s", today.Format(time.RFC3339))

	// iterate
	if err = ruleIterateDir(optConfDir, func(rulefile string, line int, rawPattern string, keep int) {
		var (
			err   error
			files []string

			patterns = expandDoubleWildcard(rawPattern)
		)
		for _, pattern := range patterns {
			if files, err = filepath.Glob(pattern); err != nil {
				log.Printf("- line: %d: pattern: %s is invalid: %s", line, pattern, err.Error())
				return
			}
			for _, file := range files {
				var date time.Time
				var ok bool
				if date, ok = dateFromFilename(file); !ok {
					log.Printf("-- unknown: %s", file)
					continue
				}
				if time.Duration(keep)*time.Hour*24 >= today.Sub(date) {
					log.Printf("-- ok: %s", file)
					continue
				}
				if optNoDelete {
					log.Printf("-- will delete: %s", file)
					continue
				}
				log.Printf("-- delete: %s", file)
				_ = os.Remove(file)
			}
		}
	}); err != nil {
		return
	}
}
