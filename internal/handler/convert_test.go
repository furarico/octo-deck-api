package handler

import (
	"testing"
	"time"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

// CardをAPIのCard型に変換するテスト
func TestConvertCardToAPI(t *testing.T) {
	tests := []struct {
		name string
		card domain.Card
		want api.Card
	}{
		{
			name: "正常にCardを変換できる",
			card: domain.Card{
				ID:       domain.NewCardID(),
				GithubID: "12345",
				UserName: "testuser",
				FullName: "Test User",
				IconUrl:  "https://example.com/avatar.png",
				Color:    domain.Color("#abcdef"),
				Blocks: domain.Blocks{
					{true, false, true, false, true},
					{false, true, false, true, false},
					{true, true, true, true, true},
					{false, false, false, false, false},
					{true, false, false, false, true},
				},
				MostUsedLanguage: domain.Language{
					LanguageName: "Go",
					Color:        "#00ADD8",
				},
			},
			want: api.Card{
				GithubId: "12345",
				UserName: "testuser",
				FullName: "Test User",
				IconUrl:  "https://example.com/avatar.png",
				Identicon: api.Identicon{
					Blocks: [][]bool{
						{true, false, true, false, true},
						{false, true, false, true, false},
						{true, true, true, true, true},
						{false, false, false, false, false},
						{true, false, false, false, true},
					},
					Color: "#abcdef",
				},
				MostUsedLanguage: api.Language{
					Name:  "Go",
					Color: "#00ADD8",
				},
			},
		},
		{
			name: "空の値でも変換できる",
			card: domain.Card{
				ID:       domain.NewCardID(),
				GithubID: "",
				UserName: "",
				FullName: "",
				IconUrl:  "",
				Color:    domain.Color(""),
				Blocks:   domain.Blocks{},
				MostUsedLanguage: domain.Language{
					LanguageName: "",
					Color:        "",
				},
			},
			want: api.Card{
				GithubId: "",
				UserName: "",
				FullName: "",
				IconUrl:  "",
				Identicon: api.Identicon{
					Blocks: [][]bool{
						{false, false, false, false, false},
						{false, false, false, false, false},
						{false, false, false, false, false},
						{false, false, false, false, false},
						{false, false, false, false, false},
					},
					Color: "",
				},
				MostUsedLanguage: api.Language{
					Name:  "",
					Color: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertCardToAPI(tt.card)
			if got.GithubId != tt.want.GithubId {
				t.Errorf("GithubId = %v, want %v", got.GithubId, tt.want.GithubId)
			}
			if got.UserName != tt.want.UserName {
				t.Errorf("UserName = %v, want %v", got.UserName, tt.want.UserName)
			}
			if got.FullName != tt.want.FullName {
				t.Errorf("FullName = %v, want %v", got.FullName, tt.want.FullName)
			}
			if got.IconUrl != tt.want.IconUrl {
				t.Errorf("IconUrl = %v, want %v", got.IconUrl, tt.want.IconUrl)
			}
			if got.Identicon.Color != tt.want.Identicon.Color {
				t.Errorf("Identicon.Color = %v, want %v", got.Identicon.Color, tt.want.Identicon.Color)
			}
			if len(got.Identicon.Blocks) != len(tt.want.Identicon.Blocks) {
				t.Errorf("Identicon.Blocks length = %v, want %v", len(got.Identicon.Blocks), len(tt.want.Identicon.Blocks))
			}
			for i := range got.Identicon.Blocks {
				if len(got.Identicon.Blocks[i]) != len(tt.want.Identicon.Blocks[i]) {
					t.Errorf("Identicon.Blocks[%d] length = %v, want %v", i, len(got.Identicon.Blocks[i]), len(tt.want.Identicon.Blocks[i]))
				}
				for j := range got.Identicon.Blocks[i] {
					if got.Identicon.Blocks[i][j] != tt.want.Identicon.Blocks[i][j] {
						t.Errorf("Identicon.Blocks[%d][%d] = %v, want %v", i, j, got.Identicon.Blocks[i][j], tt.want.Identicon.Blocks[i][j])
					}
				}
			}
			if got.MostUsedLanguage.Name != tt.want.MostUsedLanguage.Name {
				t.Errorf("MostUsedLanguage.Name = %v, want %v", got.MostUsedLanguage.Name, tt.want.MostUsedLanguage.Name)
			}
			if got.MostUsedLanguage.Color != tt.want.MostUsedLanguage.Color {
				t.Errorf("MostUsedLanguage.Color = %v, want %v", got.MostUsedLanguage.Color, tt.want.MostUsedLanguage.Color)
			}
		})
	}
}

// BlocksをAPIのBlocks型に変換するテスト
func TestConvertBlocks(t *testing.T) {
	tests := []struct {
		name   string
		blocks domain.Blocks
		want   [][]bool
	}{
		{
			name: "正常にBlocksを変換できる",
			blocks: domain.Blocks{
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, true, true, true, true},
				{false, false, false, false, false},
				{true, false, false, false, true},
			},
			want: [][]bool{
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, true, true, true, true},
				{false, false, false, false, false},
				{true, false, false, false, true},
			},
		},
		{
			name:   "空のBlocksを変換できる",
			blocks: domain.Blocks{},
			want: [][]bool{
				{false, false, false, false, false},
				{false, false, false, false, false},
				{false, false, false, false, false},
				{false, false, false, false, false},
				{false, false, false, false, false},
			},
		},
		{
			name: "すべてtrueのBlocksを変換できる",
			blocks: domain.Blocks{
				{true, true, true, true, true},
				{true, true, true, true, true},
				{true, true, true, true, true},
				{true, true, true, true, true},
				{true, true, true, true, true},
			},
			want: [][]bool{
				{true, true, true, true, true},
				{true, true, true, true, true},
				{true, true, true, true, true},
				{true, true, true, true, true},
				{true, true, true, true, true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertBlocks(tt.blocks)
			if len(got) != len(tt.want) {
				t.Errorf("Blocks length = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if len(got[i]) != len(tt.want[i]) {
					t.Errorf("Blocks[%d] length = %v, want %v", i, len(got[i]), len(tt.want[i]))
				}
				for j := range got[i] {
					if got[i][j] != tt.want[i][j] {
						t.Errorf("Blocks[%d][%d] = %v, want %v", i, j, got[i][j], tt.want[i][j])
					}
				}
			}
		})
	}
}

// CommunityをAPIのCommunity型に変換するテスト
func TestConvertCommunityToAPI(t *testing.T) {
	tests := []struct {
		name      string
		community domain.Community
		want      api.Community
	}{
		{
			name: "正常にCommunityを変換できる",
			community: domain.Community{
				ID:   domain.NewCommunityID(),
				Name: "Test Community",
			},
			want: api.Community{
				Id:   "", // UUIDは実行時に変わるので、空文字列で比較しない
				Name: "Test Community",
			},
		},
		{
			name: "空の名前でも変換できる",
			community: domain.Community{
				ID:   domain.NewCommunityID(),
				Name: "",
			},
			want: api.Community{
				Id:   "", // UUIDは実行時に変わるので、空文字列で比較しない
				Name: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertCommunityToAPI(tt.community)
			// UUIDの変換を確認
			gotUUID, err := uuid.Parse(got.Id)
			if err != nil {
				t.Errorf("Id is not a valid UUID: %v", err)
			}
			expectedUUID := uuid.UUID(tt.community.ID)
			if gotUUID != expectedUUID {
				t.Errorf("Id = %v, want %v", gotUUID, expectedUUID)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
		})
	}
}

// UserStatsをAPIのUserStats型に変換するテスト
func TestConvertUserStatsToAPI(t *testing.T) {
	tests := []struct {
		name  string
		stats *domain.Stats
		want  api.UserStats
	}{
		{
			name: "正常にUserStatsを変換できる",
			stats: &domain.Stats{
				Contributions: []domain.Contribution{
					{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Count: 10},
					{Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Count: 5},
					{Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), Count: 15},
				},
				TotalContribution: 30,
				ContributionDetail: domain.ContributionDetail{
					ReviewCount:      5,
					CommitCount:      10,
					IssueCount:       3,
					PullRequestCount: 2,
				},
				MostUsedLanguage: domain.Language{
					LanguageName: "Go",
					Color:        "#00ADD8",
				},
			},
			want: api.UserStats{
				Contributions: []api.Contribution{
					{Date: types.Date{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}, Count: 10},
					{Date: types.Date{Time: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)}, Count: 5},
					{Date: types.Date{Time: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)}, Count: 15},
				},
				TotalContribution: 30,
				ContributionDetail: api.ContributionDetail{
					ReviewCount:      5,
					CommitCount:      10,
					IssueCount:       3,
					PullRequestCount: 2,
				},
				MostUsedLanguage: api.Language{
					Name:  "Go",
					Color: "#00ADD8",
				},
			},
		},
		{
			name: "空のContributionsでも変換できる",
			stats: &domain.Stats{
				Contributions:      []domain.Contribution{},
				TotalContribution:  0,
				ContributionDetail: domain.ContributionDetail{},
				MostUsedLanguage: domain.Language{
					LanguageName: "",
					Color:        "",
				},
			},
			want: api.UserStats{
				Contributions:      []api.Contribution{},
				TotalContribution:  0,
				ContributionDetail: api.ContributionDetail{},
				MostUsedLanguage: api.Language{
					Name:  "",
					Color: "",
				},
			},
		},
		{
			name: "大きな値でも変換できる",
			stats: &domain.Stats{
				Contributions: []domain.Contribution{
					{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Count: 999999},
				},
				TotalContribution: 999999,
				ContributionDetail: domain.ContributionDetail{
					ReviewCount:      999999,
					CommitCount:      999999,
					IssueCount:       999999,
					PullRequestCount: 999999,
				},
				MostUsedLanguage: domain.Language{
					LanguageName: "Go",
					Color:        "#00ADD8",
				},
			},
			want: api.UserStats{
				Contributions: []api.Contribution{
					{Date: types.Date{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}, Count: 999999},
				},
				TotalContribution: 999999,
				ContributionDetail: api.ContributionDetail{
					ReviewCount:      999999,
					CommitCount:      999999,
					IssueCount:       999999,
					PullRequestCount: 999999,
				},
				MostUsedLanguage: api.Language{
					Name:  "Go",
					Color: "#00ADD8",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertUserStatsToAPI(tt.stats)
			if err != nil {
				t.Errorf("convertUserStatsToAPI() error = %v", err)
				return
			}
			if got.TotalContribution != tt.want.TotalContribution {
				t.Errorf("TotalContribution = %v, want %v", got.TotalContribution, tt.want.TotalContribution)
			}
			if len(got.Contributions) != len(tt.want.Contributions) {
				t.Errorf("Contributions length = %v, want %v", len(got.Contributions), len(tt.want.Contributions))
			}
			for i := range got.Contributions {
				if got.Contributions[i].Count != tt.want.Contributions[i].Count {
					t.Errorf("Contributions[%d].Count = %v, want %v", i, got.Contributions[i].Count, tt.want.Contributions[i].Count)
				}
				if !got.Contributions[i].Date.Time.Equal(tt.want.Contributions[i].Date.Time) {
					t.Errorf("Contributions[%d].Date = %v, want %v", i, got.Contributions[i].Date.Time, tt.want.Contributions[i].Date.Time)
				}
			}
			if got.ContributionDetail.ReviewCount != tt.want.ContributionDetail.ReviewCount {
				t.Errorf("ContributionDetail.ReviewCount = %v, want %v", got.ContributionDetail.ReviewCount, tt.want.ContributionDetail.ReviewCount)
			}
			if got.ContributionDetail.CommitCount != tt.want.ContributionDetail.CommitCount {
				t.Errorf("ContributionDetail.CommitCount = %v, want %v", got.ContributionDetail.CommitCount, tt.want.ContributionDetail.CommitCount)
			}
			if got.ContributionDetail.IssueCount != tt.want.ContributionDetail.IssueCount {
				t.Errorf("ContributionDetail.IssueCount = %v, want %v", got.ContributionDetail.IssueCount, tt.want.ContributionDetail.IssueCount)
			}
			if got.ContributionDetail.PullRequestCount != tt.want.ContributionDetail.PullRequestCount {
				t.Errorf("ContributionDetail.PullRequestCount = %v, want %v", got.ContributionDetail.PullRequestCount, tt.want.ContributionDetail.PullRequestCount)
			}
			if got.MostUsedLanguage.Name != tt.want.MostUsedLanguage.Name {
				t.Errorf("MostUsedLanguage.Name = %v, want %v", got.MostUsedLanguage.Name, tt.want.MostUsedLanguage.Name)
			}
			if got.MostUsedLanguage.Color != tt.want.MostUsedLanguage.Color {
				t.Errorf("MostUsedLanguage.Color = %v, want %v", got.MostUsedLanguage.Color, tt.want.MostUsedLanguage.Color)
			}
		})
	}
}

