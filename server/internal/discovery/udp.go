package discovery

import (
	"encoding/json"
	"log"
	"net"
)

func Start(log *log.Logger, port int, httpPort int) {

	go func() {

		addr := net.UDPAddr{Port: port, IP: net.IPv4zero}
		conn, _ := net.ListenUDP("udp4", &addr)

		buf := make([]byte, 2048)

		for {

			n, remote, _ := conn.ReadFromUDP(buf)

			if string(buf[:n]) == "DISCOVER" {

				resp := map[string]interface{}{
					"magic":     "MY_DEVICE",
					"http_port": httpPort,
				}

				data, _ := json.Marshal(resp)

				conn.WriteToUDP(data, remote)

				log.Println("discovery reply", remote)
			}
		}
	}()
}
