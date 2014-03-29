package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/yvasiyarov/gorelic"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	textTemplate "text/template"
	//"regexp"
	"github.com/russross/blackfriday"
	//"io/ioutil"
	"strings"
	"time"
)

// Language code

var languageTranslation = map[string]string{
	// support
	"afrikaans":          "1",
	"albanian":           "2",
	"arabic":             "3",
	"belarusian":         "4",
	"bulgarian":          "5",
	"catalan":            "6",
	"chineseSimplified":  "7",
	"chineseTraditional": "8",
	"croatian":           "9",
	"czech":              "10",
	"danish":             "11",
	"dutch":              "12",
	"english":            "13",
	"estonian":           "14",
	"filipino":           "15",
	"finnish":            "16",
	"french":             "17",
	"galician":           "18",
	"german":             "19",
	"greek":              "20",
	"haitian":            "21",
	"hebrew":             "22",
	"hindi":              "23",
	"hungarian":          "24",
	"icelandic":          "25",
	"indonesian":         "26",
	"irish":              "27",
	"italian":            "28",
	"japanese":           "29",
	"korean":             "30",
	"latvian":            "31",
	"lithuanian":         "32",
	"macedonian":         "33",
	"malay":              "34",
	"maltese":            "35",
	"norwegian":          "36",
	"persian":            "37",
	"polish":             "38",
	"portuguese":         "39",
	"romanian":           "40",
	"russian":            "41",
	"serbian":            "42",
	"slovak":             "43",
	"slovenian":          "44",
	"spanish":            "45",
	"swahili":            "46",
	"swedish":            "47",
	"thai":               "48",
	"turkish":            "49",
	"ukrainian":          "50",
	"vietnamese":         "51",
	"welsh":              "52",
	"yiddish":            "53",
}

var languageTransliteration = map[string]string{
	// add 54
	"arabic":     "1",
	"armenian":   "2",
	"bengali":    "3",
	"cyrillic":   "4",
	"devanagari": "5",
	"georgian":   "6",
	"greek":      "7",
	"gujarati":   "8",
	"hangul":     "9",
	"hebrew":     "10",
	"hiragana":   "11",
	"jamo":       "12",
	"kannada":    "13",
	"katakana":   "14",
	"latin":      "15",
	"malayalam":  "16",
	"oryia":      "17",
	"syriac":     "18",
	"tamil":      "19",
	"telugu":     "20",
	"thaana":     "21",
	"thai":       "22",
}

type Line struct {
	ScriptureId     int
	Scripture       string
	Hymn            int
	Page            int
	Line            int
	Translation     []string
	Transliteration []string
	Section         string
	MelodyId        int
	AuthorId        int
}

type RequestType int

const ( // iota is reset to 0
	RequestPage    RequestType = iota // c0 == 0
	RequestHymn    RequestType = iota // c1 == 1
	RequestLine    RequestType = iota // c2 == 2
	RequestSection RequestType = iota // c2 == 2
)

// This can be a page reply, shabad reply, etc
type Reply struct {
	Lines []Line
}

type MainPageLine struct {
	Header      string
	Description string
}

type MainPageRequest struct {
	Data []MainPageLine
}

var mainPageRequest MainPageRequest

/*
Set the following environment variables to override corresponding default values:
"GRANTHCO_DATABASE_USERNAME"
"GRANTHCO_DATABASE_PASSWORD"
"GRANTHCO_DATABASE_HOST"
"GRANTHCO_DATABASE_NAME"
*/
var defaultDbUsername = "root" 
var defaultDbPassword = "password"
var defaultDbHost     = "localhost"
var defaultDbPort     = "3306"
var defaultDbName     = "gurbanidb"
	
var databaseConnectionFormatString = "%s:%s@tcp(%s:%s)/%s"
var databaseString = ""

func databaseConnection() (*sql.DB, error) {
	if databaseString == "" {
		panic(fmt.Sprintf("Database string is empty!"))
	}
	return sql.Open("mysql", databaseString)
}

func doQueryBase(c *sql.DB, q string) (*sql.Rows, error) {
	return c.Query(q)
}

