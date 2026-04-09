package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/grevus/mcp-jira/internal/handlers"
	"github.com/grevus/mcp-jira/internal/jira"
	"github.com/stretchr/testify/require"
)

type fakeIssueFetcher struct {
	issue jira.Issue
	desc  string
	err   error
}

func (f *fakeIssueFetcher) GetIssue(_ context.Context, _ string) (jira.Issue, string, error) {
	return f.issue, f.desc, f.err
}

type fakeRetriever struct {
	gotQuery string
	gotTopK  int
	hits     []handlers.Hit
	err      error
}

func (f *fakeRetriever) Search(_ context.Context, _ string, query string, topK int) ([]handlers.Hit, error) {
	f.gotQuery = query
	f.gotTopK = topK
	return f.hits, f.err
}

func TestSimilarIssues_FiltersSelfMatch(t *testing.T) {
	fetcher := &fakeIssueFetcher{
		issue: jira.Issue{Key: "ABC-10", Summary: "Login broken", Status: "Open"},
		desc:  "users cannot login after deploy",
	}
	ret := &fakeRetriever{
		hits: []handlers.Hit{
			{IssueKey: "ABC-10", Summary: "Login broken", Score: 0.99},
			{IssueKey: "ABC-11", Summary: "SSO timeout", Score: 0.87},
			{IssueKey: "ABC-12", Summary: "Auth 500", Score: 0.80},
		},
	}

	h := handlers.SimilarIssues(fetcher, ret)
	out, err := h(context.Background(), handlers.SimilarIssuesInput{
		ProjectKey: "ABC",
		IssueKey:   "ABC-10",
		TopK:       2,
	})
	require.NoError(t, err)
	require.Equal(t, "ABC-10", out.Source.Key)
	require.Len(t, out.SimilarIssues, 2)
	require.Equal(t, "ABC-11", out.SimilarIssues[0].IssueKey)
	require.Equal(t, "ABC-12", out.SimilarIssues[1].IssueKey)
	require.Contains(t, ret.gotQuery, "Login broken")
	require.Contains(t, ret.gotQuery, "users cannot login")
	require.Equal(t, 3, ret.gotTopK) // topK+1 для self-match buffer
}

func TestSimilarIssues_RejectsBadInput(t *testing.T) {
	h := handlers.SimilarIssues(&fakeIssueFetcher{}, &fakeRetriever{})
	_, err := h(context.Background(), handlers.SimilarIssuesInput{ProjectKey: "ABC"})
	require.Error(t, err)
	_, err = h(context.Background(), handlers.SimilarIssuesInput{IssueKey: "ABC-1"})
	require.Error(t, err)
	_, err = h(context.Background(), handlers.SimilarIssuesInput{ProjectKey: "ABC", IssueKey: "ABC-1", TopK: 100})
	require.Error(t, err)
}

func TestSimilarIssues_PropagatesErrors(t *testing.T) {
	h := handlers.SimilarIssues(&fakeIssueFetcher{err: errors.New("boom")}, &fakeRetriever{})
	_, err := h(context.Background(), handlers.SimilarIssuesInput{ProjectKey: "ABC", IssueKey: "ABC-1"})
	require.Error(t, err)
}
