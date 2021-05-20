package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

const (
	owner            = "mesosphere"
	repo             = "kubernetes-base-addons"
	addonLabelPrefix = "addon/"
)

var ctx = context.Background()

func getMilestone(title string, client *github.Client) (*github.Milestone, error) {
	milestones, _, err := client.Issues.ListMilestones(ctx, owner, repo, &github.MilestoneListOptions{})
	if err != nil {
		return nil, err
	}
	for _, milestone := range milestones {
		if *milestone.Title == title {
			if milestone.GetOpenIssues() != 0 {
				return nil, fmt.Errorf("milestone, %s, still has %d open issues/PRs", title, milestone.GetOpenIssues())
			}
			return milestone, nil
		}
	}
	return nil, errors.New("milestone, " + title + ", not found.")
}

func getPullRequestIssuesInMilestone(title string, client *github.Client) (map[int64]*github.Issue, error) {
	milestone, err := getMilestone(title, client)
	if err != nil {
		return nil, err
	}
	result := map[int64]*github.Issue{}
	options := github.IssueListByRepoOptions{Milestone: strconv.Itoa(milestone.GetNumber()), State: "closed"}
	options.Page = 0
	options.PerPage = 100
	done := false
	for pullRequests, response, err := client.Issues.ListByRepo(ctx, owner, repo, &options); !done; {
		if err != nil {
			return nil, err
		}
		for _, pullRequest := range pullRequests {
			if pullRequest.Milestone.GetID() != milestone.GetID() {
				continue
			}
			if pullRequest.Milestone.GetID() == milestone.GetID() {
				result[pullRequest.GetID()] = pullRequest
			}
		}
		if options.Page == response.LastPage {
			done = true
		}
		options.Page = response.NextPage
	}
	return result, nil
}

func getLabelsFromIssues(pullRequests map[int64]*github.Issue) []string {
	labels := map[string]*string{}
	for _, pr := range pullRequests {
		for _, label := range pr.Labels {
			labels[label.GetName()] = nil
		}
	}
	result := []string{}
	for key := range labels {
		result = append(result, key)
	}
	return result
}

func labelInSlice(a string, list []*github.Label) bool {
	for _, b := range list {
		if b.GetName() == a {
			return true
		}
	}
	return false
}

func releaseNoteForIssue(issue *github.Issue) (string, error) {
	if issue == nil {
		return "", nil
	}
	matchString := "(?s)```release-note(.*?)```"
	re := regexp.MustCompile(matchString)
	m := re.FindStringSubmatch(issue.GetBody())
	if len(m) == 0 {
		return "", fmt.Errorf("PR #%d is missing its release-note: %s", issue.GetNumber(), issue.GetHTMLURL())
	}
	m1 := strings.Split(m[1], "\n")
	out := []string{}
	for _, s := range m1 {
		o := strings.Trim(s, "\n\r- \t")
		if o != "" && strings.ToUpper(o) != "NONE" {
			out = append(out, strings.Trim(s, "\n\r- \t"))
		}
	}
	if len(out) == 0 {
		return "", nil
	}
	return fmt.Sprintf("  - %s\n  #%d (@%s)\n\n", strings.Join(out, "\n  - "), issue.GetNumber(), issue.GetUser().GetLogin()), nil
}

func buildReleaseNoteForLabel(label string, issues map[int64]*github.Issue) (string, error) {
	result := fmt.Sprintf("### %s\n", label[6:])
	notes := ""
	for _, i := range issues {
		if i == nil || !labelInSlice(label, i.Labels) {
			continue
		}
		note, err := releaseNoteForIssue(i)
		if err != nil {
			return "", err
		}
		notes = notes + note
	}
	if notes == "" {
		return "", nil
	}
	result = result + notes
	return result, nil
}

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	milestone := os.Getenv("KBA_MILESTONE")
	if milestone == "" {
		log.Fatal("must set KBA_MILESTONE")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	a, err := getPullRequestIssuesInMilestone(milestone, client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	labels := getLabelsFromIssues(a)
	sort.Strings(labels)
	releaseNotes := ""
	for _, label := range labels {
		if !strings.HasPrefix(label, addonLabelPrefix) {
			continue
		}
		note, err := buildReleaseNoteForLabel(label, a)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		releaseNotes = releaseNotes + note
	}
	fmt.Println(releaseNotes)
}
