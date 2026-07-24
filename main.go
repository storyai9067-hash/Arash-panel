package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

type Server struct {
	ID       int
	IP       string
	Port     string
	Username string
	Password string
	Donor    string
}

var donatedServers []Server
var mu sync.Mutex
var idCounter = 1

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := `<!DOCTYPE html>
		<html>
		<head><title>Arash Panel</title></head>
		<body>
			<h1>🚀 Arash Panel</h1>
			<h2>اهدای سرور</h2>
			<form action="/donate" method="POST">
				<input name="ip" placeholder="آی‌پی سرور" required><br>
				<input name="port" placeholder="پورت" required><br>
				<input name="username" placeholder="نام کاربری" required><br>
				<input name="password" placeholder="رمز عبور" required><br>
				<input name="donor" placeholder="نام اهداکننده" required><br>
				<button type="submit">اهدای سرور</button>
			</form>
			<h2>سرورهای اهدا شده</h2>
			<ul>
				{{range .}}
				<li>{{.IP}}:{{.Port}} - {{.Username}} (اهداکننده: {{.Donor}})</li>
				{{else}}
				<li>هیچ سروری اهدا نشده است.</li>
				{{end}}
			</ul>
		</body>
		</html>`
		t := template.Must(template.New("index").Parse(tmpl))
		mu.Lock()
		t.Execute(w, donatedServers)
		mu.Unlock()
	})

	http.HandleFunc("/get-config", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if len(donatedServers) == 0 {
			http.Error(w, "No server donated yet!", http.StatusNotFound)
			return
		}
		s := donatedServers[0]
		config := fmt.Sprintf("vless://%s:%s@%s:%s?encryption=none#ArashVPN", s.Username, s.Password, s.IP, s.Port)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"config": config})
	})

	http.HandleFunc("/donate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		mu.Lock()
		defer mu.Unlock()
		if len(donatedServers) >= 10 {
			http.Error(w, "حداکثر ۱۰ سرور قابل اهدا است.", http.StatusBadRequest)
			return
		}
		donatedServers = append(donatedServers, Server{
			ID:       idCounter,
			IP:       r.FormValue("ip"),
			Port:     r.FormValue("port"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			Donor:    r.FormValue("donor"),
		})
		idCounter++
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fmt.Println("🚀 Arash Panel is running on :8080")
	http.ListenAndServe(":8080", nil)
}
