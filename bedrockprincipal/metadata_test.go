package bedrockprincipal

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMetadataKeyLifecycle(t *testing.T) {
	states := []struct {
		state string
		want  PrincipalError
	}{{"current", ""}, {"next", ""}, {"previous", ""}, {"revoked", KeyRevoked}, {"disabled", Trust}}
	for _, tc := range states {
		t.Run(tc.state, func(t *testing.T) {
			server := metadataServer(t, metadataDocument(tc.state), "application/json")
			provider := newTestMetadataProvider(t, server)
			key, err := provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "connect-v2-test")
			if tc.want == "" {
				require.NoError(t, err)
				require.Equal(t, testPublicKey, key)
			} else {
				require.ErrorIs(t, err, tc.want)
			}
		})
	}
}

func TestMetadataRejectsWrongContentTypeDuplicateKIDAndOversize(t *testing.T) {
	duplicate := strings.Replace(metadataDocument("current"), `]}`, `,{"kid":"connect-v2-test","kty":"OKP","crv":"Ed25519","alg":"EdDSA","use":"sig","x":"`+base64.RawURLEncoding.EncodeToString(testPublicKey)+`","state":"next"}]}`, 1)
	duplicateMember := strings.Replace(metadataDocument("current"), `{"issuer":"minekube-connect-test",`, `{"issuer":"minekube-connect-test","issuer":"minekube-connect-test",`, 1)
	cases := []struct{ name, body, contentType string }{
		{"wrong-content-type", metadataDocument("current"), "text/plain"},
		{"duplicate-kid", duplicate, "application/json"},
		{"duplicate-member", duplicateMember, "application/json"},
		{"oversize", strings.Repeat(" ", 65537), "application/json"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := metadataServer(t, tc.body, tc.contentType)
			_, err := newTestMetadataProvider(t, server).Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "connect-v2-test")
			require.ErrorIs(t, err, MetadataUnavailable)
		})
	}
}

func TestMetadataUsesFreshCacheButNotExpiredCacheDuringOutage(t *testing.T) {
	var requests atomic.Int32
	failing := atomic.Bool{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		if failing.Load() {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, metadataDocument("current"))
	}))
	t.Cleanup(server.Close)
	now := time.Unix(testNowUnix, 0)
	provider, err := NewMetadataKeyProvider(MetadataConfiguration{
		Origin: server.URL, Path: "/.well-known/minekube-connect/bedrock-principal-v2.json",
		Issuer: "minekube-connect-test", TrustDomain: "urn:minekube:connect:test:corpus-v2",
		Audience: "urn:minekube:connect:test:bedrock-principal:v2", Client: server.Client(), Now: func() time.Time { return now },
	})
	require.NoError(t, err)
	_, err = provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "connect-v2-test")
	require.NoError(t, err)
	failing.Store(true)
	_, err = provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "connect-v2-test")
	require.NoError(t, err)
	require.Equal(t, int32(1), requests.Load())
	now = now.Add(2 * time.Second)
	_, err = provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "connect-v2-test")
	require.ErrorIs(t, err, MetadataUnavailable)
}

func TestMetadataUnknownKIDTriggersOneCoalescedRefreshWindow(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		request := requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		kid := "old"
		if request > 1 {
			kid = "new"
		}
		fmt.Fprint(w, strings.Replace(metadataDocument("current"), "connect-v2-test", kid, 1))
	}))
	t.Cleanup(server.Close)
	provider := newTestMetadataProvider(t, server)
	_, err := provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "old")
	require.NoError(t, err)
	_, err = provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "new")
	require.NoError(t, err)
	require.Equal(t, int32(2), requests.Load())
	_, err = provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "missing")
	require.ErrorIs(t, err, Trust)
	require.Equal(t, int32(2), requests.Load())
}

func TestMetadata304CannotExtendPastOriginalCacheCeiling(t *testing.T) {
	var sawConditional atomic.Bool
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Header.Get("If-None-Match") == `"fixture"` {
			sawConditional.Store(true)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ETag", `"fixture"`)
		fmt.Fprint(w, metadataDocument("current"))
	}))
	t.Cleanup(server.Close)
	now := time.Unix(testNowUnix, 0)
	provider, err := NewMetadataKeyProvider(MetadataConfiguration{Origin: server.URL, Path: MetadataPathV2, Issuer: "minekube-connect-test", TrustDomain: "urn:minekube:connect:test:corpus-v2", Audience: "urn:minekube:connect:test:bedrock-principal:v2", Client: server.Client(), Now: func() time.Time { return now }})
	require.NoError(t, err)
	_, err = provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "connect-v2-test")
	require.NoError(t, err)
	now = now.Add(2 * time.Second)
	_, err = provider.Eligible(context.Background(), "urn:minekube:connect:test:corpus-v2", "connect-v2-test")
	require.ErrorIs(t, err, MetadataUnavailable)
	require.True(t, sawConditional.Load())
}

func newTestMetadataProvider(t *testing.T, server *httptest.Server) *MetadataKeyProvider {
	t.Helper()
	provider, err := NewMetadataKeyProvider(MetadataConfiguration{
		Origin: server.URL, Path: "/.well-known/minekube-connect/bedrock-principal-v2.json",
		Issuer: "minekube-connect-test", TrustDomain: "urn:minekube:connect:test:corpus-v2",
		Audience: "urn:minekube:connect:test:bedrock-principal:v2", Client: server.Client(),
		Now: func() time.Time { return time.Unix(testNowUnix, 0) },
	})
	require.NoError(t, err)
	return provider
}

func metadataServer(t *testing.T, body, contentType string) *httptest.Server {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", contentType)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)
	return server
}

func metadataDocument(state string) string {
	return fmt.Sprintf(`{"issuer":"minekube-connect-test","trust_domain":"urn:minekube:connect:test:corpus-v2","audience":"urn:minekube:connect:test:bedrock-principal:v2","cache_max_age_seconds":1,"keys":[{"kid":"connect-v2-test","kty":"OKP","crv":"Ed25519","alg":"EdDSA","use":"sig","x":"%s","state":"%s"}]}`, base64.RawURLEncoding.EncodeToString(testPublicKey), state)
}
