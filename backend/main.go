package main



import (
	"crypto/md5"
	"errors"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)



type session struct {
	security   string;
	timestamp  int64;
	info       map[string]any;
}



type user struct {
	name string;
	password string;
}



func middleware(next http.HandlerFunc) http.HandlerFunc {
	var corsOrigin string = os.Getenv("CORS_ORIGIN");
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if(corsOrigin != ""){
			//rw.Header().Set("Access-Control-Allow-Credentials", "true")
			rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
			//rw.Header().Set("Access-Control-Allow-Headers", "content-type")
			//rw.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
			rw.Header().Set("Access-Control-Allow-Methods", "GET,POST")
		}
		next.ServeHTTP(rw, req)
	})
}



const cookieName string = "PHPSESSID"; //gg
const cookieDurationInSeconds = 3600;
const cleanerInterval = 30;
var mu sync.Mutex = sync.Mutex{};
var sessions map[string]session = map[string]session{};
var users []user = []user{
	user{
		name: "admin",
		password: "admin",
	},
};



//timer cada tanto limpie sessions vencidas!
func setCleaner(timeSec uint) () {
	ticker := time.NewTicker(time.Duration(timeSec * 1000) * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("running cleaner!")
				mu.Lock();
				for k, v := range sessions {
					if(v.timestamp + (cookieDurationInSeconds * 1000) < time.Now().UnixMilli()){
						fmt.Println("cleaner found expired: "+k)
						delete(sessions, k);
					}
				}
				mu.Unlock();
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
}



//gracefull shutdown

func getSession() *session{
	asd
}

func sessionHandler(respW http.ResponseWriter, req *http.Request) {
	asd
}
func loginHandler(respW http.ResponseWriter, req *http.Request) {
	var sess *session = getSession();
	u, p, ok := req.BasicAuth();
	if(!ok){
		respW.WriteHeader(http.StatusBadRequest);
		return;
	}
	logged := false;
	mu.Lock();
	for _, val := range users {
		if(val.name == u && val.password == p){
			logged = true;
			break;
		}
	}
	mu.Unlock();
}
func logoutHandler(respW http.ResponseWriter, req *http.Request) {
	asd
}



func genericCookieHandler(respW http.ResponseWriter, req *http.Request) {
	var shouldSetCookie bool = false;
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			fmt.Println("cookie not send!");
			shouldSetCookie = true;
		default:
			http.Error(respW, "server error", http.StatusInternalServerError);
			return;
		}
	}
	//get info de origen para que no me la puedan usar de otro lado!
	//el proxy deberia pasarme el header o algo de info! como ip publica etc
	pseudoSecure := req.Header.Get("User-Agent") + req.Header.Get("Accept") + req.Header.Get("Host") + req.Header.Get("X-Forwarded-For") + req.Header.Get("Forwarded");
	//fmt.Println(req.Header.Values("User-Agent")) //recibir en x-forwarded-for por ej! para bien seguro
	if(shouldSetCookie){
		h := md5.New()
		now := time.Now()
		milliseconds := now.UnixMilli()
		io.WriteString(h, now.String())
		io.WriteString(h, pseudoSecure)
		cookieValueAndSessionKey := hex.EncodeToString(h.Sum(nil));
		newCookie := http.Cookie{
			Name:     cookieName,
			//Domain:   "localhost:3000",
			Value:    cookieValueAndSessionKey,
			MaxAge:   cookieDurationInSeconds,
			HttpOnly: true,
			Secure:   true, //CONFIGURABLE EN PROD, CORS TAMBIEN
			SameSite: http.SameSiteNoneMode, //http.SameSiteLaxMode,
		}
		newSession := session{
			timestamp: milliseconds,
			security: pseudoSecure,
		}
		mu.Lock();
		sessions[cookieValueAndSessionKey] = newSession;
		mu.Unlock();
		http.SetCookie(respW, &newCookie)
		respW.Write([]byte("cookie not send, setting one!"))
	}else{
		//check validity
		mu.Lock();
		val, ok := sessions[cookie.Value];
		mu.Unlock();
		if(ok){
			respW.Write([]byte("cookie found in session and is:" + cookie.Value+"\n"))
			if(val.timestamp + (cookieDurationInSeconds * 1000) < time.Now().UnixMilli()){
				respW.Write([]byte("expired session!\n")) //shouldSetCookie ???
			}
			if(val.security != pseudoSecure){
				respW.Write([]byte("hacker wtf\n"))
			}
			if(len(val.info) == 0){
				fmt.Println("not logged in")
			}else{
				fmt.Println("logged in")
			}
		}else{
			respW.Write([]byte("cookie must have expired found in session!\n"))
		}
	}
}



func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware(http.NotFound))
	mux.HandleFunc("GET /session", middleware(sessionHandler))
	mux.HandleFunc("POST /login", middleware(loginHandler))
	mux.HandleFunc("POST /logout", middleware(logoutHandler))
	setCleaner(cleanerInterval);
	err := http.ListenAndServe("0.0.0.0:3000", mux)
	if err != nil {
		panic(err)
	}
}