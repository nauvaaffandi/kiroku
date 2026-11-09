package logger

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
)

type dailyRotator struct {
	mu     sync.Mutex
	dir    string
	day    string
	writer io.Writer
}

func newDailyRotator(dir string) *dailyRotator {
	return &dailyRotator{dir: dir}
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func (d *dailyRotator) activePath() string {
	return filepath.Join(d.dir, today()+logFileExt)
}

func (d *dailyRotator) newWriter(path string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    parseMaxSize(),
		MaxBackups: parseMaxBackups(),
		MaxAge:     parseMaxAge(),
		Compress:   parseCompress(),
		LocalTime:  false,
	}
}

func (d *dailyRotator) Write(p []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	day := today()
	if d.writer == nil || day != d.day {
		d.day = day
		if err := os.MkdirAll(d.dir, 0o755); err != nil {
			d.writer = io.Discard
			return 0, err
		}
		d.writer = d.newWriter(d.activePath())
	}
	return d.writer.Write(p)
}
