package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
	"net/http"
)

func Off(w http.ResponseWriter, r *http.Request) {
	Database.LogOff()
	http.Redirect(w,r,"/",http.StatusPermanentRedirect)
}