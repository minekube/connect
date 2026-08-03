package bedrockprincipal

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const MetadataPathV2 = "/.well-known/minekube-connect/bedrock-principal-v2.json"

type MetadataConfiguration struct {
	Origin      string
	Path        string
	Issuer      string
	TrustDomain string
	Audience    string
	Client      *http.Client
	Now         func() time.Time
}

type MetadataKeyProvider struct {
	configuration       MetadataConfiguration
	client              *http.Client
	now                 func() time.Time
	mu                  sync.Mutex
	keys                map[string]metadataKey
	freshUntil          time.Time
	nextRefresh         time.Time
	unknownRefreshAfter time.Time
	backoff             time.Duration
	etag                string
}

type metadataDocumentV2 struct {
	Issuer             string        `json:"issuer"`
	TrustDomain        string        `json:"trust_domain"`
	Audience           string        `json:"audience"`
	CacheMaxAgeSeconds int64         `json:"cache_max_age_seconds"`
	Keys               []metadataKey `json:"keys"`
}

type metadataKey struct {
	KID   string `json:"kid"`
	KTY   string `json:"kty"`
	CRV   string `json:"crv"`
	ALG   string `json:"alg"`
	Use   string `json:"use"`
	X     string `json:"x"`
	State string `json:"state"`
	key   ed25519.PublicKey
}

func NewMetadataKeyProvider(configuration MetadataConfiguration) (*MetadataKeyProvider, error) {
	origin, err := url.Parse(configuration.Origin)
	if err != nil || origin.Scheme != "https" || origin.Host == "" || origin.User != nil || origin.RawQuery != "" || origin.Fragment != "" || origin.Path != "" {
		return nil, fmt.Errorf("invalid metadata origin")
	}
	if configuration.Path != MetadataPathV2 || configuration.Issuer == "" || configuration.TrustDomain == "" || configuration.Audience == "" {
		return nil, fmt.Errorf("invalid metadata trust configuration")
	}
	if configuration.TrustDomain == "urn:minekube:connect:production" && configuration.Origin != "https://connect.minekube.com" {
		return nil, fmt.Errorf("production metadata origin is pinned")
	}
	baseClient := configuration.Client
	if baseClient == nil {
		transport := http.DefaultTransport
		if defaultTransport, ok := transport.(*http.Transport); ok {
			clonedTransport := defaultTransport.Clone()
			clonedTransport.DialContext = (&net.Dialer{Timeout: 2 * time.Second, KeepAlive: 30 * time.Second}).DialContext
			transport = clonedTransport
		}
		baseClient = &http.Client{Transport: transport}
	}
	clientCopy := *baseClient
	clientCopy.Timeout = 5 * time.Second
	clientCopy.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
	now := configuration.Now
	if now == nil {
		now = time.Now
	}
	return &MetadataKeyProvider{configuration: configuration, client: &clientCopy, now: now, keys: map[string]metadataKey{}}, nil
}

func (p *MetadataKeyProvider) Eligible(ctx context.Context, trustDomain, kid string) (ed25519.PublicKey, error) {
	if trustDomain != p.configuration.TrustDomain || kid == "" {
		return nil, Trust
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	if !now.Before(p.freshUntil) {
		if now.Before(p.nextRefresh) {
			return nil, MetadataUnavailable
		}
		if err := p.refresh(ctx, now); err != nil {
			if p.backoff == 0 {
				p.backoff = time.Second
			} else {
				p.backoff *= 2
			}
			if p.backoff > 30*time.Second {
				p.backoff = 30 * time.Second
			}
			p.nextRefresh = now.Add(p.backoff)
			return nil, MetadataUnavailable
		}
		p.backoff = 0
		p.nextRefresh = time.Time{}
	}
	key, ok := p.keys[kid]
	if !ok {
		if now.Before(p.unknownRefreshAfter) {
			return nil, Trust
		}
		p.unknownRefreshAfter = now.Add(5 * time.Second)
		if err := p.refresh(ctx, now); err != nil {
			return nil, MetadataUnavailable
		}
		key, ok = p.keys[kid]
		if !ok {
			return nil, Trust
		}
	}
	switch key.State {
	case "current", "next", "previous":
		return append(ed25519.PublicKey(nil), key.key...), nil
	case "revoked":
		return nil, KeyRevoked
	default:
		return nil, Trust
	}
}

func (p *MetadataKeyProvider) refresh(ctx context.Context, now time.Time) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, p.configuration.Origin+p.configuration.Path, nil)
	if err != nil {
		return err
	}
	if p.etag != "" {
		request.Header.Set("If-None-Match", p.etag)
	}
	response, err := p.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotModified {
		return fmt.Errorf("metadata 304 cannot extend expired cache")
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("metadata status %d", response.StatusCode)
	}
	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
		return fmt.Errorf("metadata content type")
	}
	raw, err := io.ReadAll(io.LimitReader(response.Body, 65537))
	if err != nil || len(raw) > 65536 {
		return fmt.Errorf("metadata size")
	}
	if err := rejectDuplicateMembers(raw, 4); err != nil {
		return err
	}
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	var document metadataDocumentV2
	if err := decoder.Decode(&document); err != nil {
		return err
	}
	if err := ensureJSONEOF(decoder); err != nil {
		return err
	}
	if document.Issuer != p.configuration.Issuer || document.TrustDomain != p.configuration.TrustDomain || document.Audience != p.configuration.Audience || document.CacheMaxAgeSeconds < 1 || document.CacheMaxAgeSeconds > 300 || len(document.Keys) < 1 || len(document.Keys) > 16 {
		return fmt.Errorf("metadata trust or bounds")
	}
	keys := make(map[string]metadataKey, len(document.Keys))
	for _, value := range document.Keys {
		if _, duplicate := keys[value.KID]; duplicate {
			return fmt.Errorf("duplicate kid")
		}
		if value.KID == "" || len(value.KID) > 128 || value.KTY != "OKP" || value.CRV != "Ed25519" || value.ALG != "EdDSA" || value.Use != "sig" {
			return fmt.Errorf("invalid key metadata")
		}
		if value.State != "current" && value.State != "next" && value.State != "previous" && value.State != "revoked" && value.State != "disabled" {
			return fmt.Errorf("invalid key state")
		}
		decoded, err := decodeCanonical(value.X, ed25519.PublicKeySize)
		if err != nil || len(decoded) != ed25519.PublicKeySize {
			return fmt.Errorf("invalid public key")
		}
		value.key = append(ed25519.PublicKey(nil), decoded...)
		keys[value.KID] = value
	}
	p.keys = keys
	p.freshUntil = now.Add(time.Duration(document.CacheMaxAgeSeconds) * time.Second)
	p.etag = response.Header.Get("ETag")
	return nil
}

func ensureJSONEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("trailing json")
	}
	return nil
}
