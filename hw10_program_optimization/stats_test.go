//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestCountDomains_Minimal(t *testing.T) {
	t.Run("counts matching domains", func(t *testing.T) {
		users := users{
			{Email: "user1@sub.example.com"},
			{Email: "user2@test.example.com"},
			{Email: "user3@example.org"},
		}

		result, _ := countDomains(users, "com")

		expected := DomainStat{
			"sub.example.com":  1,
			"test.example.com": 1,
		}

		for domain, count := range expected {
			if result[domain] != count {
				t.Errorf("expected %s: %d, got %d", domain, count, result[domain])
			}
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		users := users{
			{Email: "USER@EXAMPLE.COM"},
		}

		result, _ := countDomains(users, "com")

		if result["example.com"] != 1 {
			t.Error("should match case insensitive")
		}
	})

	t.Run("skips invalid emails", func(t *testing.T) {
		users := users{
			{Email: "invalid-email"},
			{Email: "valid@example.com"},
		}

		result, _ := countDomains(users, "com")

		if result["example.com"] != 1 {
			t.Error("should count valid emails only")
		}
	})
}
