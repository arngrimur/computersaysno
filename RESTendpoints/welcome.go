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

//var lookUpMap = make(map[string]*models.IpRecord)

func (welcome *WelcomeModel) Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	//ipEntry, ipEntryExist := lookUpMap[ip]
	findIpRecord := models.NewIpRecord(ip)
	readIpRecord, err := findIpRecord.Read(welcome.DB)
	if err != nil {
		findIpRecord.Create(welcome.DB)
		generateOutput(w, findIpRecord)
	} else {
		readIpRecord.IncreaseHitCount()
		readIpRecord.Update(welcome.DB)
		generateOutput(w, readIpRecord)
	}

	//if ipEntryExist {
	//	ipEntry.IncreaseHitCount()
	//	generateOutput(w, ipEntry)
	//} else {
	//	lookUpMap[ip] = models.NewIpRecord(ip)
	//	//save(lookUpMap[ip])
	//	generateOutput(w, lookUpMap[ip])
	//}
}

func generateOutput(w http.ResponseWriter, entry *models.IpRecord) {
	switch entry.GetHitCount() {
	case uint8(1):
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Welcome, your IP is %s!", entry.GetIp())
		break
	case uint8(2):
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Yeah, I know your IP is %s!", entry.GetIp())
		break
	default:
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "The computer says NO!!!")
		break
	}
}
