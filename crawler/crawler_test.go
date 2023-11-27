package crawler

import "testing"

var cases = []struct {
	name     string
	title    string
	teams    []string
	terms    []string
	excludes []string
	expected bool
}{
	{
		"PassingCase",
		"AC Milan vs Inter highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"highlights"},
		[]string{"reports"},
		true,
	},
	{
		"PassingCaseCaseInsensitive",
		"AC milan vs inter highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"HIGHLIGHTS"},
		[]string{"reports"},
		true,
	},
	{
		"KeyWordsMatchesTeamDoesNot",
		"AC Milan vs Inter highlights",
		[]string{"Napoli", "Juventus"},
		[]string{"highlights"},
		[]string{"reports"},
		false,
	},
	{
		"KeyWordsDontMatchTeamDoes",
		"AC Milan vs Inter highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"Incorrect"},
		[]string{"reports"},
		false,
	},
	{
		"OnlyKeywords is there",
		"Lyon highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"highlights"},
		[]string{"reports"},
		false,
	},
	{
		"Everything is ok, but it contains a keyword that should be excluded",
		"Reports AC Milan vs Inter highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"highlights"},
		[]string{"reports"},
		false,
	},
}

func TestMatch(t *testing.T) {
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			configs := &Configs{
				[]Channel{
					{ID: "BLAH"},
				},
				c.teams,
				c.terms,
				c.excludes,
				"FAKE",
			}
			out := matches(c.title, configs)
			if out != c.expected {
				t.Errorf("=> %v, want %v", out, c.expected)
			}
		})
	}
}
