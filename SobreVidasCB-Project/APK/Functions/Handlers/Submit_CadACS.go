package Handlers

import "net/http"

func submitcadacs(w http.ResponseWriter, r *http.Request){
	ProcessCadACS(r)
	http.Redirect(w,r,"/",http.StatusPermanentRedirect)
}