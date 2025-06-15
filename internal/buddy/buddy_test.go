package buddy

import (
	"testing"
	"time"
)

func TestNewBuddy(t *testing.T) {
	buddy := NewBuddy()

	if buddy.State != StateEgg {
		t.Errorf("新しいbuddyの状態は卵であるべき, got %v", buddy.State)
	}

	if buddy.Health != 100 {
		t.Errorf("新しいbuddyの健康状態は100であるべき, got %d", buddy.Health)
	}

	if buddy.CreatedAt.IsZero() {
		t.Error("CreatedAtが設定されていない")
	}

	if buddy.LastFedAt.IsZero() {
		t.Error("LastFedAtが設定されていない")
	}
}

func TestBuddyHatch(t *testing.T) {
	buddy := NewBuddy()

	// 卵の状態で孵化を試行
	err := buddy.Hatch()
	if err != nil {
		t.Errorf("卵の状態から孵化できるべき: %v", err)
	}

	if buddy.State != StateHatched {
		t.Errorf("孵化後の状態はHatchedであるべき, got %v", buddy.State)
	}

	// 既に孵化済みの場合はエラー
	err = buddy.Hatch()
	if err == nil {
		t.Error("既に孵化済みの場合はエラーになるべき")
	}
}

func TestBuddyFeed(t *testing.T) {
	buddy := NewBuddy()

	// 健康状態を下げる
	buddy.Health = 50

	// 餌やり
	buddy.Feed()

	if buddy.Health != 100 {
		t.Errorf("餌やり後の健康状態は100になるべき, got %d", buddy.Health)
	}

	// LastFedAtが更新されているかチェック
	now := time.Now()
	if buddy.LastFedAt.After(now) || buddy.LastFedAt.Before(now.Add(-time.Second)) {
		t.Error("LastFedAtが正しく更新されていない")
	}
}

func TestBuddyAge(t *testing.T) {
	buddy := NewBuddy()
	buddy.CreatedAt = time.Now().Add(-24 * time.Hour) // 24時間前に作成

	age := buddy.Age()
	if age < 23*time.Hour || age > 25*time.Hour {
		t.Errorf("年齢が正しく計算されていない, got %v", age)
	}
}

func TestBuddyIsAlive(t *testing.T) {
	buddy := NewBuddy()

	// 健康状態が1以上なら生きている
	buddy.Health = 1
	if !buddy.IsAlive() {
		t.Error("健康状態が1以上なら生きているべき")
	}

	// 健康状態が0なら死んでいる
	buddy.Health = 0
	if buddy.IsAlive() {
		t.Error("健康状態が0なら死んでいるべき")
	}
}

func TestBuddyValidation(t *testing.T) {
	tests := []struct {
		name    string
		buddy   Buddy
		wantErr bool
	}{
		{
			name: "正常なbuddy",
			buddy: Buddy{
				State:     StateEgg,
				Health:    100,
				CreatedAt: time.Now(),
				LastFedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "不正な状態",
			buddy: Buddy{
				State:     "invalid",
				Health:    100,
				CreatedAt: time.Now(),
				LastFedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "不正な健康状態（負数）",
			buddy: Buddy{
				State:     StateEgg,
				Health:    -1,
				CreatedAt: time.Now(),
				LastFedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "不正な健康状態（101以上）",
			buddy: Buddy{
				State:     StateEgg,
				Health:    101,
				CreatedAt: time.Now(),
				LastFedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.buddy.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
