package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ktny/ccbuddy/internal/buddy"
)

func TestNewStore(t *testing.T) {
	store := NewStore()

	if store.configDir == "" {
		t.Error("configDirが設定されていない")
	}
}

func TestStoreSaveBuddy(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	store := Store{configDir: tempDir}

	buddy := &buddy.Buddy{
		State:     buddy.StateEgg,
		Health:    100,
		CreatedAt: time.Now(),
		LastFedAt: time.Now(),
	}

	err := store.SaveBuddy(buddy)
	if err != nil {
		t.Errorf("SaveBuddy()でエラーが発生: %v", err)
	}

	// ファイルが作成されているかチェック
	filePath := filepath.Join(tempDir, "buddy.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("buddy.jsonファイルが作成されていない")
	}
}

func TestStoreLoadBuddy(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	store := Store{configDir: tempDir}

	originalBuddy := &buddy.Buddy{
		State:     buddy.StateHatched,
		Health:    75,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		LastFedAt: time.Now().Add(-2 * time.Hour),
	}

	// 保存
	err := store.SaveBuddy(originalBuddy)
	if err != nil {
		t.Fatalf("SaveBuddy()でエラーが発生: %v", err)
	}

	// 読み込み
	loadedBuddy, err := store.LoadBuddy()
	if err != nil {
		t.Errorf("LoadBuddy()でエラーが発生: %v", err)
	}

	// 比較
	if loadedBuddy.State != originalBuddy.State {
		t.Errorf("状態が一致しない: expected %v, got %v", originalBuddy.State, loadedBuddy.State)
	}

	if loadedBuddy.Health != originalBuddy.Health {
		t.Errorf("健康状態が一致しない: expected %d, got %d", originalBuddy.Health, loadedBuddy.Health)
	}

	// 時刻は秒単位で比較（JSONシリアライゼーションで精度が変わる可能性があるため）
	if abs(loadedBuddy.CreatedAt.Unix()-originalBuddy.CreatedAt.Unix()) > 1 {
		t.Errorf("作成日時が一致しない: expected %v, got %v", originalBuddy.CreatedAt, loadedBuddy.CreatedAt)
	}

	if abs(loadedBuddy.LastFedAt.Unix()-originalBuddy.LastFedAt.Unix()) > 1 {
		t.Errorf("最後の餌やり時刻が一致しない: expected %v, got %v", originalBuddy.LastFedAt, loadedBuddy.LastFedAt)
	}
}

func TestStoreLoadBuddyNotExists(t *testing.T) {
	// テスト用の一時ディレクトリを作成（ファイルなし）
	tempDir := t.TempDir()
	store := Store{configDir: tempDir}

	buddy, err := store.LoadBuddy()
	if err == nil {
		t.Error("ファイルが存在しない場合はエラーになるべき")
	}

	if buddy != nil {
		t.Error("ファイルが存在しない場合はnilを返すべき")
	}
}

func TestStoreBuddyExists(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	store := Store{configDir: tempDir}

	// ファイルが存在しない場合
	if store.BuddyExists() {
		t.Error("ファイルが存在しない場合はfalseを返すべき")
	}

	// buddyを保存
	buddy := &buddy.Buddy{
		State:     buddy.StateEgg,
		Health:    100,
		CreatedAt: time.Now(),
		LastFedAt: time.Now(),
	}
	err := store.SaveBuddy(buddy)
	if err != nil {
		t.Fatalf("SaveBuddy()でエラーが発生: %v", err)
	}

	// ファイルが存在する場合
	if !store.BuddyExists() {
		t.Error("ファイルが存在する場合はtrueを返すべき")
	}
}

func TestStoreInvalidJSON(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	store := Store{configDir: tempDir}

	// 無効なJSONファイルを作成
	filePath := filepath.Join(tempDir, "buddy.json")
	err := os.WriteFile(filePath, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("テスト用ファイルの作成に失敗: %v", err)
	}

	// 読み込みを試行
	buddy, err := store.LoadBuddy()
	if err == nil {
		t.Error("無効なJSONの場合はエラーになるべき")
	}

	if buddy != nil {
		t.Error("無効なJSONの場合はnilを返すべき")
	}
}

// ヘルパー関数：絶対値を返す
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
