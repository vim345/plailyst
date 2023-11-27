package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	yfile, err := os.ReadFile(path)
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
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/run", http.StatusTemporaryRedirect)
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

func (ws *wrapperStruct) main(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Great job out of you!")
}

func (ws *wrapperStruct) handler(w http.ResponseWriter, r *http.Request) {

	if ws.Service == nil {
		fmt.Fprintf(w, "Manually log in")
	} else {
		configsFile := os.Getenv("CONFIGS")
		var configs *crawler.Configs
		if len(configsFile) == 0 {
			configs = readYaml("configs.yaml")
		} else {
			configs = readYaml(configsFile)
		}
		c := crawler.NewCrawler(configs, ws.Service)
		c.Run()
		fmt.Fprintf(w, "Updated videos = %+v\n", c.GetUpdatedVideos())
		fmt.Fprintf(w, "Updated the playist. Have fun!")
	}
}

func (ws *wrapperStruct) login(w http.ResponseWriter, r *http.Request) {
	ws.oauth.Login(w, r)
}

func main() {
	handlers := newWrapperStruct(oauth.NewOauth())
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/testing", fs)
	http.HandleFunc("/run", handlers.handler)
	http.HandleFunc("/login/", handlers.login)
	http.HandleFunc("/login/callback", handlers.callBack)
	http.HandleFunc("/", handlers.main)

	cert := os.Getenv("CERT")
	privateKey := os.Getenv("PRIVATE_KEY")
	if len(privateKey) == 0 || len(cert) == 0 {
		log.Printf("Starting HTTP Server. Listening at 8765")
		if err := http.ListenAndServe(":8765", nil); err != http.ErrServerClosed {
			log.Printf("%v", err)
		} else {
			log.Println("Server closed!")
		}
	} else {
		log.Printf("Starting HTTPS Server. Listening at 8443")
		if err := http.ListenAndServeTLS(":8443", cert, privateKey, nil); err != http.ErrServerClosed {
			log.Printf("%v", err)
		} else {
			log.Println("Server closed!")
		}
	}
}
