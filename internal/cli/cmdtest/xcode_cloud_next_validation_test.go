package cmdtest

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

func runXcodeCloudInvalidNextURLCases(
	t *testing.T,
	argsPrefix []string,
	wantErrPrefix string,
) {
	t.Helper()

	tests := []struct {
		name    string
		next    string
		wantErr string
	}{
		{
			name:    "invalid scheme",
			next:    "http://api.appstoreconnect.apple.com/v1/ciBuildActions/action-1/artifacts?cursor=AQ",
			wantErr: wantErrPrefix + "--next must be an App Store Connect URL",
		},
		{
			name:    "malformed URL",
			next:    "https://api.appstoreconnect.apple.com/%zz",
			wantErr: wantErrPrefix + "--next must be a valid URL:",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := append(append([]string{}, argsPrefix...), "--next", test.next)

			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			var runErr error
			stdout, stderr := captureOutput(t, func() {
				if err := root.Parse(args); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				runErr = root.Run(context.Background())
			})

			if runErr == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(runErr.Error(), test.wantErr) {
				t.Fatalf("expected error %q, got %v", test.wantErr, runErr)
			}
			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
			if stderr != "" {
				t.Fatalf("expected empty stderr, got %q", stderr)
			}
		})
	}
}

func runXcodeCloudPaginateFromNext(
	t *testing.T,
	argsPrefix []string,
	firstURL string,
	secondURL string,
	firstBody string,
	secondBody string,
	wantIDs ...string,
) {
	t.Helper()

	setupAuth(t)
	t.Setenv("ASC_CONFIG_PATH", filepath.Join(t.TempDir(), "nonexistent.json"))

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	requestCount := 0
	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requestCount++
		switch requestCount {
		case 1:
			if req.Method != http.MethodGet || req.URL.String() != firstURL {
				t.Fatalf("unexpected first request: %s %s", req.Method, req.URL.String())
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(firstBody)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		case 2:
			if req.Method != http.MethodGet || req.URL.String() != secondURL {
				t.Fatalf("unexpected second request: %s %s", req.Method, req.URL.String())
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(secondBody)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		default:
			t.Fatalf("unexpected extra request: %s %s", req.Method, req.URL.String())
			return nil, nil
		}
	})

	args := append(append([]string{}, argsPrefix...), "--paginate", "--next", firstURL)

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse(args); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	for _, id := range wantIDs {
		needle := `"id":"` + id + `"`
		if !strings.Contains(stdout, needle) {
			t.Fatalf("expected output to contain %q, got %q", needle, stdout)
		}
	}
}

func TestXcodeCloudBuildRunsListRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "build-runs", "list"},
		"xcode-cloud build-runs: ",
	)
}

func TestXcodeCloudBuildRunsListPaginateFromNextWithoutWorkflowID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciWorkflows/workflow-1/buildRuns?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciWorkflows/workflow-1/buildRuns?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciBuildRuns","id":"ci-run-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciBuildRuns","id":"ci-run-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "build-runs", "list"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-run-next-1",
		"ci-run-next-2",
	)
}

func TestXcodeCloudBuildRunsBuildsRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "build-runs", "builds"},
		"xcode-cloud build-runs builds: ",
	)
}

func TestXcodeCloudBuildRunsBuildsPaginateFromNextWithoutRunID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciBuildRuns/run-1/builds?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciBuildRuns/run-1/builds?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"builds","id":"ci-run-build-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"builds","id":"ci-run-build-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "build-runs", "builds"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-run-build-next-1",
		"ci-run-build-next-2",
	)
}

func TestXcodeCloudIssuesListRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "issues", "list"},
		"xcode-cloud issues list: ",
	)
}

func TestXcodeCloudIssuesListPaginateFromNextWithoutActionID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciBuildActions/action-1/issues?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciBuildActions/action-1/issues?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciIssues","id":"ci-issue-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciIssues","id":"ci-issue-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "issues", "list"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-issue-next-1",
		"ci-issue-next-2",
	)
}

func TestXcodeCloudTestResultsListRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "test-results", "list"},
		"xcode-cloud test-results list: ",
	)
}

func TestXcodeCloudTestResultsListPaginateFromNextWithoutActionID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciBuildActions/action-1/testResults?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciBuildActions/action-1/testResults?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciTestResults","id":"ci-test-result-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciTestResults","id":"ci-test-result-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "test-results", "list"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-test-result-next-1",
		"ci-test-result-next-2",
	)
}

func TestXcodeCloudArtifactsListRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "artifacts", "list"},
		"xcode-cloud artifacts list: ",
	)
}

func TestXcodeCloudArtifactsListPaginateFromNextWithoutActionID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciBuildActions/action-1/artifacts?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciBuildActions/action-1/artifacts?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciArtifacts","id":"ci-artifact-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciArtifacts","id":"ci-artifact-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "artifacts", "list"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-artifact-next-1",
		"ci-artifact-next-2",
	)
}

func TestXcodeCloudProductsListRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "products", "list"},
		"xcode-cloud products: ",
	)
}

func TestXcodeCloudProductsListPaginateFromNext(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciProducts?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciProducts?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciProducts","id":"ci-product-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciProducts","id":"ci-product-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "products", "list"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-product-next-1",
		"ci-product-next-2",
	)
}

func TestXcodeCloudProductsBuildRunsRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "products", "build-runs"},
		"xcode-cloud products build-runs: ",
	)
}

func TestXcodeCloudProductsBuildRunsPaginateFromNextWithoutID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciProducts/prod-1/buildRuns?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciProducts/prod-1/buildRuns?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciBuildRuns","id":"ci-product-run-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciBuildRuns","id":"ci-product-run-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "products", "build-runs"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-product-run-next-1",
		"ci-product-run-next-2",
	)
}

func TestXcodeCloudProductsWorkflowsRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "products", "workflows"},
		"xcode-cloud products workflows: ",
	)
}

func TestXcodeCloudProductsWorkflowsPaginateFromNextWithoutID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciProducts/prod-1/workflows?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciProducts/prod-1/workflows?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciWorkflows","id":"ci-product-workflow-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciWorkflows","id":"ci-product-workflow-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "products", "workflows"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-product-workflow-next-1",
		"ci-product-workflow-next-2",
	)
}

func TestXcodeCloudProductsPrimaryRepositoriesRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "products", "primary-repositories"},
		"xcode-cloud products primary-repositories: ",
	)
}

func TestXcodeCloudProductsPrimaryRepositoriesPaginateFromNextWithoutID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciProducts/prod-1/primaryRepositories?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciProducts/prod-1/primaryRepositories?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"scmRepositories","id":"ci-product-primary-repo-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"scmRepositories","id":"ci-product-primary-repo-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "products", "primary-repositories"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-product-primary-repo-next-1",
		"ci-product-primary-repo-next-2",
	)
}

func TestXcodeCloudProductsAdditionalRepositoriesRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "products", "additional-repositories"},
		"xcode-cloud products additional-repositories: ",
	)
}

func TestXcodeCloudProductsAdditionalRepositoriesPaginateFromNextWithoutID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciProducts/prod-1/additionalRepositories?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciProducts/prod-1/additionalRepositories?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"scmRepositories","id":"ci-product-additional-repo-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"scmRepositories","id":"ci-product-additional-repo-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "products", "additional-repositories"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-product-additional-repo-next-1",
		"ci-product-additional-repo-next-2",
	)
}

func TestXcodeCloudMacOSVersionsListRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "macos-versions", "list"},
		"xcode-cloud macos-versions: ",
	)
}

func TestXcodeCloudMacOSVersionsListPaginateFromNext(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciMacOsVersions?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciMacOsVersions?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciMacOsVersions","id":"ci-macos-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciMacOsVersions","id":"ci-macos-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "macos-versions", "list"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-macos-next-1",
		"ci-macos-next-2",
	)
}

func TestXcodeCloudXcodeVersionsListRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "xcode-versions", "list"},
		"xcode-cloud xcode-versions: ",
	)
}

func TestXcodeCloudXcodeVersionsListPaginateFromNext(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciXcodeVersions?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciXcodeVersions?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciXcodeVersions","id":"ci-xcode-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciXcodeVersions","id":"ci-xcode-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "xcode-versions", "list"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-xcode-next-1",
		"ci-xcode-next-2",
	)
}

func TestXcodeCloudMacOSVersionsXcodeVersionsRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "macos-versions", "xcode-versions"},
		"xcode-cloud macos-versions xcode-versions: ",
	)
}

func TestXcodeCloudMacOSVersionsXcodeVersionsPaginateFromNextWithoutID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciMacOsVersions/macos-1/xcodeVersions?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciMacOsVersions/macos-1/xcodeVersions?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciXcodeVersions","id":"ci-macos-xcode-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciXcodeVersions","id":"ci-macos-xcode-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "macos-versions", "xcode-versions"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-macos-xcode-next-1",
		"ci-macos-xcode-next-2",
	)
}

func TestXcodeCloudXcodeVersionsMacOSVersionsRejectsInvalidNextURL(t *testing.T) {
	runXcodeCloudInvalidNextURLCases(
		t,
		[]string{"xcode-cloud", "xcode-versions", "macos-versions"},
		"xcode-cloud xcode-versions macos-versions: ",
	)
}

func TestXcodeCloudXcodeVersionsMacOSVersionsPaginateFromNextWithoutID(t *testing.T) {
	const firstURL = "https://api.appstoreconnect.apple.com/v1/ciXcodeVersions/xcode-1/macOsVersions?cursor=AQ&limit=200"
	const secondURL = "https://api.appstoreconnect.apple.com/v1/ciXcodeVersions/xcode-1/macOsVersions?cursor=BQ&limit=200"

	firstBody := `{"data":[{"type":"ciMacOsVersions","id":"ci-xcode-macos-next-1"}],"links":{"next":"` + secondURL + `"}}`
	secondBody := `{"data":[{"type":"ciMacOsVersions","id":"ci-xcode-macos-next-2"}],"links":{"next":""}}`

	runXcodeCloudPaginateFromNext(
		t,
		[]string{"xcode-cloud", "xcode-versions", "macos-versions"},
		firstURL,
		secondURL,
		firstBody,
		secondBody,
		"ci-xcode-macos-next-1",
		"ci-xcode-macos-next-2",
	)
}
