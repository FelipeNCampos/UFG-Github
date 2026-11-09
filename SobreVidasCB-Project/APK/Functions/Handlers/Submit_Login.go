package Handlers

import "net/http"

func SubmitLogin(w http.ResponseWriter, r* http.Request){
	ProcessLogin(r)
	http.Redirect(w,r,"/home",http.StatusPermanentRedirect)
}