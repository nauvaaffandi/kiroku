// Package webpush implements the real Web Push (RFC 8291 / RFC 8292) sender using
// only the Go standard library. No external push-delivery dependency is required:
// VAPID JWT (ES256), ECDH P-256 key agreement, HKDF-SHA-256 and AES-GCM all come from
// the standard library.
package webpush

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/ecdh"
	"crypto/elliptic"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/domain"
)

// Sender delivers a push notification to a single browser subscription over HTTP.
type Sender struct {
	client      *http.Client
	vapidPublic string
	vapidPriv   *ecdsa.PrivateKey
	vapidSubj   string
}

// NewSender builds a Web Push sender from the configured VAPID keys.
func NewSender() (*Sender, error) {
	var priv *ecdsa.PrivateKey
	if config.VAPIDPrivateKey() != "" {
		d, err := base64.RawURLEncoding.DecodeString(config.VAPIDPrivateKey())
		if err != nil {
			log.Warn().Err(err).Msg("webpush: failed to decode VAPID private key; push delivery unavailable")
		} else {
			priv = ecdsaP256FromRaw(d)
			if priv == nil {
				log.Warn().Msg("webpush: invalid VAPID private key")
			}
		}
	}
	return &Sender{
		client:      &http.Client{Timeout: 30 * time.Second},
		vapidPublic: config.VAPIDPublicKey(),
		vapidPriv:   priv,
		vapidSubj:   config.VAPIDSubject(),
	}, nil
}

// Send delivers a notification to a subscription.
func (s *Sender) Send(ctx context.Context, sub *domain.PushSubscription, n *domain.Notification) error {
	if s.vapidPriv == nil {
		return errors.New("webpush: VAPID private key not configured")
	}
	if sub == nil || sub.Endpoint == "" {
		return errors.New("webpush: missing push subscription")
	}
	if sub.P256DH == "" || sub.Auth == "" {
		return errors.New("webpush: missing subscription encryption keys")
	}

	payload, err := json.Marshal(notificationToPushPayload(n))
	if err != nil {
		return fmt.Errorf("webpush: marshal payload: %w", err)
	}
	return s.sendEncrypted(ctx, sub, payload)
}

func (s *Sender) sendEncrypted(ctx context.Context, sub *domain.PushSubscription, plaintext []byte) error {
	jwt, err := s.vapidJWT(sub.Endpoint)
	if err != nil {
		return fmt.Errorf("webpush: vapid jwt: %w", err)
	}

	encrypted, err := encryptAESGCM(sub, plaintext)
	if err != nil {
		return fmt.Errorf("webpush: encrypt: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.Endpoint, bytes.NewReader(encrypted))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "vapid t="+jwt+", k="+s.vapidPublic)
	req.Header.Set("Content-Encoding", "aesgcm")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("TTL", "86400")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("webpush: request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webpush: delivery failed: status=%d body=%q", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

// notificationToPushPayload builds a small, safe push payload. No financial data.
func notificationToPushPayload(n *domain.Notification) map[string]any {
	payload := map[string]any{
		"title": n.Title,
		"body":  n.Body,
	}
	if n.ActionURL != "" {
		payload["url"] = n.ActionURL
	}
	if n.Data != nil {
		payload["data"] = n.Data
	}
	return payload
}

// vapidJWT builds the VAPID JWT (ES256 signed with the VAPID private key).
func (s *Sender) vapidJWT(endpoint string) (string, error) {
	host := endpoint
	if u, err := url.Parse(endpoint); err == nil && u.Host != "" {
		host = u.Host
	}

	header := map[string]any{"typ": "JWT", "alg": "ES256"}
	now := time.Now()
	claims := map[string]any{"aud": host, "exp": now.Add(12 * time.Hour).Unix(), "sub": s.vapidSubj}

	h, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	c, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	signingInput := base64.RawURLEncoding.EncodeToString(h) + "." + base64.RawURLEncoding.EncodeToString(c)

	r, sBig, err := ecdsa.Sign(rand.Reader, s.vapidPriv, []byte(signingInput))
	if err != nil {
		return "", err
	}

	// P-256 signatures: r and s padded to 32 bytes each.
	rBytes := padBigInt(r, 32)
	sBytes := padBigInt(sBig, 32)

	sig64 := append([]byte{}, rBytes...)
	sig64 = append(sig64, sBytes...)

	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig64), nil
}

// padBigInt left-pads a big integer to n bytes (truncates if longer).
func padBigInt(b *big.Int, n int) []byte {
	out := make([]byte, n)
	bs := b.Bytes()
	if len(bs) >= n {
		return bs[len(bs)-n:]
	}
	copy(out[n-len(bs):], bs)
	return out
}

// encryptAESGCM encrypts plaintext per RFC 8291 with the subscription's keys.
func encryptAESGCM(sub *domain.PushSubscription, plaintext []byte) ([]byte, error) {
	curve := ecdh.P256()

	recipientPub, err := base64.RawURLEncoding.DecodeString(sub.P256DH)
	if err != nil {
		return nil, fmt.Errorf("decode p256dh: %w", err)
	}
	auth, err := base64.RawURLEncoding.DecodeString(sub.Auth)
	if err != nil {
		return nil, fmt.Errorf("decode auth: %w", err)
	}

	recvPub, err := curve.NewPublicKey(recipientPub)
	if err != nil {
		return nil, fmt.Errorf("parse recipient public key: %w", err)
	}

	ephemeralKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	_ = ephemeralKey.PublicKey

	sharedSecret, err := ephemeralKey.ECDH(recvPub)
	if err != nil {
		return nil, fmt.Errorf("ecdh: %w", err)
	}

	// Salt is the 16-byte auth token; pad if shorter.
	salt := auth
	if len(salt) < 16 {
		padded := make([]byte, 16)
		copy(padded, auth)
		salt = padded
	}

	// PRK = HKDF-Extract(salt, ikm=sharedSecret)
	prk, err := hkdf.Extract(sha256.New, salt, sharedSecret)
	if err != nil {
		return nil, err
	}

	cek, err := hkdf.Expand(sha256.New, prk, "Content-Encoding: aesgcm\x00", 32)
	if err != nil {
		return nil, err
	}
	nonce, err := hkdf.Expand(sha256.New, prk, "Content-Encoding: nonce\x00", 12)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// AES-GCM Seal appends the 16-byte tag to the ciphertext.
	encryptedData := aead.Seal(nil, nonce, plaintext, nil)

	// Record format (RFC 8291): 0x01 || 2-byte BE length || encrypted data
	record := make([]byte, 1+2+len(encryptedData))
	record[0] = 0x01
	record[1] = byte(len(encryptedData) >> 8)
	record[2] = byte(len(encryptedData))
	copy(record[3:], encryptedData)

	body := make([]byte, len(record))
	copy(body, record)
	return body, nil
}

// ecdsaP256FromRaw reconstructs an ECDSA P-256 private key from raw private bytes
// (as stored from VAPID_PRIVATE_KEY), computing X,Y via scalar base multiplication.
func ecdsaP256FromRaw(d []byte) *ecdsa.PrivateKey {
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: elliptic.P256()},
		D:         new(big.Int).SetBytes(d),
	}
	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P256().ScalarBaseMult(d)
	if priv.PublicKey.X == nil || priv.PublicKey.Y == nil {
		return nil
	}
	return priv
}
