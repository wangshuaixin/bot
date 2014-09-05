package url

import (
	"github.com/fabioxgn/go-bot"
	"github.com/fabioxgn/go-bot/web"
	"html"
	"net/url"
	"regexp"
	"strings"
)

const (
	minDomainLength = 3
)

var (
	re = regexp.MustCompile("<title>\\n*?(.*?)\\n*?<\\/title>")
)

func canBeURLWithoutProtocol(text string) bool {
	return len(text) > minDomainLength &&
		!strings.HasPrefix(text, "http") &&
		strings.Contains(text, ".")
}

func extractURL(text string) string {
	extractedURL := ""
	for _, value := range strings.Split(text, " ") {
		if canBeURLWithoutProtocol(value) {
			value = "http://" + value
		}

		parsedURL, err := url.Parse(value)
		if err != nil {
			continue
		}
		if strings.HasPrefix(parsedURL.Scheme, "http") {
			extractedURL = parsedURL.String()
			break
		}
	}
	return extractedURL
}

func getTitle(text string, get web.GetBodyFunc) (string, error) {
	URL := extractURL(text)

	if URL == "" {
		return "", nil
	}

	body, err := get(URL)
	if err != nil {
		return "", err
	}

	title := re.FindString(string(body))
	if title == "" {
		return "", nil
	}

	title = strings.Replace(title, "\n", "", -1)
	title = title[strings.Index(title, ">")+1 : strings.LastIndex(title, "<")]

	return html.UnescapeString(title), nil
}

func urlTitle(command *bot.PassiveCmd) (string, error) {
	return getTitle(command.Raw, web.GetBody)
}

func init() {
	bot.RegisterPassiveCommand(
		"url",
		urlTitle)
}
