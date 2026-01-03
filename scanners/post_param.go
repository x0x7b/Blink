package scanners

import (
	"Blink/core"
	"Blink/output"
	"Blink/types"
	"bufio"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func TestForms(baseline types.BlinkResponse, fc types.FlagCondition, report func(types.Progress)) (types.BlinkResponse, [][]types.BlinkResponse, types.BlinkError) {
	var forms []types.Form
	response := types.BlinkResponse{}
	var baselineSumb types.BlinkResponse
	var results [][]types.BlinkResponse
	parsed := strings.NewReader(string(baseline.Body))
	doc, err := html.Parse(parsed)
	if err != nil {
		return response, results, types.BlinkError{Message: err.Error()}
	}
	file, err := os.Open(fc.Wordlist)
	if err != nil {
		return response, results, types.BlinkError{Message: err.Error()}
	}
	defer file.Close()
	var payloads []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		payloads = append(payloads, scanner.Text())
	}
	parseForms(doc, &forms)
	for _, form := range forms {
		var formResult []types.BlinkResponse
		baseURL, err := url.Parse(baseline.URL)
		if err != nil {
			return response, results, types.BlinkError{Message: err.Error()}
		}

		newURLurl, err := baseURL.Parse(form.Action)
		if err != nil {
			return response, results, types.BlinkError{Message: err.Error()}
		}
		newURL := newURLurl.String()
		values := make(url.Values)
		for _, input := range form.Inputs {
			values.Add(input.Name, "test")
		}
		body := values.Encode()

		fc2 := fc
		fc2.FollowRedirects = false
		fc2.Data = body
		baselineSumb, _, errb := core.HttpRequest(form.Method, newURL, fc2)
		if errb.Stage != "OK" && errb.Stage != "INFO" {
			continue
		}
		output.ErrorOutput(errb)
		formResult = append(formResult, baselineSumb)
		for key := range values {
			for i, payload := range payloads {
				mut := make(url.Values)
				for k, v := range values {
					cp := make([]string, len(v))
					copy(cp, v)
					mut[k] = cp
				}
				mut.Set(key, payload)
				fc2 := fc
				fc2.FollowRedirects = false
				fc2.Data = mut.Encode()
				if report != nil {
					report(types.Progress{
						Stage:   "URL_PARAMS",
						Current: i + 1,
						Target:  form.Name,
						Total:   len(payloads),
					})
				}
				test, _, errb := core.HttpRequest(form.Method, newURL, fc2)
				if errb.Stage != "OK" && errb.Stage != "INFO" {
					continue
				}
				formResult = append(formResult, test)

			}

		}
		results = append(results, formResult)

	}
	return baselineSumb, results, types.BlinkError{Stage: "OK"}
}

func parseForms(f *html.Node, forms *[]types.Form) {
	for n := range f.Descendants() {
		if n.Type == html.ElementNode && n.Data == "form" {
			*forms = append(*forms, parseForm(n))
		}
	}
}

func parseForm(n *html.Node) types.Form {
	var parsedForm types.Form
	var inputs []types.Input
	for _, a := range n.Attr {
		if a.Key == "name" {
			parsedForm.Name = a.Val
		}
		if a.Key == "method" {
			parsedForm.Method = strings.ToUpper(a.Val)
		}
		if a.Key == "action" {
			parsedForm.Action = a.Val
		}

	}
	if parsedForm.Method == "" {
		parsedForm.Method = "GET"
	}
	for c := range n.Descendants() {
		if c.Type == html.ElementNode && c.Data == "input" {
			var input types.Input
			for _, b := range c.Attr {
				if b.Key == "name" {
					input.Name = b.Val
				}
				if b.Key == "type" {
					input.Type = b.Val
				}
			}
			if input.Name != "" && input.Type != "submit" && input.Type != "button" && input.Type != "image" {
				inputs = append(inputs, input)
			}

		}
	}
	parsedForm.Inputs = inputs
	return parsedForm
}