func doQuery(q string) (*sql.Rows, error) {
	fmt.Println(q)
	c, err := databaseConnection()
	if err != nil {
		fmt.Println("---- ", err)
		return nil, err
	}
	return doQueryBase(c, q)
}

func queryBuilder(translationLanguages []string,
	transliterationLanguages []string, page string, hymn string,
	lineBegin string, lineEnd string, melodyId string,
	requestType RequestType) string {

	// TODO: check arguments

	base := "scripture.id, scripture.scripture, scripture.section, scripture.hymn, scripture.page, scripture.line, translation.text, transliteration.text, melody.id, author.id, translation.language_id, transliteration.language_id"

	tables := "scripture inner join melody on scripture.melody_id = melody.id inner join author on scripture.author_id = author.id inner join translation on scripture.id = translation.scripture_id inner join transliteration on scripture.id = transliteration.scripture_id"

	var tLangs []string
	for i := 0; i < len(translationLanguages); i++ {
		if _, ok := languageTranslation[translationLanguages[i]]; ok {
			tLangs = append(tLangs, fmt.Sprintf("translation.language_id = %s",
				languageTranslation[translationLanguages[i]]))
		}
	}

	if len(tLangs) == 0 {
		tLangs = append(tLangs, fmt.Sprintf("translation.language_id = %s", languageTranslation["english"]))
	}

	translationBase := "(" + strings.Join(tLangs, " or ") + ") "

	transliterationBase := "(transliteration.language_id = 69) "

	// build up the content, default it to page 1
	request := "scripture.page = 1"

	switch requestType {
	case RequestPage:
		request = fmt.Sprintf("scripture.page = %s", page)
		break
	case RequestLine:
		request = fmt.Sprintf("scripture.id >= %s and scripture.id <= %s",
			lineBegin, lineEnd)
		break
	case RequestHymn:
		request = fmt.Sprintf("scripture.hymn = %s", hymn)
		break
	case RequestSection:
		request = fmt.Sprintf("scripture.melody_id = %s", melodyId)
	}

	query := fmt.Sprintf("select %s from %s where %s and %s and %s order by scripture.id asc;",
		base, tables, translationBase, transliterationBase, request)

	return query
}

func executeQuery(translationLanguages []string,
	transliterationLanguages []string, page string, hymn string,
	lineBegin string, lineEnd string, melodyId string,
	requestType RequestType) (*sql.Rows, error) {

	return doQuery(queryBuilder(translationLanguages, transliterationLanguages, page,
		hymn, lineBegin, lineEnd, melodyId, requestType))
}

func request(translationLanguages []string,
	transliterationLanguages []string, page string, hymn string,
	lineBegin string, lineEnd string, melodyId string,
	requestType RequestType) Reply {

	var reply Reply

	rows, err := executeQuery(translationLanguages, transliterationLanguages,
		page, hymn, lineBegin, lineEnd, melodyId, requestType)

	if err != nil {
		fmt.Println("Error executing the requst", err)
		return reply
	}

	currentId := 0
	var currentTranslationId = 0
	var currentTransliterationId = 0

	for rows.Next() {
		l := new(Line)
		translation := new(string)
		transliteration := new(string)
		var translationId int
		var transliterationId int

		rows.Scan(&l.ScriptureId, &l.Scripture, &l.Section, &l.Hymn, &l.Page, &l.Line,
			&translation, &transliteration, &l.MelodyId, &l.AuthorId,
			&translationId, &transliterationId)

		switch {
		case currentId != l.ScriptureId:
			currentId = l.ScriptureId
			currentTranslationId = translationId
			currentTransliterationId = transliterationId
			l.Translation = append(l.Translation, *translation)
			l.Transliteration = append(l.Transliteration, *transliteration)
			reply.Lines = append(reply.Lines, *l)
			break
		case currentId == l.ScriptureId:
			if currentTranslationId != translationId {
				currentTranslationId = translationId
				reply.Lines[len(reply.Lines)-1].Translation =
					append(reply.Lines[len(reply.Lines)-1].Translation, *translation)
			}
			if currentTransliterationId != transliterationId {
				currentTransliterationId = transliterationId
				reply.Lines[len(reply.Lines)-1].Transliteration =
					append(reply.Lines[len(reply.Lines)-1].Transliteration, *transliteration)
			}
		}
	}

	return reply
}

