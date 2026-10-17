package tunnel

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTunnelManager_RegexCapture(t *testing.T) {
	tmpDir := t.TempDir()
	urlFile := filepath.Join(tmpDir, "backend-url.txt")

	mgr := NewManager("http://localhost:8080", urlFile, "", "")
	capture := &logCapture{manager: mgr}

	logLine := "2026-10-16T10:00:00Z INF | Your quick Tunnel has been created! Visit it at (it may take some time to be reachable): https://abc-123.trycloudflare.com\n"
	n, err := capture.Write([]byte(logLine))
	require.NoError(t, err)
	assert.Equal(t, len(logLine), n)

	select {
	case url := <-mgr.URLChan():
		assert.Equal(t, "https://abc-123.trycloudflare.com", url)
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for URL on channel")
	}

	assert.Equal(t, "https://abc-123.trycloudflare.com", mgr.GetURL())
}
