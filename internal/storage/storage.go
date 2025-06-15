package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ktny/ccbuddy/internal/buddy"
)

// Store handles buddy state persistence
type Store struct {
	configDir string
}

// NewStore creates a new storage instance with default config directory
func NewStore() *Store {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// フォールバック: カレントディレクトリに.ccbuddyディレクトリを作成
		homeDir = "."
	}

	configDir := filepath.Join(homeDir, ".ccbuddy")
	return &Store{
		configDir: configDir,
	}
}

// SaveBuddy saves the buddy state to a JSON file
func (s *Store) SaveBuddy(b *buddy.Buddy) error {
	if b == nil {
		return errors.New("buddyがnilです")
	}

	// バリデーション
	if err := b.Validate(); err != nil {
		return err
	}

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(s.configDir, 0755); err != nil {
		return err
	}

	// JSONエンコード
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}

	// ファイルに書き込み
	filePath := s.getFilePath()
	return os.WriteFile(filePath, data, 0600)
}

// LoadBuddy loads the buddy state from a JSON file
func (s *Store) LoadBuddy() (*buddy.Buddy, error) {
	filePath := s.getFilePath()

	// ファイルが存在するかチェック
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("buddyファイルが存在しません")
	}

	// ファイル読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// JSONデコード
	var b buddy.Buddy
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}

	// バリデーション
	if err := b.Validate(); err != nil {
		return nil, err
	}

	return &b, nil
}

// BuddyExists checks if a buddy file exists
func (s *Store) BuddyExists() bool {
	filePath := s.getFilePath()
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// getFilePath returns the full path to the buddy JSON file
func (s *Store) getFilePath() string {
	return filepath.Join(s.configDir, "buddy.json")
}
