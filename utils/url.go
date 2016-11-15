// Package url created to provide Marshaling and Unmarshaling for url.URL
package utils

import (
	"bytes"
	"fmt"
	nativeurl "net/url"
	"strings"
)

// URL is native url.URL struct
type URL nativeurl.URL

// Parse parses rawurl into a URL structure.
// The rawurl may be relative or absolute.
func Parse(rawurl string) (*URL, error) {
	un, err := nativeurl.Parse(rawurl)
	u := URL(*un)
	return &u, err
}

// ParseRequestURI parses rawurl into a URL structure. It assumes that rawurl
// was received in an HTTP request, so the rawurl is interpreted only as an
// absolute URI or an absolute path. The string rawurl is assumed not to have a
// #fragment suffix. (Web browsers strip #fragment before sending the URL to a
// web server.)
func ParseRequestURI(rawurl string) (*URL, error) {
	un, err := nativeurl.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}
	u := URL(*un)
	return &u, nil
}

func NewUrl(rawurl string) *URL {
	un, _ := nativeurl.Parse(rawurl)
	u := URL(*un)
	return &u
}

// UnmarshalText calls url.Parse
func (u *URL) UnmarshalText(p []byte) error {
	nu, err := nativeurl.Parse(string(p))
	if err != nil {
		return err
	}
	(*u) = URL(*nu)
	return nil
}

// UnmarshalBinary calls url.Parse
func (u *URL) UnmarshalBinary(p []byte) error {
	nu, err := nativeurl.Parse(string(p))
	if err != nil {
		return err
	}
	(*u) = URL(*nu)
	return nil
}

// MarshalText just calls String()
func (u *URL) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// MarshalBinary just calls String()
func (u *URL) MarshalBinary() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalJSON parses JSON string into url.URL
func (u *URL) UnmarshalJSON(p []byte) error {
	nu, err := nativeurl.Parse(string(bytes.Trim(p, `"`)))
	if err != nil {
		return err
	}
	(*u) = URL(*nu)
	return nil
}

// MarshalJSON turns url into a JSON string
func (u *URL) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"%s"`, u.String())
	return []byte(s), nil
}

// GetPath returns url.Path with leading '/' removed
func (u *URL) GetPath() string {
	if u == nil {
		return ""
	}
	return strings.TrimLeft(u.Path, "/")
}

// GetHost url.Host without the port suffix
func (u *URL) GetHost() string {
	if u == nil {
		return ""
	}
	i := strings.Index(u.Host, ":")
	if i == -1 {
		return u.Host
	}
	return u.Host[0:i]
}

// String returns the string representation
func (u *URL) String() string {
	return (*nativeurl.URL)(u).String()
}

func (u *URL) Query() nativeurl.Values {
	nu := nativeurl.URL(*u)
	return nu.Query()
}
