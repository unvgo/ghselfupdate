package ghselfupdate

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/go-github/v30/github"
)

// Updater is responsible for managing the context of self-update.
// It contains GitHub client and its context.
type Updater struct {
	api       *github.Client
	apiCtx    context.Context
	validator Validator
	filters   []*regexp.Regexp
}

// Config represents the configuration of self-update.
type Config struct {
	// EnterpriseBaseURL is a base URL of GitHub API. If you want to use this library with GitHub Enterprise,
	// please set "https://{your-organization-address}/api/v3/" to this field.
	EnterpriseBaseURL string
	// EnterpriseUploadURL is a URL to upload stuffs to GitHub Enterprise instance. This is often the same as an API base URL.
	// So if this field is not set and EnterpriseBaseURL is set, EnterpriseBaseURL is also set to this field.
	EnterpriseUploadURL string
	// Validator represents types which enable additional validation of downloaded release.
	Validator Validator
	// Filters are regexp used to filter on specific assets for releases with multiple assets.
	// An asset is selected if it matches any of those, in addition to the regular tag, os, arch, extensions.
	// Please make sure that your filter(s) uniquely match an asset.
	Filters []string
	// HTTPClient allows to set a custom http.Client for GitHub API requests. If nil, a default client will be used.
	HTTPClient *http.Client
}

// NewUpdater creates a new updater instance. It initializes GitHub API client.
// If you set your API token to $GITHUB_TOKEN, the client will use it.
func NewUpdater(ctx context.Context, config Config) (*Updater, error) {
	filtersRe := make([]*regexp.Regexp, 0, len(config.Filters))
	for _, filter := range config.Filters {
		re, err := regexp.Compile(filter)
		if err != nil {
			return nil, fmt.Errorf("could not compile regular expression %q for filtering releases: %v", filter, err)
		}
		filtersRe = append(filtersRe, re)
	}

	if config.EnterpriseBaseURL == "" {
		client := github.NewClient(config.HTTPClient)
		return &Updater{api: client, apiCtx: ctx, validator: config.Validator, filters: filtersRe}, nil
	}

	u := config.EnterpriseUploadURL
	if u == "" {
		u = config.EnterpriseBaseURL
	}
	client, err := github.NewEnterpriseClient(config.EnterpriseBaseURL, u, config.HTTPClient)
	if err != nil {
		return nil, err
	}
	return &Updater{api: client, apiCtx: ctx, validator: config.Validator, filters: filtersRe}, nil
}

// DefaultUpdater creates a new updater instance with default configuration.
// It initializes GitHub API client with default API base URL.
// If you set your API token to $GITHUB_TOKEN, the client will use it.
func DefaultUpdater() *Updater {
	ctx := context.Background()
	return &Updater{api: github.NewClient(http.DefaultClient), apiCtx: ctx}
}
