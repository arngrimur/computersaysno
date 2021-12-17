package RESTendpoints

import (
	"datastructures"
	"fmt"
	"net"
	"net/http"
)

var lookUpMap = make(map[string]*datastructures.IpRecord)

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	ipEntry, ipEntryExist := lookUpMap[ip]
	if ipEntryExist {
		ipEntry.IncreaseHitCount()
		generateOutput(w, ipEntry)
	} else {
		lookUpMap[ip] = datastructures.NewIpRecord(ip)
		//save(lookUpMap[ip])
		generateOutput(w,lookUpMap[ip])
	}
}

func generateOutput(w http.ResponseWriter, entry *datastructures.IpRecord) {
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


