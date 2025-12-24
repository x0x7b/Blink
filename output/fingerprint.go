package output

import (
	"Blink/types"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func serverFingerprint(bl types.BlinkResponse, fc types.FlagCondition) {
	var out strings.Builder
	if !fc.ShowFp {
		return
	}
	out.WriteString(types.Cyan + "\nServer Fingerprint:\n" + types.Reset)
	out.WriteString("   Server: " + bl.Headers.Get("Server") + "\n")
	if bl.Headers.Get("X-Powered-By") != "" {
		out.WriteString("   Tech: " + bl.Headers.Get("X-Powered-By") + "\n")

	}
	if bl.Headers.Get("X-Powered-CMS") != "" {
		out.WriteString(types.Cyan + "   CMS: " + types.Reset + bl.Headers.Get("X-Powered-CMS") + "\n")

	}
	if bl.Headers.Get("X-Frame-Options") == "" {
		out.WriteString(types.Cyan + "   [INFO] " + types.Reset + "Missing X-Frame-Options header. (Clickjacking-related behavior)\n")
	}
	out.WriteString(types.Cyan + "Defined endpoints: \n" + types.Reset)
	out.WriteString(getLinks(bl))

	fmt.Println(out.String())

}

func getLinks(bl types.BlinkResponse) string {
	var out strings.Builder
	parsed := strings.NewReader(string(bl.Body))
	doc, err := html.Parse(parsed)
	if err != nil {
		ErrorOutput(types.BlinkError{Message: err.Error()})
		return out.String()
	}
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					out.WriteString(a.Val + "\n")
				}
			}
		}
	}
	return out.String()
}