// HighlightedCardをAPIのHighlightedCard型に変換するテスト
func TestConvertHighlightedCardToAPI(t *testing.T) {
	tests := []struct {
		name            string
		highlightedCard domain.HighlightedCard
		want            api.HighlightedCard
	}{
		{
			name: "正常にHighlightedCardを変換できる",
			highlightedCard: domain.HighlightedCard{
				BestContributor: domain.Card{
					ID:       domain.NewCardID(),
					GithubID: "contributor1",
					UserName: "contributor1",
					FullName: "Contributor One",
					IconUrl:  "https://example.com/contributor1.png",
					Color:    domain.Color("#111111"),
					Blocks:   domain.Blocks{},
					MostUsedLanguage: domain.Language{
						LanguageName: "Go",
						Color:        "#00ADD8",
					},
				},
				BestCommitter: domain.Card{
					ID:       domain.NewCardID(),
					GithubID: "committer1",
					UserName: "committer1",
					FullName: "Committer One",
					IconUrl:  "https://example.com/committer1.png",
					Color:    domain.Color("#222222"),
					Blocks:   domain.Blocks{},
					MostUsedLanguage: domain.Language{
						LanguageName: "Python",
						Color:        "#3776AB",
					},
				},
				BestIssuer: domain.Card{
					ID:       domain.NewCardID(),
					GithubID: "issuer1",
					UserName: "issuer1",
					FullName: "Issuer One",
					IconUrl:  "https://example.com/issuer1.png",
					Color:    domain.Color("#333333"),
					Blocks:   domain.Blocks{},
					MostUsedLanguage: domain.Language{
						LanguageName: "JavaScript",
						Color:        "#F7DF1E",
					},
				},
				BestPullRequester: domain.Card{
					ID:       domain.NewCardID(),
					GithubID: "pr1",
					UserName: "pr1",
					FullName: "PR One",
					IconUrl:  "https://example.com/pr1.png",
					Color:    domain.Color("#444444"),
					Blocks:   domain.Blocks{},
					MostUsedLanguage: domain.Language{
						LanguageName: "TypeScript",
						Color:        "#3178C6",
					},
				},
				BestReviewer: domain.Card{
					ID:       domain.NewCardID(),
					GithubID: "reviewer1",
					UserName: "reviewer1",
					FullName: "Reviewer One",
					IconUrl:  "https://example.com/reviewer1.png",
					Color:    domain.Color("#555555"),
					Blocks:   domain.Blocks{},
					MostUsedLanguage: domain.Language{
						LanguageName: "Rust",
						Color:        "#000000",
					},
				},
			},
			want: api.HighlightedCard{
				BestContributor: api.Card{
					GithubId: "contributor1",
					UserName: "contributor1",
					FullName: "Contributor One",
					IconUrl:  "https://example.com/contributor1.png",
					MostUsedLanguage: api.Language{
						Name:  "Go",
						Color: "#00ADD8",
					},
				},
				BestCommitter: api.Card{
					GithubId: "committer1",
					UserName: "committer1",
					FullName: "Committer One",
					IconUrl:  "https://example.com/committer1.png",
					MostUsedLanguage: api.Language{
						Name:  "Python",
						Color: "#3776AB",
					},
				},
				BestIssuer: api.Card{
					GithubId: "issuer1",
					UserName: "issuer1",
					FullName: "Issuer One",
					IconUrl:  "https://example.com/issuer1.png",
					MostUsedLanguage: api.Language{
						Name:  "JavaScript",
						Color: "#F7DF1E",
					},
				},
				BestPullRequester: api.Card{
					GithubId: "pr1",
					UserName: "pr1",
					FullName: "PR One",
					IconUrl:  "https://example.com/pr1.png",
					MostUsedLanguage: api.Language{
						Name:  "TypeScript",
						Color: "#3178C6",
					},
				},
				BestReviewer: api.Card{
					GithubId: "reviewer1",
					UserName: "reviewer1",
					FullName: "Reviewer One",
					IconUrl:  "https://example.com/reviewer1.png",
					MostUsedLanguage: api.Language{
						Name:  "Rust",
						Color: "#000000",
					},
				},
			},
		},
		{
			name: "空のCardでも変換できる",
			highlightedCard: domain.HighlightedCard{
				BestContributor:   domain.Card{},
				BestCommitter:     domain.Card{},
				BestIssuer:        domain.Card{},
				BestPullRequester: domain.Card{},
				BestReviewer:      domain.Card{},
			},
			want: api.HighlightedCard{
				BestContributor:   api.Card{},
				BestCommitter:     api.Card{},
				BestIssuer:        api.Card{},
				BestPullRequester: api.Card{},
				BestReviewer:      api.Card{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertHighlightedCardToAPI(tt.highlightedCard)
			if got.BestContributor.GithubId != tt.want.BestContributor.GithubId {
				t.Errorf("BestContributor.GithubId = %v, want %v", got.BestContributor.GithubId, tt.want.BestContributor.GithubId)
			}
			if got.BestCommitter.GithubId != tt.want.BestCommitter.GithubId {
				t.Errorf("BestCommitter.GithubId = %v, want %v", got.BestCommitter.GithubId, tt.want.BestCommitter.GithubId)
			}
			if got.BestIssuer.GithubId != tt.want.BestIssuer.GithubId {
				t.Errorf("BestIssuer.GithubId = %v, want %v", got.BestIssuer.GithubId, tt.want.BestIssuer.GithubId)
			}
			if got.BestPullRequester.GithubId != tt.want.BestPullRequester.GithubId {
				t.Errorf("BestPullRequester.GithubId = %v, want %v", got.BestPullRequester.GithubId, tt.want.BestPullRequester.GithubId)
			}
			if got.BestReviewer.GithubId != tt.want.BestReviewer.GithubId {
				t.Errorf("BestReviewer.GithubId = %v, want %v", got.BestReviewer.GithubId, tt.want.BestReviewer.GithubId)
			}
		})
	}
}
