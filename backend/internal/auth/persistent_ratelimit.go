package auth

import (
	"database/sql"
	"time"
)

type PersistentRateLimiter struct {
	db *sql.DB
}

func NewPersistentRateLimiter(database *sql.DB) *PersistentRateLimiter {
	return &PersistentRateLimiter{db: database}
}

type LimitStatus struct {
	Blocked      bool
	CooldownUntil time.Time
	Tier         int
}

func (p *PersistentRateLimiter) CheckAndRecord(identifier string) (LimitStatus, error) {
	now := time.Now()
	var id int64
	var attempts int
	var tier int
	var cooldownUntil sql.NullTime

	err := p.db.QueryRow(
		"SELECT id, attempts, tier, cooldown_until FROM login_attempts WHERE identifier = ?",
		identifier,
	).Scan(&id, &attempts, &tier, &cooldownUntil)

	if err == sql.ErrNoRows {
		_, err = p.db.Exec(
			"INSERT INTO login_attempts (identifier, attempts, tier, last_failure, updated_at) VALUES (?, 1, 0, ?, ?)",
			identifier, now, now,
		)
		return LimitStatus{Blocked: false}, err
	} else if err != nil {
		return LimitStatus{}, err
	}

	if cooldownUntil.Valid && cooldownUntil.Time.After(now) {
		return LimitStatus{
			Blocked:       true,
			CooldownUntil: cooldownUntil.Time,
			Tier:          tier,
		}, nil
	}

	// If cooldown expired, reset or increment tier
	if cooldownUntil.Valid && cooldownUntil.Time.Before(now) {
		attempts = 0
	}

	attempts++
	newTier := tier
	var newCooldown time.Time

	if attempts >= 10 {
		newTier = 3
		newCooldown = now.Add(24 * time.Hour)
	} else if attempts >= 5 {
		newTier = 2
		newCooldown = now.Add(3 * time.Hour)
	} else if attempts >= 3 {
		newTier = 1
		newCooldown = now.Add(1 * time.Hour)
	}

	blocked := !newCooldown.IsZero() && newCooldown.After(now)

	_, err = p.db.Exec(
		"UPDATE login_attempts SET attempts = ?, tier = ?, cooldown_until = ?, last_failure = ?, updated_at = ? WHERE id = ?",
		attempts, newTier, func() interface{} {
			if newCooldown.IsZero() {
				return nil
			}
			return newCooldown
		}(), now, now, id,
	)

	return LimitStatus{
		Blocked:       blocked,
		CooldownUntil: newCooldown,
		Tier:          newTier,
	}, err
}

func (p *PersistentRateLimiter) Reset(identifier string) error {
	_, err := p.db.Exec("DELETE FROM login_attempts WHERE identifier = ?", identifier)
	return err
}
