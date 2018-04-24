package counter

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	path        = os.Getenv("GOPATH")
	search_page = path + "/src/Wordcounter/assets/search.html"
)

// web portal
func Portal(w http.ResponseWriter, r *http.Request) {

	input_tpl := template.Must(template.ParseFiles(search_page))
	input_tpl.Execute(w, struct{ Success bool }{true})

}

//this section of the code will strip all the htnl tags of the webpage and returns the wordount
func Wordcounter(w http.ResponseWriter, r *http.Request) {
	contents := make([]string, 0)
	wordcount := make(map[string]int)
	result := make([]string, 0)
	defer r.Body.Close()
	u, _ := ioutil.ReadAll(r.Body)
	url := string(u)
	if !strings.Contains(url, "http") {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if reflect.ValueOf(resp).IsNil() {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Invalid URL"))
		return
	}
	if err != nil {
		fmt.Printf("got error %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("could not access the URL"))
		return
	}

	defer resp.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	z := html.NewTokenizerFragment(resp.Body, "body")
	stt := z.Token()

loop:
	for {
		switch tt := z.Next(); tt {
		case html.ErrorToken:
			break loop
		case html.StartTagToken:
			stt = z.Token()
		case html.TextToken:
			if contains([]string{"script", "style"}, stt.Data) {
				continue
			}
			content := strings.TrimSpace(html.UnescapeString(string(z.Text())))
			if len(content) > 0 {
				regSpecChar := regexp.MustCompile(`[^a-zA-Z'\s-]+`)
				content = regSpecChar.ReplaceAllString(content, "")
				regSpaceDash := regexp.MustCompile(`[\s\p{Zs}]{2,}|-|\n+`)
				content = regSpaceDash.ReplaceAllString(content, " ")
				contentSlice := strings.Split(content, " ")
				contents = append(contents, contentSlice...)
			}
		}
	}
	for _, data := range contents {
		words := strings.Fields(data)
		for _, word := range words {
			wordcount[word]++
		}
	}
	for a, count := range wordcount {
		result = append(result, a+":"+strconv.Itoa(count))
	}

	ReturnJSONAPISuccess(w, result)
}

func contains(strSlice []string, s string) bool {
	for _, value := range strSlice {
		if value == s {
			return true
		}
	}
	return false
}

// return success output
func ReturnJSONAPISuccess(w http.ResponseWriter, extra interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(struct {
		Status bool
		Msg    string
		Extra  interface{}
	}{Msg: "success", Status: true, Extra: extra})
	fmt.Println(string(j))
	w.Write(j)
}
