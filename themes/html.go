package themes

// Functions directly related to HTML, CSS, Amber or GCSS

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	// DefaultStyleFilename is the default stylesheet filename (GCSS)
	DefaultStyleFilename = "style.gcss"
)

var (
	// MetaKeywords contains a selection of allowed keywords for the HTML meta tag
	MetaKeywords = []string{"application-name", "author", "description", "generator", "keywords", "robots", "language", "googlebot", "Slurp", "bingbot", "geo.position", "geo.placename", "geo.region", "ICBM", "viewport"}
)

// MessagePage is an easy way to output a HTML page only given a title, the body
// (will be placed between the <body></body> tags) and the name of one of the
// built-in themes.
// Deprecated
func MessagePage(title, body, theme string) string {
	return fmt.Sprintf("<!doctype html><html><head><title>%s</title>%s<head><body><h1>%s</h1>%s</body></html>", title, StyleHead(theme), title, body)
}

// StyleHead returns contents that goes in "<head>", as bytes.
// This is either CSS wrapped in a "<style>" tag, or "<link>" tags to CSS and JS.
func StyleHead(theme string) []byte {
	var buf bytes.Buffer
	if theme == "material" {
		buf.WriteString(MaterialHead())
	}
	buf.WriteString("<style>")
	buf.Write(builtinThemes[theme])
	buf.WriteString("</style>")
	return buf.Bytes()
}

// MessagePageBytes provides the same functionalityt as MessagePage,
// but with []byte instead of string, and without closing </body></html>
func MessagePageBytes(title string, body []byte, theme string) []byte {
	var buf bytes.Buffer
	buf.WriteString("<!doctype html><html><head><title>")
	buf.WriteString(title)
	buf.WriteString("</title>")
	buf.Write(StyleHead(theme))
	buf.WriteString("<head><body><h1>")
	buf.WriteString(title)
	buf.WriteString("</h1>")
	buf.Write(body)
	return buf.Bytes()
}

// SimpleHTMLPage provides a quick way to build a HTML page
func SimpleHTMLPage(title, headline, inhead, body []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString("<!doctype html><html><head><title>")
	buf.Write(title)
	buf.WriteString("</title>")
	buf.Write(inhead)
	buf.WriteString("<head><body><h1>")
	buf.Write(headline)
	buf.WriteString("</h1>")
	buf.Write(body)
	return buf.Bytes()
}

// HTMLLink builds an HTML link given the link text, the URL to a file/directory
// and a boolean that is true if the given URL is to a directory.
func HTMLLink(text, url string, isDirectory bool) string {
	// Add a final slash, if needed
	if isDirectory {
		text += "/"
		url += "/"
	}
	return "<a href=\"/" + url + "\">" + text + "</a><br>"
}

// HighlightHead creates the HTML code for linking with a stylesheet for
// highlighting <code> tags with the given highlight.js code style
func HighlightHead(codeStyle string) string {
	return `<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/` + HighlightJSversion + `/styles/` + codeStyle + `.min.css"><script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/` + HighlightJSversion + `/highlight.min.js"></script><script>hljs.initHighlightingOnLoad();</script>`
}

// MaterialHead creates the HTML code for a material design framework called MUI (www.muicss.com)
func MaterialHead() string {
	return `<link rel="stylesheet" href="//cdn.muicss.com/mui-` + MUIVersion + `/css/mui.min.css"><script src="//cdn.muicss.com/mui-` + MUIVersion + `/js/mui.min.js"></script>`
}

// StyleAmber modifies Amber source code so that a link to the given stylesheet URL is added
func StyleAmber(amberdata *[]byte, url string) {
	// If the given url is not already mentioned and the data contains "body"
	if !bytes.Contains(*amberdata, []byte(url)) && bytes.Contains(*amberdata, []byte("html")) && bytes.Contains(*amberdata, []byte("body")) {
		// Extract one level of indendation
		whitespace := OneLevelOfIndentation(amberdata, "body")
		// Check if there already is a head section
		if bytes.Contains(*amberdata, []byte("head")) {
			// Add a link to the stylesheet
			*amberdata = bytes.Replace(*amberdata, []byte("head\n"), []byte("head\n"+whitespace+whitespace+`link[href="`+url+`"][rel="stylesheet"][type="text/css"]`+"\n"), 1)

		} else if bytes.Contains(*amberdata, []byte("body")) {

			// Add a link to the stylesheet
			*amberdata = bytes.Replace(*amberdata, []byte("html\n"), []byte("html\n"+whitespace+"head\n"+whitespace+whitespace+`link[href="`+url+`"][rel="stylesheet"][type="text/css"]`+"\n"), 1)
		}
	}
}

// StyleHTML modifies HTML source code so that a link to the given stylesheet URL is added
func StyleHTML(htmldata *[]byte, url string) {
	// If the given url is not already mentioned and the data contains "body"
	if !bytes.Contains(*htmldata, []byte(url)) && bytes.Contains(*htmldata, []byte("body")) {
		if bytes.Contains(*htmldata, []byte("</head>")) {
			*htmldata = bytes.Replace(*htmldata, []byte("</head>"), []byte("  <link rel=\"stylesheet\" href=\""+url+"\">\n  </head>"), 1)
		} else if bytes.Contains(*htmldata, []byte("<body>")) {
			*htmldata = bytes.Replace(*htmldata, []byte("<body>"), []byte("  <head>\n  <link rel=\"stylesheet\" href=\""+url+"\">\n  </head>\n  <body>"), 1)
		}
	}
}

// InsertDoctype inserts <doctype html> to the HTML, if missing.
// Does not check if the given data is HTML. Assumes it to be HTML.
func InsertDoctype(htmldata []byte) []byte {
	// If there are more than two lines
	if bytes.Count(htmldata, []byte("\n")) > 2 {
		fields := bytes.SplitN(htmldata, []byte("\n"), 3)
		line1 := strings.ToLower(string(fields[0]))
		line2 := strings.ToLower(string(fields[1]))
		if strings.Contains(line1, "doctype") || strings.Contains(line2, "doctype") {
			return htmldata
		}
		// Doctype is missing from the first two lines, add it
		return []byte("<!doctype html>\n" + string(htmldata))
	}
	return htmldata
}

// NoPage provides the same functionality as NoPage, but returns []byte
func NoPage(filename, theme string) []byte {
	return MessagePageBytes("Not found", []byte("File not found: "+filename), theme)
}