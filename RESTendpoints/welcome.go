package RESTendpoints

import (
	"csn/models"
	"database/sql"
	"fmt"
	"net"
	"net/http"
)

type WelcomeModel struct {
	DB *sql.DB
}

func (welcome *WelcomeModel) Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	findIpRecord := models.NewIpRecord(ip)
	readIpRecord, err := findIpRecord.Read(welcome.DB)
	if err != nil {
		_, err := findIpRecord.Create(welcome.DB)
		if err != nil {
			return
		}
		generateOutput(w, findIpRecord)
	} else {
		readIpRecord.IncreaseHitCount()
		_, err := readIpRecord.Update(welcome.DB)
		if err != nil {
			return
		}
		generateOutput(w, readIpRecord)
	}
}

func generateOutput(w http.ResponseWriter, entry *models.IpRecord) {
	switch entry.GetHitCount() {
	case uint8(1):
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Welcome, your IP is %s!", entry.GetIp())
	case uint8(2):
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Yeah, I know your IP is %s!", entry.GetIp())
	default:
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "The computer says NO!!!")
	}

}
