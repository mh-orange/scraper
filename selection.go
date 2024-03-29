// Copyright 2019 Andrew Bates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scraper

import (
	"strings"

	"golang.org/x/net/html"
)

type selection struct {
	*html.Node
}

func (n *selection) value(t *tag) (value string) {
	switch t.typ {
	case text:
		value = n.text()
	case attr:
		value = n.attr(t.detail)
	}
	return value
}

func (n *selection) text() (val string) {
	var buf strings.Builder

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n == nil {
			return
		}

		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n.Node)

	return buf.String()
}

func (n *selection) attr(name string) (val string) {
	for i, a := range n.Attr {
		if a.Key == name {
			val = n.Attr[i].Val
		}
	}
	return
}
