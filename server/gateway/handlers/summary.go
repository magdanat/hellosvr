package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pageURL := r.FormValue("url")
	if len(pageURL) < 1 {
		http.Error(w, "Invalid input params!", http.StatusBadRequest)
		return
	}

	htmlStream, err := fetchHTML(pageURL)
	if err != nil {
		http.Error(w, "Error! Bad Request.", 400)
		return
	}
	pageSum, err := extractSummary(pageURL, htmlStream)
	if err != nil {
		http.Error(w, "Server-generated an unexpected error, please try again!", 500)
		return
	}
	defer htmlStream.Close()

	w.Header().Set("Content-Type", "application/json")

	var jsonData []byte
	jsonData, err = json.Marshal(pageSum)
	if err != nil {
		http.Error(w, "Invalid JSON", 500)
		return
	}
	w.Write(jsonData)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.
	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML
	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/

	// check response status code
	resp, err := http.Get(pageURL)
	if resp.StatusCode >= 400 || err != nil {
		return nil, errors.New("error")
	}
	// check response content type
	ctype := resp.Header.Get("Content-Type")
	if strings.HasPrefix(ctype, "text/html") {
		return resp.Body, nil
	}
	return nil, errors.New("error")
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	pageSum := new(PageSummary)
	tokenizer := html.NewTokenizer(htmlStream)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}

			// process the token according to the token type...
		} else if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()

			if token.Data == "title" {
				if tokenizer.Next() == html.TextToken {
					titleToken := tokenizer.Token()
					if pageSum.Title == "" {
						pageSum.Title = titleToken.Data
					}
				}
			}

			if token.Data == "meta" {
				key := ""
				prop := ""
				name := ""
				content := ""
				for _, attr := range token.Attr {
					if attr.Key == "property" {
						prop = attr.Val
					}
					if attr.Key == "name" {
						name = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}

				// Determine property type
				if prop == "og:type" {
					key = "type"
				} else if prop == "og:url" {
					key = "url"
				} else if prop == "og:site_name" {
					key = "site_name"
				} else if prop == "og:title" {
					key = "title"
				} else if prop == "og:description" {
					key = "description"
				} else if prop == "og:image" {
					// Need to handle images differently
					// Need to create a PreviewImage struct for each
					// OpenGraph image on the page
					pageSum.Images = append(pageSum.Images, new(PreviewImage))
					key = "image url"
				} else if prop == "og:image:secure_url" {
					key = "image surl"
				} else if prop == "og:image:type" {
					key = "image type"
				} else if prop == "og:image:width" {
					key = "image width"
				} else if prop == "og:image:height" {
					key = "image height"
				} else if prop == "og:image:alt" {
					key = "image alt"
				}

				if name == "author" {
					key = "author"
				} else if name == "keywords" {
					key = "keywords"
				} else if name == "description" && len(pageSum.Description) == 0 {
					key = "description"
				}

				// Set values based on key
				if key != "" {
					if key == "type" {
						pageSum.Type = content
					} else if key == "url" {
						pageSum.URL = content
					} else if key == "site_name" {
						pageSum.SiteName = content
					} else if key == "title" {
						if content != "" {
							pageSum.Title = content
						}
					} else if key == "author" {
						pageSum.Author = content
					} else if key == "keywords" {
						s := strings.Split(content, ",")
						for i := range s {
							s[i] = strings.TrimSpace(s[i])
						}
						pageSum.Keywords = s
					} else if key == "description" {
						pageSum.Description = content
					} else if key == "image url" {
						base := pageURL
						baseURL, _ := url.Parse(base)
						path, _ := url.Parse(content)
						url := baseURL.ResolveReference(path)
						pageSum.Images[len(pageSum.Images)-1].URL = url.String()
					} else if key == "image surl" {
						baseSecure := pageURL
						baseURLSecure, _ := url.Parse(baseSecure)
						pathSecure, _ := url.Parse(content)
						urlSecure := baseURLSecure.ResolveReference(pathSecure)
						pageSum.Images[len(pageSum.Images)-1].SecureURL = urlSecure.String()
					} else if key == "image type" {
						pageSum.Images[len(pageSum.Images)-1].Type = content
					} else if key == "image width" {
						if n, err := strconv.Atoi(content); err == nil {
							pageSum.Images[len(pageSum.Images)-1].Width = n
						}
					} else if key == "image height" {
						if n, err := strconv.Atoi(content); err == nil {
							pageSum.Images[len(pageSum.Images)-1].Height = n
						}
					} else if key == "image alt" {
						pageSum.Images[len(pageSum.Images)-1].Alt = content
					}
					key = ""
				}
			}

			if token.Data == "link" {
				pageSum.Icon = new(PreviewImage)
				linkKey := ""
				for _, attr := range token.Attr {
					if attr.Key == "rel" && attr.Val == "icon" {
						linkKey = "yes"
					}
					if linkKey == "yes" {
						if attr.Key == "href" {
							baseIcon := pageURL
							baseURLIcon, _ := url.Parse(baseIcon)
							pathIcon, _ := url.Parse(attr.Val)
							urlIcon := baseURLIcon.ResolveReference(pathIcon)
							pageSum.Icon.URL = urlIcon.String()
						} else if attr.Key == "type" {
							pageSum.Icon.Type = attr.Val
						} else if attr.Key == "sizes" {
							size := attr.Val
							if strings.Contains(size, "x") {
								heightWidth := strings.Split(size, "x")
								height, err := strconv.Atoi(heightWidth[0])
								if err == nil {
									pageSum.Icon.Height = height
								}
								width, err := strconv.Atoi(heightWidth[1])
								if err == nil {
									pageSum.Icon.Width = width
								}
							}
						} else if attr.Key == "alt" {
							pageSum.Icon.Alt = attr.Val
						}
					}
				}
			}
		} else if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "head" {
				break
			}
		}
	}
	return pageSum, nil
}
