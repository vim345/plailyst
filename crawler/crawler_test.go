package crawler

import "testing"

var cases = []struct {
	name     string
	title    string
	teams    []string
	keywords []string
	expected bool
}{
	{
		"PassingCase",
		"AC Milan vs Inter highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"highlights"},
		true,
	},
	{
		"PassingCaseCaseInsensitive",
		"AC milan vs inter highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"HIGHLIGHTS"},
		true,
	},
	{
		"KeyWordsMatchesTeamDoesNot",
		"AC Milan vs Inter highlights",
		[]string{"Napoli", "Juventus"},
		[]string{"highlights"},
		false,
	},
	{
		"KeyWordsDontMatchTeamDoes",
		"AC Milan vs Inter highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"Incorrect"},
		false,
	},
	{
		"OnlyKeywords is there",
		"Lyon highlights",
		[]string{"AC Milan", "Juventus"},
		[]string{"highlights"},
		false,
	},
}

func TestMatch(t *testing.T) {
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			out := matches(c.title, c.teams, c.keywords)
			if out != c.expected {
				t.Errorf("=> %v, want %v", out, c.expected)
			}
		})
	}
}
