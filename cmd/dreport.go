/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"time"

	"github.com/aoyouer/message-link/collector"
	"github.com/aoyouer/message-link/messenger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// dreportCmd represents the dreport command
var dreportCmd = &cobra.Command{
	Use:   "dreport",
	Short: "发送日报 (PR)",
	Run: func(cmd *cobra.Command, args []string) {
		sendGithubPRToFeishu()
	},
}

func init() {
	rootCmd.AddCommand(dreportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dreportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dreportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 向飞书发送github 24小时pr信息
func sendGithubPRToFeishu() {
	var msg = messenger.FeishuMessage{
		Title:   "过去24小时pr",
		Content: []messenger.FeishuMessageContent{},
	}

	repoList := collector.GetWatchedRepos()

	if repoList == nil {
		zap.L().Error("repoMap is nil")
		return
	}

	for _, repo := range repoList {
		// 获取指定repo的一天内的pr
		prs := collector.GetGithubCollector().ListPR(repo, time.Now().AddDate(0, 0, -1))
		if prs == nil {
			continue
		}
		// 每一个repo的标题部分
		var feishuMessageContent = messenger.FeishuMessageContent{
			messenger.FeishuMessageContentItem{
				Tag:  "text",
				Text: "\n[Repository] ",
			},
			messenger.FeishuMessageContentItem{
				Tag:  "a",
				Text: repo.GetFullName(),
				Href: repo.GetHTMLURL(),
			},
		}
		msg.Content = append(msg.Content, feishuMessageContent)

		// 每一个repo过去一天的pr
		for _, pr := range prs {
			var feishuMessageContent = messenger.FeishuMessageContent{
				messenger.FeishuMessageContentItem{
					Tag:  "text",
					Text: pr.GetTitle() + " ",
				},
				messenger.FeishuMessageContentItem{
					Tag:  "a",
					Text: pr.GetHTMLURL(),
					Href: pr.GetHTMLURL(),
				},
				messenger.FeishuMessageContentItem{
					Tag:  "a",
					Text: " @" + pr.User.GetLogin(),
					Href: pr.User.GetHTMLURL(),
				},
			}
			msg.Content = append(msg.Content, feishuMessageContent)
		}
	}
	messenger.GetFeishuMessenger().SendHyperTextMessage(msg)
}
