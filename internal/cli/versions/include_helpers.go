package versions

import (
	"fmt"
	"strings"
)

func normalizeAppStoreVersionInclude(value string) ([]string, error) {
	return normalizeInclude(value, appStoreVersionIncludeList(), "--include")
}

func normalizeInclude(value string, allowed []string, flagName string) ([]string, error) {
	include := splitCSV(value)
	if len(include) == 0 {
		return nil, nil
	}
	allowedMap := map[string]struct{}{}
	for _, option := range allowed {
		allowedMap[option] = struct{}{}
	}
	for _, option := range include {
		if _, ok := allowedMap[option]; !ok {
			return nil, fmt.Errorf("%s must be one of: %s", flagName, strings.Join(allowed, ", "))
		}
	}
	return include, nil
}

func appStoreVersionIncludeList() []string {
	return []string{
		"ageRatingDeclaration",
		"appStoreReviewDetail",
		"appClipDefaultExperience",
		"appStoreVersionExperiments",
		"appStoreVersionExperimentsV2",
		"appStoreVersionSubmission",
		"customerReviews",
		"routingAppCoverage",
		"alternativeDistributionPackage",
		"gameCenterAppVersion",
	}
}
