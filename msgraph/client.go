package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
	"github.com/askasoft/pango/ret"
	"github.com/askasoft/pango/str"
)

type Credential struct {
	Time        time.Time `json:"-"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	AccessToken string    `json:"access_token"`
}

func (c *Credential) IsValid() bool {
	return c != nil && c.Time.Add(time.Duration(c.ExpiresIn-60)*time.Second).After(time.Now())
}

type GraphClient struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	Scope        string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	MaxRetries  int
	RetryAfter  time.Duration
	ShouldRetry func(error) bool // default retry on not canceled error or (status = 429 || (status >= 500 && status <= 599))

	credential Credential
}

func (gc *GraphClient) Endpoint(format string, a ...any) string {
	return "https://graph.microsoft.com/v1.0" + fmt.Sprintf(format, a...)
}

func (gc *GraphClient) shouldRetry(err error) bool {
	sr := gc.ShouldRetry
	if sr == nil {
		sr = shouldRetry
	}
	return sr(err)
}

func (gc *GraphClient) call(req *http.Request) (res *http.Response, err error) {
	client := &http.Client{
		Transport: gc.Transport,
		Timeout:   gc.Timeout,
	}

	res, err = httplog.TraceClientDo(gc.Logger, client, req)
	if err != nil {
		if gc.shouldRetry(err) {
			err = ret.NewRetryError(err, gc.RetryAfter)
		}
		return res, err
	}

	return res, nil
}

func (gc *GraphClient) RetryForError(ctx context.Context, api func() error) (err error) {
	return ret.RetryForError(ctx, api, gc.MaxRetries, gc.Logger)
}

func (gc *GraphClient) authenticate(ctx context.Context, req *http.Request) error {
	if !gc.credential.IsValid() {
		if err := gc.DoAuth(ctx); err != nil {
			return err
		}
	}
	req.Header.Set("Authorization", gc.credential.TokenType+" "+gc.credential.AccessToken)
	return nil
}

func (gc *GraphClient) DoAuth(ctx context.Context) error {
	return gc.RetryForError(ctx, func() error {
		return gc.doAuth(ctx)
	})
}

func (gc *GraphClient) doAuth(ctx context.Context) error {
	vals := url.Values{}
	vals.Add("client_id", gc.ClientID)
	vals.Add("client_secret", gc.ClientSecret)
	vals.Add("grant_type", "client_credentials")
	vals.Add("scope", str.IfEmpty(gc.Scope, "https://graph.microsoft.com/.default"))

	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", gc.TenantID)
	//TODO:
	// val := urlx.EncodeValues(
	// 	"client_id", gc.ClientID,
	// 	"client_secret", gc.ClientSecret,
	// 	"grant_type", "client_credentials",
	// 	"scope", str.IfEmpty(gc.Scope, "https://graph.microsoft.com/.default"),
	// )

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(vals.Encode()))
	if err != nil {
		return err
	}

	res, err := gc.call(req)
	if err != nil {
		return err
	}
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK {
		if err = decoder.Decode(&gc.credential); err != nil {
			return err
		}
		gc.credential.Time = time.Now()
		return nil
	}

	ae := newAuthError(res)
	_ = decoder.Decode(ae)

	if gc.shouldRetry(ae) {
		ae.RetryAfter = gc.RetryAfter
	}
	return ae
}

func (gc *GraphClient) authAndCall(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := gc.authenticate(ctx, req); err != nil {
		return nil, err
	}
	return gc.call(req)
}

func (gc *GraphClient) doCall(ctx context.Context, req *http.Request, result any) error {
	res, err := gc.authAndCall(ctx, req)
	if err != nil {
		return err
	}

	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK {
		if result != nil {
			return decoder.Decode(result)
		}
		return nil
	}

	re := newResultError(res)
	_ = decoder.Decode(re)

	if gc.shouldRetry(re) {
		re.RetryAfter = gc.RetryAfter
	}
	return re
}

func (gc *GraphClient) DoGet(ctx context.Context, url string, result any) error {
	return gc.RetryForError(ctx, func() error {
		return gc.doGet(ctx, url, result)
	})
}

func (gc *GraphClient) doGet(ctx context.Context, url string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return gc.doCall(ctx, req, result)
}

func (gc *GraphClient) DoPost(ctx context.Context, url string, body io.Reader, result any) error {
	return gc.RetryForError(ctx, func() error {
		return gc.doPost(ctx, url, body, result)
	})
}

func (gc *GraphClient) doPost(ctx context.Context, url string, body io.Reader, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}

	return gc.doCall(ctx, req, result)
}

func (gc *GraphClient) DoCopyFile(ctx context.Context, url string, w io.Writer) error {
	return gc.RetryForError(ctx, func() error {
		return gc.doCopyFile(ctx, url, w)
	})
}

func (gc *GraphClient) doCopyFile(ctx context.Context, url string, w io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := gc.authAndCall(ctx, req)
	if err != nil {
		return err
	}

	return copyResponse(res, w)
}

func (gc *GraphClient) DoReadFile(ctx context.Context, url string) (buf []byte, err error) {
	err = gc.RetryForError(ctx, func() error {
		buf, err = gc.doReadFile(ctx, url)
		return err
	})
	return
}

func (gc *GraphClient) doReadFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := gc.authAndCall(ctx, req)
	if err != nil {
		return nil, err
	}

	return readResponse(res)
}

func (gc *GraphClient) DoSaveFile(ctx context.Context, url string, path string) error {
	return gc.RetryForError(ctx, func() error {
		return gc.doSaveFile(ctx, url, path)
	})
}

func (gc *GraphClient) doSaveFile(ctx context.Context, url string, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := gc.authAndCall(ctx, req)
	if err != nil {
		return err
	}

	return saveResponse(res, path)
}

func DoGets[T any](ctx context.Context, gc *GraphClient, url string) ([]T, string, error) {
	type R struct {
		NextLink string `json:"@odata.nextLink"`
		Values   []T    `json:"value"`
	}

	r := &R{}
	err := gc.DoGet(ctx, url, r)
	return r.Values, r.NextLink, err
}

func DoList[T any](ctx context.Context, gc *GraphClient, url string) (vs []T, err error) {
	type R struct {
		NextLink string `json:"@odata.nextLink"`
		Values   []T    `json:"value"`
	}

	for url != "" {
		r := &R{}
		if err = gc.DoGet(ctx, url, r); err != nil {
			return
		}

		vs = append(vs, r.Values...)
		url = r.NextLink
	}
	return
}

func DoIter[T any](ctx context.Context, gc *GraphClient, url string, itf func(T) error) error {
	type R struct {
		NextLink string `json:"@odata.nextLink"`
		Values   []T    `json:"value"`
	}

	for url != "" {
		r := &R{}
		if err := gc.DoGet(ctx, url, r); err != nil {
			return err
		}

		for _, v := range r.Values {
			if err := itf(v); err != nil {
				return err
			}
		}

		url = r.NextLink
	}

	return nil
}

func optionsQuery(options ...string) string {
	// TODO: return urlx.EncodeQuery(options...)
	if len(options) == 0 {
		return ""
	}
	return "?" + encodeOptions(options...)
}

func encodeOptions(options ...string) string {
	z := len(options)
	if z == 0 {
		return ""
	}

	vs := url.Values{}
	for i := 0; i < z; i += 2 {
		k := options[i]
		v := ""
		if i+1 < z {
			v = options[i+1]
		}
		vs.Add(k, v)
	}
	return vs.Encode()
}

func copyResponse(res *http.Response, w io.Writer) error {
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return newResultError(res)
	}

	_, err := io.Copy(w, res.Body)
	return err
}

func readResponse(res *http.Response) ([]byte, error) {
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, newResultError(res)
	}

	return io.ReadAll(res.Body)
}

func saveResponse(res *http.Response, path string) error {
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return newResultError(res)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0770); err != nil {
		return err
	}

	return fsu.WriteReader(path, res.Body, 0660)
}

func toString(o any) string {
	return jsonx.Prettify(o)
}
