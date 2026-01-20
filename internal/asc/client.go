package asc

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/auth"
)

const (
	// BaseURL is the App Store Connect API base URL
	BaseURL = "https://api.appstoreconnect.apple.com"
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
	tokenLifetime  = 20 * time.Minute
)

// Client is an App Store Connect API client
type Client struct {
	httpClient *http.Client
	keyID      string
	issuerID   string
	privateKey *ecdsa.PrivateKey
}

// Resource is a generic ASC API resource wrapper.
type Resource[T any] struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes T      `json:"attributes"`
}

// Response is a generic ASC API response wrapper.
type Response[T any] struct {
	Data  []Resource[T] `json:"data"`
	Links Links         `json:"links,omitempty"`
}

// FeedbackAttributes describes beta feedback screenshot submissions.
type FeedbackAttributes struct {
	CreatedDate string `json:"createdDate"`
	Comment     string `json:"comment"`
	Email       string `json:"email"`
	DeviceModel string `json:"deviceModel,omitempty"`
	OSVersion   string `json:"osVersion,omitempty"`
}

// CrashAttributes describes beta feedback crash submissions.
type CrashAttributes struct {
	CreatedDate string `json:"createdDate"`
	Comment     string `json:"comment"`
	Email       string `json:"email"`
	DeviceModel string `json:"deviceModel,omitempty"`
	OSVersion   string `json:"osVersion,omitempty"`
	CrashLog    string `json:"crashLog,omitempty"`
}

// ReviewAttributes describes App Store customer reviews.
type ReviewAttributes struct {
	Rating           int    `json:"rating"`
	Title            string `json:"title"`
	Body             string `json:"body"`
	ReviewerNickname string `json:"reviewerNickname"`
	CreatedDate      string `json:"createdDate"`
	Territory        string `json:"territory"`
}

// FeedbackResponse is the response from beta feedback screenshots endpoint.
type FeedbackResponse = Response[FeedbackAttributes]

// CrashesResponse is the response from beta feedback crashes endpoint.
type CrashesResponse = Response[CrashAttributes]

// ReviewsResponse is the response from customer reviews endpoint.
type ReviewsResponse = Response[ReviewAttributes]

type reviewQuery struct {
	rating    int
	territory string
}

// ReviewOption is a functional option for GetReviews.
type ReviewOption func(*reviewQuery)

// WithRating filters reviews by star rating (1-5).
func WithRating(rating int) ReviewOption {
	return func(r *reviewQuery) {
		if rating >= 1 && rating <= 5 {
			r.rating = rating
		}
	}
}

// WithTerritory filters reviews by territory code (e.g. US, GBR).
func WithTerritory(territory string) ReviewOption {
	return func(r *reviewQuery) {
		if territory != "" {
			r.territory = strings.ToUpper(territory)
		}
	}
}

// NewClient creates a new ASC client
func NewClient(keyID, issuerID, privateKeyPath string) (*Client, error) {
	if err := auth.ValidateKeyFile(privateKeyPath); err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	key, err := auth.LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		keyID:      keyID,
		issuerID:   issuerID,
		privateKey: key,
	}, nil
}

// newRequest creates a new HTTP request with JWT authentication
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	// Generate JWT token
	token, err := c.generateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	url := BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// generateJWT generates a JWT for ASC API authentication
