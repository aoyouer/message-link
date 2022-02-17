package collector

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type GithubCollector struct {
	githubClient *github.Client
	repos        []*github.Repository // 以 owner(org):repo 的形式记录 不对个人仓库以及组织仓库做区分
}

//github消息收集
var (
	githubCollector *GithubCollector
)

func initGithubClient() {
	githubCollector = new(GithubCollector)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetViper().GetString("github.token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubCollector.githubClient = github.NewClient(tc)

	// 查看所有关注的repo
	repoUrls := viper.GetStringSlice("github.repos")
	for _, r := range repoUrls {
		paths := strings.Split(r, "/")
		repo, owner := paths[len(paths)-1], paths[len(paths)-2]
		if gr, _, err := githubCollector.githubClient.Repositories.Get(ctx, owner, repo); err != nil {
			zap.L().Error("failed to get repo: " + err.Error())
		} else {
			zap.S().Infof("load repo: %s owned by: %s", gr.GetName(), gr.Owner.GetLogin())
			githubCollector.repos = append(githubCollector.repos, gr)
		}
	}
}

func GetGithubCollector() *GithubCollector {
	return githubCollector
}

func GetWatchedRepos() []*github.Repository {
	return githubCollector.repos
}

// 获取的pr，传入起始时间
func (gh *GithubCollector) ListPR(repo *github.Repository, since time.Time) []*github.PullRequest {
	ctx := context.Background()
	var prList []*github.PullRequest
	if repo == nil {
		zap.L().Error("repo is nil")
		return nil
	}
	prs, _, err := gh.githubClient.PullRequests.List(ctx, repo.GetOwner().GetLogin(), repo.GetName(), &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 60,
		},
	})

	if err != nil {
		zap.L().Error(err.Error())
		return nil
	} else {
		for _, pr := range prs {
			if since.Before(pr.GetCreatedAt()) {
				prList = append(prList, pr)
			}
		}
	}
	return prList
}
