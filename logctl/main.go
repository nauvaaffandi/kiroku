package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const defaultLogsDir = "backend/logs"

func logsDir() string {
	if d := os.Getenv("KIROKU_LOG_DIR"); d != "" {
		return d
	}
	return defaultLogsDir
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "list":
		cmdList(os.Args[2:])
	case "search":
		cmdSearch(os.Args[2:])
	case "filter":
		cmdFilter(os.Args[2:])
	case "clean":
		cmdClean(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "usage: logctl <command> [args]\n")
	fmt.Fprintf(os.Stderr, "commands:\n")
	fmt.Fprintf(os.Stderr, "  list                                  list log files with sizes\n")
	fmt.Fprintf(os.Stderr, "  search [query]                        search text in logs (supports field=value)\n")
	fmt.Fprintf(os.Stderr, "  filter --from YYYY-MM-DD --to YYYY-MM-DD  filter logs by date range\n")
	fmt.Fprintf(os.Stderr, "  clean                                 delete expired logs per LOG_RETENTION\n")
}

type logEntry map[string]interface{}

func parseLogLine(line string) (logEntry, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false
	}
	var e logEntry
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return nil, false
	}
	return e, true
}

func renderEntry(e logEntry) string {
	ts, _ := e["time"].(string)
	level, _ := e["level"].(string)
	msg, _ := e["message"].(string)
	var parts []string
	if ts != "" {
		parts = append(parts, ts)
	}
	if level != "" {
		parts = append(parts, level)
	}
	if msg != "" {
		parts = append(parts, msg)
	}
	out := strings.Join(parts, " ")
	delete(e, "time")
	delete(e, "level")
	delete(e, "message")
	if len(e) > 0 {
		b, _ := json.Marshal(e)
		out = out + " " + string(b)
	}
	return out
}

func logFiles() ([]string, error) {
	entries, err := os.ReadDir(logsDir())
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".log") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)
	return files, nil
}

func dateFromFile(name string) (time.Time, bool) {
	base := strings.TrimSuffix(name, ".log")
	t, err := time.Parse("2006-01-02", base)
	return t, err == nil
}

func cmdList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Parse(args)

	files, err := logFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading logs dir: %v\n", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("(no log files)")
		return
	}
	fmt.Printf("%-20s %10s\n", "FILE", "SIZE")
	fmt.Println(strings.Repeat("-", 32))
	for _, name := range files {
		info, err := os.Stat(filepath.Join(logsDir(), name))
		size := "-"
		if err == nil {
			size = fmt.Sprintf("%d", info.Size())
		}
		fmt.Printf("%-20s %10s\n", name, size)
	}
}

func cmdSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: search requires a query")
		os.Exit(2)
	}
	query := strings.ToLower(fs.Arg(0))
	fieldMode := strings.Contains(query, "=")

	files, err := logFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading logs dir: %v\n", err)
		os.Exit(1)
	}

	count := 0
	for _, name := range files {
		path := filepath.Join(logsDir(), name)
		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			match := false
			if fieldMode {
				kv := strings.SplitN(query, "=", 2)
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				if e, ok := parseLogLine(line); ok {
					if v, ok := e[key]; ok && fmt.Sprintf("%v", v) == val {
						match = true
					}
				}
			} else {
				if strings.Contains(strings.ToLower(line), query) {
					match = true
				}
			}
			if match {
				if e, ok := parseLogLine(line); ok {
					fmt.Printf("[%s] %s\n", name, renderEntry(e))
				} else {
					fmt.Printf("[%s] %s\n", name, line)
				}
				count++
			}
		}
		f.Close()
	}
	fmt.Printf("\n%d match(es)\n", count)
}

func cmdFilter(args []string) {
	fs := flag.NewFlagSet("filter", flag.ExitOnError)
	from := fs.String("from", "", "start date YYYY-MM-DD")
	to := fs.String("to", "", "end date YYYY-MM-DD")
	fs.Parse(args)

	var fromT, toT time.Time
	var hasFrom, hasTo bool
	if *from != "" {
		fromT, _ = time.Parse("2006-01-02", *from)
		hasFrom = true
	}
	if *to != "" {
		toT, _ = time.Parse("2006-01-02", *to)
		hasTo = true
	}

	files, err := logFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading logs dir: %v\n", err)
		os.Exit(1)
	}

	count := 0
	for _, name := range files {
		t, ok := dateFromFile(name)
		if !ok {
			continue
		}
		if hasFrom && t.Before(fromT) {
			continue
		}
		if hasTo && t.After(toT) {
			continue
		}
		path := filepath.Join(logsDir(), name)
		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			continue
		}
		fmt.Printf("== %s ==\n", name)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if e, ok := parseLogLine(line); ok {
				fmt.Println(renderEntry(e))
			} else {
				fmt.Println(line)
			}
			count++
		}
		f.Close()
	}
	fmt.Printf("\n%d entr(y|ies)\n", count)
}

func parseRetentionDays(s string) int {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "d")
	s = strings.TrimSuffix(s, "D")
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 90
	}
	return n
}

func cmdClean(args []string) {
	fs := flag.NewFlagSet("clean", flag.ExitOnError)
	fs.Parse(args)

	retentionDays := parseRetentionDays(os.Getenv("LOG_RETENTION"))
	if retentionDays <= 0 {
		retentionDays = 90
	}

	entries, err := os.ReadDir(logsDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading logs dir: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	removed := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".log") {
			continue
		}
		t, ok := dateFromFile(e.Name())
		if !ok {
			continue
		}
		if now.Sub(t) > time.Duration(retentionDays)*24*time.Hour {
			os.Remove(filepath.Join(logsDir(), e.Name()))
			removed++
		}
	}
	fmt.Printf("removed %d expired log file(s) (retention: %dd)\n", removed, retentionDays)
}
