// Package logger provides a logger which will write log entries to Plex Media Server
package logger

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
)

// Log levels known by Plex
const (
	PlexLogError = iota
	PlexLogWarning
	PlexLogInfo
	PlexLogDebug
	PlexLogVerbose
)

// NewPlexLogger returns a PlexLogger instance that has URL preset
//
// URL should be the base url for Plex Media Server, `/log` path will be added
func NewPlexLogger(name, token, plexurl string, opts Options) (logr.Logger, error) {
	u, err := url.Parse(plexurl)
	if err != nil {
		return logr.Logger{}, fmt.Errorf("unable to parse url: %v", err)
	}

	u.Path = strings.TrimSuffix(u.Path, "/") + "/log"

	fopts := funcr.Options{
		Verbosity: opts.Verbosity,
	}

	l := &PlexLogSink{
		Formatter: funcr.NewFormatter(fopts),

		plexURL:   u,
		plexToken: token,
	}
	return logr.New(l.WithName(name)), nil
}

type Options struct {
	// Verbosity tells the logger which V logs to write. Higher values enable more logs.
	// Default is 0
	Verbosity int
}

// PlexLogSink is a single instance of Plex which is used for logging
type PlexLogSink struct {
	funcr.Formatter

	plexURL   *url.URL // Plex url, includes plex source. (http://127.0.0.1:32400/?source=Transcoder)
	plexToken string   // Plex token for authentication
}

// Init is not implemented for PlexLogSink
func (PlexLogSink) Init(_ logr.RuntimeInfo) {}

// Info level logs are written directly to Plex
func (l PlexLogSink) Info(level int, msg string, kvs ...interface{}) {
	prefix, m := l.FormatInfo(level, msg, kvs)
	l.send(prefix, PlexLogInfo, m)
}

// Error logs will have the error message passed as error key
func (l PlexLogSink) Error(err error, msg string, kvs ...interface{}) {
	prefix, m := l.FormatError(err, msg, kvs)
	l.send(prefix, PlexLogError, m)
}

// WithName adds an element to the logger name
func (l PlexLogSink) WithName(name string) logr.LogSink {
	l.Formatter.AddName(name)
	return &l
}

// WithValues adds key value pairs to the logger
func (l PlexLogSink) WithValues(kvs ...interface{}) logr.LogSink {
	l.Formatter.AddValues(kvs)
	return &l
}

// send message to PMS. Wrap all key value pairs to a text string since Plex has
// no concept of metadata other than message level.
//
// The request includes Plex token if it's available through the environment
func (l PlexLogSink) send(prefix string, level int, msg string) {
	u := l.getURL()
	q := u.Query()
	q.Set("level", fmt.Sprintf("%d", level))
	q.Set("message", msg)
	q.Set("source", prefix)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		// We have an error, but no place to report it. Bail out!
		return
	}

	plexToken := l.getPlexToken()
	if plexToken != "" {
		req.Header.Add("X-Plex-Token", plexToken)
	}
	req.Header.Add("User-Agent", "PlexLogger")

	// Ignore results
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("ERROR! %v", err)
	}
}

// getURL returns either the set URL or the default URL if unset
func (l PlexLogSink) getURL() url.URL {
	if l.plexURL == nil {
		u, _ := url.Parse("http://127.0.0.1:32400/log")
		return *u
	}
	return *l.plexURL
}

// getPlexToken returns the plex token from struct if it exists or falls back to
// X_PLEX_TOKEN environment variable
func (l PlexLogSink) getPlexToken() string {
	if l.plexToken != "" {
		return l.plexToken
	}
	return os.Getenv("X_PLEX_TOKEN")
}
