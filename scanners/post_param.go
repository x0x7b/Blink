package scanners

import (
	"Blink/types"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func TestForms(baseline types.BlinkResponse, fc types.FlagCondition) (types.BlinkResponse, []types.BlinkResponse, types.BlinkError) {
	var forms []types.Form
	var response types.BlinkResponse
	var results []types.BlinkResponse
	parsed := strings.NewReader(string(baseline.Body))
	doc, err := html.Parse(parsed)
	if err != nil {
		return response, results, types.BlinkError{Stage: "No"}
	}
	parseForms(doc, &forms)
	for _, form := range forms {
		fmt.Printf("%v\n", form)
	}
	return response, results, types.BlinkError{Stage: "No"}
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
			parsedForm.Method = a.Val
		}
		if a.Key == "action" {
			parsedForm.Action = a.Val
		}
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
			inputs = append(inputs, input)
		}
	}
	parsedForm.Inputs = inputs
	return parsedForm
}
