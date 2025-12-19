package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/identicon"
	"github.com/furarico/octo-deck-api/internal/repository"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// GitHubMember represents a member from the GitHub Organization Members API
type GitHubMember struct {
	ID        int64  `json:"id"`
	NodeID    string `json:"node_id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	FullName  string `json:"full_name"`
}

func main() {
	// 環境変数から設定を読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	orgName := os.Getenv("ORG_NAME")
	if orgName == "" {
		orgName = "p2hacks2025"
	}

	log.Printf("Starting card creation for organization: %s", orgName)

	// DB接続
	db, err := database.ConnectWithConnectorIAMAuthN()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// リポジトリとジェネレーターの初期化
	cardRepo := repository.NewCardRepository(db)
	identiconGen := identicon.NewGenerator()
	githubClient := github.NewClient(token)

	// Organization メンバー一覧を取得
	members, err := fetchOrgMembers(token, orgName)
	if err != nil {
		log.Fatalf("Failed to fetch organization members: %v", err)
	}

	log.Printf("Found %d members in organization %s", len(members), orgName)

	// 全メンバーのログイン名を収集してMostUsedLanguageを一括取得
	logins := make([]string, len(members))
	for i, member := range members {
		logins[i] = member.Login
	}

	log.Printf("Fetching most used languages for %d members...", len(logins))
	languageMap, err := githubClient.GetMostUsedLanguages(context.Background(), logins)
	if err != nil {
		log.Printf("Warning: Failed to fetch languages: %v (continuing with empty languages)", err)
		languageMap = make(map[string]github.LanguageInfo)
	}

	// 各メンバーのカードを作成
	created := 0
	skipped := 0
	failed := 0

	for _, member := range members {
		githubID := strconv.FormatInt(member.ID, 10)

		// 既存カードの確認
		_, err := cardRepo.FindByGitHubID(githubID)
		if err == nil {
			// カードが既に存在する
			log.Printf("Skipping %s (ID: %s): card already exists", member.Login, githubID)
			skipped++
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 予期しないエラー
			log.Printf("Error checking card for %s (ID: %s): %v", member.Login, githubID, err)
			failed++
			continue
		}

		// Identicon生成
		color, blocks, err := identiconGen.Generate(githubID)
		if err != nil {
			log.Printf("Failed to generate identicon for %s (ID: %s): %v", member.Login, githubID, err)
			failed++
			continue
		}

		// MostUsedLanguageを取得
		langInfo := languageMap[member.Login]
		mostUsedLanguage := domain.Language{
			LanguageName: langInfo.Name,
			Color:        langInfo.Color,
		}

		// カード作成
		card := domain.NewCard(githubID, member.NodeID, color, blocks, mostUsedLanguage, member.Login, member.FullName, member.AvatarURL)
		if err := cardRepo.Create(card); err != nil {
			log.Printf("Failed to create card for %s (ID: %s): %v", member.Login, githubID, err)
			failed++
			continue
		}

		log.Printf("Created card for %s (ID: %s)", member.Login, githubID)
		created++
	}

	log.Printf("Completed: %d created, %d skipped, %d failed", created, skipped, failed)
}

// fetchOrgMembers fetches all members from a GitHub organization with pagination
func fetchOrgMembers(token, orgName string) ([]GitHubMember, error) {
	var allMembers []GitHubMember
	page := 1
	perPage := 100

	client := &http.Client{}

	for {
		url := fmt.Sprintf("https://api.github.com/orgs/%s/members?per_page=%d&page=%d", orgName, perPage, page)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch members: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
		}

		var members []GitHubMember
		if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		allMembers = append(allMembers, members...)

		// ページネーション: 取得数がperPage未満なら最後のページ
		if len(members) < perPage {
			break
		}

		page++
	}

	return allMembers, nil
}
