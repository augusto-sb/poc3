package main

// imports

import (
	"context"
	"crypto/md5"
	"errors"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// types

type session struct {
	security   string;
	timestamp  int64;
	info       map[string]any;
}

type user struct {
	name string;
	password string;
}

// consts

const cookieName string = "PHPSESSID"; //gg
const cookieDurationInSeconds = 3600;
const cleanerInterval = 30;
const ctxKey string = "yourKey";

// vars

var muS sync.Mutex = sync.Mutex{};
var muU sync.Mutex = sync.Mutex{};
var sessions map[string]session = map[string]session{};
var users []user = []user{
	user{
		name: "admin",
		password: "admin",
	},
};

// helpers

func logger(msj string){
	if(os.Getenv("LOGGER")=="true"){
		fmt.Println(msj);
	}
}

/*func getSession(req *http.Request) *session{ // mmmmm
	cookie, err := req.Cookie(cookieName)
	if (err != nil){
		return nil;
	}
	muS.Lock();
	val, ok := sessions[cookie.Value];
	muS.Unlock();
	if(!ok){
		return nil;
	}
	return &val;
}*/

//gracefull shutdown implementar!!!

func setCleaner(timeSec uint) () {
	//timer cada tanto limpia sessions vencidas!
	ticker := time.NewTicker(time.Duration(timeSec * 1000) * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				logger("running cleaner!")
				muS.Lock();
				for k, v := range sessions {
					if(v.timestamp + (cookieDurationInSeconds * 1000) < time.Now().UnixMilli()){
						logger("cleaner found expired: "+k)
						delete(sessions, k);
					}
				}
				muS.Unlock();
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
}

// middlewares

func sessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var shouldSetCookie bool = false;
		cookie, err := req.Cookie(cookieName)
		if err != nil {
			if(errors.Is(err, http.ErrNoCookie)){
				shouldSetCookie = true;
			}else{
				logger("error retrieving cookie")
				http.Error(rw, "server error", http.StatusInternalServerError);
				return;
			}
		}
		//get info de origen para que no me la puedan usar de otro lado!
		//el proxy deberia pasarme el header o algo de info! como ip publica etc
		//logger(req.Header.Values("User-Agent")) //recibir en x-forwarded-for por ej! para bien seguro
		pseudoSecure := req.Header.Get("User-Agent") + req.Header.Get("Accept") + req.Header.Get("Host") + req.Header.Get("X-Forwarded-For") + req.Header.Get("Forwarded");
		logger(pseudoSecure)
		if(!shouldSetCookie){
			//check validity
			muS.Lock();
			val, ok := sessions[cookie.Value];
			muS.Unlock();
			if(ok){
				//rw.Write([]byte("cookie found in session and is:" + cookie.Value+"\n"))
				if(val.timestamp + (cookieDurationInSeconds * 1000) < time.Now().UnixMilli()){
					//rw.Write([]byte("expired session!\n")) //shouldSetCookie ???
					muS.Lock();
					delete(sessions, cookie.Value)
					muS.Unlock();
					shouldSetCookie = true;
				}else{
					//actualizar timestamp
					val.timestamp = time.Now().UnixMilli()
				}
				if(val.security != pseudoSecure){
					rw.Write([]byte("hacker wtf\n"))
					return;
				}
				if(len(val.info) == 0){
					logger("not logged in")
				}else{
					logger("logged in")
				}
			}else{
				//rw.Write([]byte("cookie must have expired found in session!\n"))
				shouldSetCookie = true
			}
		}
		if(shouldSetCookie){
			h := md5.New()
			now := time.Now()
			io.WriteString(h, now.String())
			io.WriteString(h, pseudoSecure)
			cookieValueAndSessionKey := hex.EncodeToString(h.Sum(nil));
			cookie = &http.Cookie{
				Name:     cookieName,
				//Domain:   "localhost:3000",
				Value:    cookieValueAndSessionKey,
				MaxAge:   cookieDurationInSeconds,
				HttpOnly: true,
				Secure:   true, //CONFIGURABLE EN PROD, CORS TAMBIEN
				SameSite: http.SameSiteLaxMode, //http.SameSiteNoneMode,
			}
			newSession := session{
				timestamp: now.UnixMilli(),
				security: pseudoSecure,
			}
			muS.Lock();
			sessions[cookieValueAndSessionKey] = newSession;
			muS.Unlock();
			//rw.Write([]byte("cookie not send, setting one!"))
		}
		http.SetCookie(rw, cookie)
		ctx := context.WithValue(req.Context(), ctxKey, cookie.Value)
        req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	var corsOrigin string = os.Getenv("CORS_ORIGIN");
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if(corsOrigin != ""){
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
			rw.Header().Set("Access-Control-Allow-Headers", "authorization")
			rw.Header().Set("Access-Control-Allow-Methods", "GET,POST")
		}
		if(req.Method == http.MethodOptions){
			return;
		}
		next.ServeHTTP(rw, req)
	})
}

func middleware(next http.HandlerFunc) http.HandlerFunc {
	// chain de todos los middlewares
	return corsMiddleware(sessionMiddleware(next));
}

// handlers

func sessionHandler(rw http.ResponseWriter, req *http.Request) {
	/*var sess *session = getSession(req);
	if(sess == nil){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}*/
	cookieValueAndSessionKey, ok := req.Context().Value(ctxKey).(string)
	if(!ok){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	muS.Lock();
	sess, ok := sessions[cookieValueAndSessionKey];
	muS.Unlock();
	if(!ok){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	byteArr, err := json.Marshal(len(sess.info));
	if(err != nil){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	rw.Write(byteArr)
}

func loginHandler(rw http.ResponseWriter, req *http.Request) {
	cookieValueAndSessionKey, ok := req.Context().Value(ctxKey).(string)
	if(!ok){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	u, p, ok := req.BasicAuth();
	if(!ok){
		rw.WriteHeader(http.StatusBadRequest);
		return;
	}
	validUserPass := false;
	muU.Lock();
	for _, val := range users {
		if(val.name == u && val.password == p){
			validUserPass = true;
			break;
		}
	}
	muU.Unlock();
	if(validUserPass){
		muS.Lock();
		sess, ok := sessions[cookieValueAndSessionKey];
		if(!ok){
			rw.WriteHeader(http.StatusBadRequest);
		}else{
			sess.info = map[string]any{"asd":"asd"};
			rw.WriteHeader(http.StatusOK)
		}
		muS.Unlock();
	}else{
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func logoutHandler(rw http.ResponseWriter, req *http.Request) {
	cookieValueAndSessionKey, ok := req.Context().Value(ctxKey).(string)
	if(!ok){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	_, ok = sessions[cookieValueAndSessionKey];
	if(!ok){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	delete(sessions, cookieValueAndSessionKey)
}

// main

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware(http.NotFound))
	mux.HandleFunc("GET /entities", middleware(http.NotFound))
	mux.HandleFunc("GET /session", middleware(sessionHandler))
	mux.HandleFunc("POST /login", middleware(loginHandler))
	mux.HandleFunc("POST /logout", middleware(logoutHandler))
	setCleaner(cleanerInterval);
	err := http.ListenAndServe("0.0.0.0:3000", mux)
	if err != nil {
		panic(err)
	}
}