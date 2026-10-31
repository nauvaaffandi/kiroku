package tunnel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Manager struct {
	mu           sync.Mutex
	cmd          *exec.Cmd
	ctx          context.Context
	cancel       context.CancelFunc
	url          string
	urlFile      string
	restartCount int
	upstashURL   string
	upstashToken string
	localURL     string
	urlChan      chan string
}

func NewManager(localURL, urlFile, upstashURL, upstashToken string) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		ctx:          ctx,
		cancel:       cancel,
		urlFile:      urlFile,
		upstashURL:   upstashURL,
		upstashToken: upstashToken,
		localURL:     localURL,
		urlChan:      make(chan string, 10),
	}
}

func (m *Manager) URLChan() <-chan string {
	return m.urlChan
}

func (m *Manager) Start() error {
	if m.upstashURL == "" || m.upstashToken == "" {
		log.Warn().Msg("tunnel: Upstash credentials not set, skipping tunnel")
		return nil
	}

	go m.run()
	return nil
}

func (m *Manager) run() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			if err := m.runOnce(); err != nil {
				log.Error().Err(err).Msg("tunnel: run failed")
			}

			m.mu.Lock()
			m.restartCount++
			count := m.restartCount
			m.mu.Unlock()

			if count >= 3 {
				log.Error().Int("restarts", count).Msg("tunnel: max restarts reached, giving up")
				return
			}

			backoff := time.Duration(count) * 10 * time.Second
			log.Info().Dur("backoff", backoff).Int("attempt", count+1).Msg("tunnel: restarting after backoff")
			select {
			case <-m.ctx.Done():
				return
			case <-time.After(backoff):
			}
		}
	}
}

func (m *Manager) runOnce() error {
	os.MkdirAll(filepath.Dir(m.urlFile), 0755)
	logFile := filepath.Join(filepath.Dir(m.urlFile), "..", "logs", "cloudflared.log")
	os.MkdirAll(filepath.Dir(logFile), 0755)
	logOut, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	defer logOut.Close()

	m.cmd = exec.CommandContext(m.ctx, "cloudflared", "tunnel", "--url", m.localURL, "--no-autoupdate")
	m.cmd.Stdout = io.MultiWriter(logOut, &logCapture{manager: m})
	m.cmd.Stderr = io.MultiWriter(logOut, &logCapture{manager: m})

	log.Info().Str("url", m.localURL).Msg("tunnel: starting cloudflared")
	if err := m.cmd.Run(); err != nil {
		if m.ctx.Err() != nil {
			return nil // Context cancelled, normal shutdown
		}
		return fmt.Errorf("cloudflared exited: %w", err)
	}
	return nil
}

type logCapture struct {
	manager *Manager
}

func (l *logCapture) Write(p []byte) (n int, err error) {
	output := string(p)
	urlRegex := regexp.MustCompile(`https://[a-z0-9-]+\.trycloudflare\.com`)
	if matches := urlRegex.FindStringSubmatch(output); len(matches) > 0 {
		url := matches[0]
		l.manager.setURL(url)
	}
	return len(p), nil
}

func (m *Manager) setURL(url string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.url = url
	m.restartCount = 0

	log.Info().Str("url", url).Msg("tunnel: URL detected / forced update")

	// Write to file
	if err := os.WriteFile(m.urlFile, []byte(url), 0644); err != nil {
		log.Error().Err(err).Msg("tunnel: write url file")
	}

	// Notify channel non-blockingly
	select {
	case m.urlChan <- url:
	default:
	}

	// Push to Upstash
	go m.pushToUpstash(url)
}

func (m *Manager) pushToUpstash(url string) {
	if m.upstashURL == "" || m.upstashToken == "" {
		return
	}

	body, err := json.Marshal([]any{"SET", "kiroku:backend:url", url, "EX", 86400})
	if err != nil {
		log.Error().Err(err).Msg("tunnel: marshal upstash payload")
		return
	}

	req, err := http.NewRequest("POST", m.upstashURL, bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Msg("tunnel: create upstash request")
		return
	}
	req.Header.Set("Authorization", "Bearer "+m.upstashToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("tunnel: upstash request failed")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Msg("tunnel: upstash non-OK response")
		return
	}

	log.Info().Msg("tunnel: URL pushed to Upstash")
}

func (m *Manager) Stop() {
	m.cancel()
	m.mu.Lock()
	cmd := m.cmd
	m.mu.Unlock()
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
	}
}

func (m *Manager) GetURL() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.url
}
