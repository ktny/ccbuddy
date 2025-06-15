package buddy

import (
	"errors"
	"time"
)

// State represents the current state of a buddy
type State string

const (
	// StateEgg represents an unhatched buddy
	StateEgg State = "egg"
	// StateHatched represents a hatched buddy
	StateHatched State = "hatched"
)

// Buddy represents a virtual pet that grows based on ClaudeCode usage
type Buddy struct {
	// State is the current lifecycle state of the buddy
	State State `json:"state"`
	// Health is the health level (0-100)
	Health int `json:"health"`
	// CreatedAt is when the buddy was first created
	CreatedAt time.Time `json:"created_at"`
	// LastFedAt is when the buddy was last fed
	LastFedAt time.Time `json:"last_fed_at"`
}

// NewBuddy returns a new Buddy instance initialized in the egg state with full health and current timestamps for creation and last feeding.
func NewBuddy() *Buddy {
	now := time.Now()
	return &Buddy{
		State:     StateEgg,
		Health:    100,
		CreatedAt: now,
		LastFedAt: now,
	}
}

// Hatch transitions the buddy from egg to hatched state
func (b *Buddy) Hatch() error {
	if b.State != StateEgg {
		return errors.New("buddyは卵の状態でのみ孵化できます")
	}

	b.State = StateHatched
	return nil
}

// Feed restores the buddy's health to full and updates LastFedAt
func (b *Buddy) Feed() {
	b.Health = 100
	b.LastFedAt = time.Now()
}

// Age returns how long the buddy has been alive
func (b *Buddy) Age() time.Duration {
	return time.Since(b.CreatedAt)
}

// IsAlive returns true if the buddy is still alive (health > 0)
func (b *Buddy) IsAlive() bool {
	return b.Health > 0
}

// Validate checks if the buddy's state is valid
func (b *Buddy) Validate() error {
	// Validate state
	switch b.State {
	case StateEgg, StateHatched:
		// Valid states
	default:
		return errors.New("不正な状態です")
	}

	// Validate health
	if b.Health < 0 || b.Health > 100 {
		return errors.New("健康状態は0-100の範囲である必要があります")
	}

	// Validate timestamps
	if b.CreatedAt.IsZero() {
		return errors.New("作成日時が設定されていません")
	}

	if b.LastFedAt.IsZero() {
		return errors.New("最後の餌やり時刻が設定されていません")
	}

	return nil
}