func reduceReply(r Reply) Reply {

	toReturn := make([]Line, r.Lines[len(r.Lines)-1].Line)

	for _, line := range r.Lines {
		var index = line.Line - 1
		toReturn[index].Scripture += " " + line.Scripture
		toReturn[index].Page = line.Page
		toReturn[index].Hymn = line.Hymn
		toReturn[index].ScriptureId = line.ScriptureId
		for i, translation := range line.Translation {
			if len(toReturn[index].Translation) == 0 {
				toReturn[index].Translation = make([]string, len(line.Translation))
			}
			toReturn[index].Translation[i] += " " + translation
		}
		for i, transliterate := range line.Transliteration {
			if len(toReturn[index].Transliteration) == 0 {
				toReturn[index].Transliteration = make([]string, len(line.Transliteration))
			}
			toReturn[index].Transliteration[i] += " " + transliterate
		}
	}

	return Reply{toReturn}
}

func requestPage(translationLanguages []string,
	transliterationLanguages []string, page string) Reply {
	return request(translationLanguages, transliterationLanguages, page, "0",
		"0", "0", "0", RequestPage)
}

func requestReducedPage(translationLanguages []string,
	transliterationLanguages []string, page string) Reply {
	return reduceReply(requestPage(translationLanguages, transliterationLanguages, page))
}

func requestLines(translationLanguages []string,
	transliterationLanguages []string, begin string, end string) Reply {
	return request(translationLanguages, transliterationLanguages, "0", "0",
		begin, end, "0", RequestLine)
}

func requestHymn(translationLanguages []string,
	transliterationLanguages []string, hymn string) Reply {
	return request(translationLanguages, transliterationLanguages, "0", hymn,
		"0", "0", "0", RequestHymn)
}

func requestReducedHymn(translationLanguages []string,
	transliterationLanguages []string, hymn string) Reply {
	return reduceReply(requestHymn(translationLanguages, transliterationLanguages, hymn))
}

func requestSection(translationLanguages []string,
	transliterationLanguages []string, section string) Reply {
	return request(translationLanguages, transliterationLanguages, "0", "0",
		"0", "0", section, RequestSection)
}

func printReply(res http.ResponseWriter, r Reply) {
	for _, l := range r.Lines {
		fmt.Fprintln(res, l)
	}
}