func (c *Client) generateJWT() (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    c.issuerID,
		Audience:  jwt.ClaimStrings{"appstoreconnect-v1"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = c.keyID

	// Sign with the private key
	signedToken, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// do performs an HTTP request and returns the response
func (c *Client) do(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		if err := ParseError(respBody); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func buildReviewQuery(opts []ReviewOption) string {
	query := &reviewQuery{}
	for _, opt := range opts {
		opt(query)
	}

	values := url.Values{}
	if query.territory != "" {
		values.Set("filter[territory]", query.territory)
	}
	if query.rating >= 1 && query.rating <= 5 {
		values.Set("filter[rating]", fmt.Sprintf("%d", query.rating))
	}

	return values.Encode()
}

// GetFeedback retrieves TestFlight feedback
func (c *Client) GetFeedback(ctx context.Context, appID string) (*FeedbackResponse, error) {
	path := fmt.Sprintf("/v1/apps/%s/betaFeedbackScreenshotSubmissions", appID)

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response FeedbackResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetCrashes retrieves TestFlight crash reports
func (c *Client) GetCrashes(ctx context.Context, appID string) (*CrashesResponse, error) {
	path := fmt.Sprintf("/v1/apps/%s/betaFeedbackCrashSubmissions", appID)

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response CrashesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetReviews retrieves App Store reviews
func (c *Client) GetReviews(ctx context.Context, appID string, opts ...ReviewOption) (*ReviewsResponse, error) {
	path := fmt.Sprintf("/v1/apps/%s/customerReviews", appID)
	query := buildReviewQuery(opts)
	if query != "" {
		path += "?" + query
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response ReviewsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// Links represents pagination links
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

// PrintJSON prints data as minified JSON (best for AI agents)
func PrintJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(data)
}

// PrintPrettyJSON prints data as indented JSON (best for debugging).
func PrintPrettyJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// PrintMarkdown prints data as Markdown table
func PrintMarkdown(data interface{}) error {
	switch v := data.(type) {
	case *FeedbackResponse:
		return printFeedbackMarkdown(v)
	case *CrashesResponse:
		return printCrashesMarkdown(v)
	case *ReviewsResponse:
		return printReviewsMarkdown(v)
	default:
		return PrintJSON(data)
	}
}

// PrintTable prints data as a formatted table
func PrintTable(data interface{}) error {
	switch v := data.(type) {
	case *FeedbackResponse:
		return printFeedbackTable(v)
	case *CrashesResponse:
		return printCrashesTable(v)
	case *ReviewsResponse:
		return printReviewsTable(v)
	default:
		return PrintJSON(data)
	}
}

// BuildRequestBody builds a JSON request body
func BuildRequestBody(data interface{}) (io.Reader, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}
	return &buf, nil
}

// ParseError parses an error response
func ParseError(body []byte) error {
	var errResp struct {
		Errors []struct {
			Code   string `json:"code"`
			Title  string `json:"title"`
			Detail string `json:"detail"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("unknown error: %s", string(body))
	}

	if len(errResp.Errors) > 0 {
		return fmt.Errorf("%s: %s", errResp.Errors[0].Title, errResp.Errors[0].Detail)
	}

	return fmt.Errorf("unknown error: %s", string(body))
}

// IsNotFound checks if the error is a "not found" error
func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "NOT_FOUND")
}

// IsUnauthorized checks if the error is an "unauthorized" error
func IsUnauthorized(err error) bool {
	return strings.Contains(err.Error(), "UNAUTHORIZED")
}

func compactWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func escapeMarkdown(input string) string {
	clean := compactWhitespace(input)
	return strings.ReplaceAll(clean, "|", "\\|")
}

func printFeedbackTable(resp *FeedbackResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Created\tEmail\tComment")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			item.Attributes.CreatedDate,
			item.Attributes.Email,
			compactWhitespace(item.Attributes.Comment),
		)
	}
	return w.Flush()
}

func printCrashesTable(resp *CrashesResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Created\tEmail\tDevice\tOS\tComment")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.Attributes.CreatedDate,
			item.Attributes.Email,
			item.Attributes.DeviceModel,
			item.Attributes.OSVersion,
			compactWhitespace(item.Attributes.Comment),
		)
	}
	return w.Flush()
}

func printReviewsTable(resp *ReviewsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Created\tRating\tTerritory\tTitle")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
			item.Attributes.CreatedDate,
			item.Attributes.Rating,
			item.Attributes.Territory,
			compactWhitespace(item.Attributes.Title),
		)
	}
	return w.Flush()
}

func printFeedbackMarkdown(resp *FeedbackResponse) error {
	fmt.Fprintln(os.Stdout, "| Created | Email | Comment |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
			escapeMarkdown(item.Attributes.CreatedDate),
			escapeMarkdown(item.Attributes.Email),
			escapeMarkdown(item.Attributes.Comment),
		)
	}
	return nil
}

func printCrashesMarkdown(resp *CrashesResponse) error {
	fmt.Fprintln(os.Stdout, "| Created | Email | Device | OS | Comment |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s |\n",
			escapeMarkdown(item.Attributes.CreatedDate),
			escapeMarkdown(item.Attributes.Email),
			escapeMarkdown(item.Attributes.DeviceModel),
			escapeMarkdown(item.Attributes.OSVersion),
			escapeMarkdown(item.Attributes.Comment),
		)
	}
	return nil
}

func printReviewsMarkdown(resp *ReviewsResponse) error {
	fmt.Fprintln(os.Stdout, "| Created | Rating | Territory | Title |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %d | %s | %s |\n",
			escapeMarkdown(item.Attributes.CreatedDate),
			item.Attributes.Rating,
			escapeMarkdown(item.Attributes.Territory),
			escapeMarkdown(item.Attributes.Title),
		)
	}
	return nil
}
