package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/youtube/v3"

	"github.com/vim345/plailyst/crawler"
	"github.com/vim345/plailyst/oauth"
	"gopkg.in/yaml.v3"
)

type wrapperStruct struct {
	oauth   *oauth.LocalOauth
	Service *youtube.Service
}

func newWrapperStruct(oauth *oauth.LocalOauth) *wrapperStruct {
	return &wrapperStruct{oauth: oauth}
}

func readYaml(path string) *crawler.Configs {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Cannot find the configs file = %+v\n", err)
		os.Exit(1)
	}
	data := &crawler.Configs{}

	err = yaml.Unmarshal(yfile, data)

	if err != nil {
		log.Errorf("Configs is not valid = %+v\n", err)
		os.Exit(1)
	}

	return data
}

func setupLogger() {
	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash("./logs/plailyst/plailist.log"),
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = time.RFC1123Z // or RFC3339
	logFormatter.FullTimestamp = true

	log.SetFormatter(logFormatter)
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)
}

func init() {
	setupLogger()
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "creds.json")
}

func (ws *wrapperStruct) callBack(w http.ResponseWriter, r *http.Request) {
	// url := "http://localhost:8765/login/callback?state=rCRxv20Odp7DeEaOgxHRyg%3D%3D&code=4/0AWgavdc117ValCjz7TL1KIAYUJ8P-jLj4_ESLy6w6WKMEBX51qvO2sY6CXEkTZbpu2PQQw&scope=https://www.googleapis.com/auth/youtube"
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	service, err := ws.oauth.GetService(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	ws.Service = service
	log.Println("Redirecting to the main page")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (ws *wrapperStruct) handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Service: %s\n", ws.Service.BasePath)
	if ws.Service == nil {
		fmt.Fprintf(w, "Go to login page to proceed!")
	} else {
		configs := readYaml("configs.yaml")
		c := crawler.NewCrawler(configs, ws.Service)
		c.Run()
		fmt.Fprintf(w, "Updated the playist. Have fun!")
	}
}

func (ws *wrapperStruct) login(w http.ResponseWriter, r *http.Request) {
	ws.oauth.Login(w, r)
}

func main() {
	server := &http.Server{
		Addr: fmt.Sprintf(":8765"),
	}
	handlers := newWrapperStruct(oauth.NewOauth())
	http.HandleFunc("/", handlers.handler)
	http.HandleFunc("/login/", handlers.login)
	http.HandleFunc("/login/callback", handlers.callBack)

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}