func returnJsonReply(res http.ResponseWriter, r Reply) {
	b, err := json.Marshal(r)

	if err != nil {
		fmt.Fprintln(res, "got an error", err)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(b)
}

func returnHtmlreply(res http.ResponseWriter, template string, r Reply) {

	renderReplyTemplate(res, template, &r)
}

func returnHtmlMain(w http.ResponseWriter, template string) {

	err := textTemplate.Must(textTemplate.New("base.html").
		Delims("<<<", ">>>").
		ParseFiles(templateLocationPrefix+"base.html",
		templateLocationPrefix+"main_page.html")).ExecuteTemplate(w, "base", mainPageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//
// Handler Bases
//

func prepareLanguageString(queryString url.Values) []string {
	var splitLangs []string

	if _, ok := queryString["lang"]; ok {
		splitLangs = strings.Split(queryString["lang"][0], ",")
	}

	if len(splitLangs) == 0 {
		splitLangs = []string{"english"}
	}
	return splitLangs
}

func mainHandlerBase(req *http.Request) Reply {
	var r Reply
	return r
}

func hymnHandlerBase(req *http.Request) Reply {
	vars := mux.Vars(req)
	hymn := vars["hymn"]
	fmt.Println(hymn)
	splitLangs := prepareLanguageString(req.URL.Query())

	r := requestHymn(splitLangs, []string{}, hymn)
	return r
}

func reducedHymnHandlerBase(req *http.Request) Reply {
	vars := mux.Vars(req)
	hymn := vars["hymn"]
	fmt.Println(hymn)
	splitLangs := prepareLanguageString(req.URL.Query())

	r := requestReducedHymn(splitLangs, []string{}, hymn)
	return r
}

func lineHandlerBase(req *http.Request) Reply {
	vars := mux.Vars(req)
	begin := vars["begin"]
	end := vars["end"]
	fmt.Println(begin, end)

	splitLangs := prepareLanguageString(req.URL.Query())
	r := requestLines(splitLangs, []string{}, begin, end)
	return r
}

func pageHandlerBase(req *http.Request) Reply {
	vars := mux.Vars(req)
	page := vars["page"]
	fmt.Println(page)
	splitLangs := prepareLanguageString(req.URL.Query())
	//r := requestReducedPage(splitLangs, []string{}, page)
	r := requestPage(splitLangs, []string{}, page)
	//fmt.Println(r)
	return r
}

func reducedPageHandlerBase(req *http.Request) Reply {
	vars := mux.Vars(req)
	page := vars["page"]
	fmt.Println(page)
	splitLangs := prepareLanguageString(req.URL.Query())
	r := requestReducedPage(splitLangs, []string{}, page)
	//r := requestPage(splitLangs, []string{}, page)
	//fmt.Println(r)
	return r
}

//
// Handler Rest
//
func restMainHandler(res http.ResponseWriter, req *http.Request) {
	returnJsonReply(res, mainHandlerBase(req))
}

func restHymnHandler(res http.ResponseWriter, req *http.Request) {
	returnJsonReply(res, hymnHandlerBase(req))
}

func restLineHandler(res http.ResponseWriter, req *http.Request) {
	returnJsonReply(res, lineHandlerBase(req))
}

func restPageHandler(res http.ResponseWriter, req *http.Request) {
	returnJsonReply(res, pageHandlerBase(req))
}

//
// Handler HTML
//
func mainHandler(res http.ResponseWriter, req *http.Request) {
	returnHtmlMain(res, "main_display")
}

func hymnHandler(res http.ResponseWriter, req *http.Request) {
	returnHtmlreply(res, "reply_display", hymnHandlerBase(req))
}

func reducedHymnHander(res http.ResponseWriter, req *http.Request) {
	returnHtmlreply(res, "reply_display", reducedHymnHandlerBase(req))
}

func lineHandler(res http.ResponseWriter, req *http.Request) {
	returnHtmlreply(res, "reply_display", lineHandlerBase(req))
}

func pageHandler(res http.ResponseWriter, req *http.Request) {
	returnHtmlreply(res, "reply_display", pageHandlerBase(req))
}

func reducedPageHandler(res http.ResponseWriter, req *http.Request) {
	returnHtmlreply(res, "reply_display", reducedPageHandlerBase(req))
}

func aboutHandler(res http.ResponseWriter, req *http.Request) {
	returnHtmlreply(res, "about_display", mainHandlerBase(req))
}

//
// Handler HTML
//

func randomHymnHandler(res http.ResponseWriter, req *http.Request) {
	rand.Seed(time.Now().UnixNano())
	var hymn = rand.Intn(3621) // Number of shabads
	fmt.Println(hymn)

	splitLangs := prepareLanguageString(req.URL.Query())

	r := requestHymn(splitLangs, []string{}, strconv.Itoa(hymn))

	returnHtmlreply(res, "reply_display", r)
}

func randomLineHandler(res http.ResponseWriter, req *http.Request) {
	returnHtmlreply(res, "reply_display", lineHandlerBase(req))
}

func randomPageHandler(res http.ResponseWriter, req *http.Request) {
	returnHtmlreply(res, "reply_display", pageHandlerBase(req))
}

//
// -----------------------------------------------------------------------------
//

var templateLocationPrefix = "static/templates/html/"

var templates = map[string]*template.Template{

	"reply_display": template.Must(template.New("base.html").
		Delims("<<<", ">>>").
		ParseFiles(templateLocationPrefix+"base.html",
		templateLocationPrefix+"reply_display.html")),

	"main_display": template.Must(template.New("base.html").
		Delims("<<<", ">>>").
		ParseFiles(templateLocationPrefix+"base.html",
		templateLocationPrefix+"main_page.html")),

	"about_display": template.Must(template.New("base.html").
		Delims("<<<", ">>>").
		ParseFiles(templateLocationPrefix+"base.html",
		templateLocationPrefix+"about_page.html")),
}

func renderReplyTemplate(w http.ResponseWriter, tmpl string, s *Reply) {

	err := templates[tmpl].ExecuteTemplate(w, "base", s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//
// -----------------------------------------------------------------------------
//

func setupRoutes() *mux.Router {

	// TODO: Rethink routes
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	r.HandleFunc("/h{hymn:[0-9]+}", hymnHandler)
	r.HandleFunc("/{page:[0-9]+}", pageHandler)
	r.HandleFunc("/full/{page:[0-9]+}", reducedPageHandler)
	r.HandleFunc("/full/h{hymn:[0-9]+}", reducedHymnHander)
	r.HandleFunc("/{begin:[0-9]+}-{end:[0-9]+}", lineHandler)
	r.HandleFunc("/randompage", randomPageHandler)
	r.HandleFunc("/about", aboutHandler)
	r.HandleFunc("/randomhymn", randomHymnHandler)
	r.HandleFunc("/randomline", randomLineHandler)

	restr := r.PathPrefix("/rest").Subrouter()
	restr.HandleFunc("/", restMainHandler)
	restr.HandleFunc("/h{hymn:[0-9]+}", restHymnHandler)
	restr.HandleFunc("/{page:[0-9]+}", restPageHandler)
	restr.HandleFunc("/{begin:[0-9]+}-{end:[0-9]+}", restLineHandler)

	r.PathPrefix("/static").Handler(http.StripPrefix("/static",
		http.FileServer(http.Dir("./static/"))))
	return r
}

func runLogProcesses() {
	fmt.Println("Initialize relic ")
	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicLicense = os.Getenv("GRANTHCO_LOG_CREDENTIALS")
	agent.Run()

}

func setEnvironment() {

	// Get the environment
	databaseUsername := os.Getenv("GRANTHCO_DATABASE_USERNAME")
	databasePassword := os.Getenv("GRANTHCO_DATABASE_PASSWORD")
	databaseHost := os.Getenv("GRANTHCO_DATABASE_HOST")
	databasePort := os.Getenv("GRANTHCO_DATABASE_PORT")
	databaseName := os.Getenv("GRANTHCO_DATABASE_NAME")

	if databaseUsername == "" {
		databaseUsername = defaultDbUsername
	}
	if databasePassword == "" {
		databasePassword = defaultDbPassword
	}
	if databaseHost == "" {
		databaseHost = defaultDbHost
	}
	if databasePort == "" {
		databasePort = defaultDbPort
	}
	if databaseName == "" {
		databaseName = defaultDbName
	}

	databaseString = fmt.Sprintf(databaseConnectionFormatString, databaseUsername, databasePassword, 
		databaseHost, databasePort, databaseName)

	//databaseString = localDatabaseString
	fmt.Println(databaseString)
}

func readMainPageData() {

	//var reply MainPageRsequest
	var pageLine MainPageLine

	file, err := os.Open("data/mainpage.markup")
	if err != nil {
		fmt.Println("Big problem reading file", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		begin := strings.Index(line, "|")

		if begin == 0 {

			if pageLine.Header != "" {
				mainPageRequest.Data = append(mainPageRequest.Data, pageLine)
			}

			pageLine.Header = strings.TrimSpace(line[2:])
			pageLine.Description = ""

		} else {
			if len(line) > 0 {
				pageLine.Description = strings.Join([]string{pageLine.Description,
					string(blackfriday.MarkdownBasic(scanner.Bytes()))}, " ")
			}
		}
	}

	//fmt.Println(mainPageRequest)
	//data, err := ioutil.ReadFile("data/mainpage.markup")

	//output := blackfriday.MarkdownBasic(data)
	//fmt.Println("===========", string(output))

}

func main() {

	setEnvironment()

	http.Handle("/", setupRoutes())

	runLogProcesses()

	readMainPageData()

	fmt.Println("listening...")
	
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		http.ListenAndServe(":8888", nil)
	}

}
