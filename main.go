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
		<head>
			<title>Arash Panel</title>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body { font-family: Arial, sans-serif; background: #0f172a; color: #fff; padding: 20px; }
				.container { max-width: 800px; margin: 0 auto; background: #1e293b; padding: 20px; border-radius: 15px; }
				h1 { color: #38bdf8; text-align: center; }
				.social-buttons { display: flex; justify-content: center; gap: 20px; margin: 20px 0; }
				.social-btn { display: inline-block; padding: 10px 20px; background: #334155; color: #fff; border-radius: 30px; text-decoration: none; font-weight: bold; transition: 0.3s; }
				.social-btn:hover { transform: scale(1.05); background: #475569; }
				.social-btn.github { background: #24292e; }
				.social-btn.telegram { background: #0088cc; }
				form { background: #0f172a; padding: 20px; border-radius: 10px; margin: 20px 0; }
				input { width: 100%; padding: 10px; margin: 8px 0; border: 1px solid #334155; border-radius: 8px; background: #1e293b; color: #fff; }
				button { background: #38bdf8; color: #0f172a; padding: 12px 20px; border: none; border-radius: 8px; font-weight: bold; cursor: pointer; }
				button:hover { background: #0ea5e9; }
				ul { list-style: none; padding: 0; }
				li { background: #0f172a; margin: 5px 0; padding: 10px; border-radius: 8px; border-left: 4px solid #38bdf8; }
				.footer { text-align: center; margin-top: 30px; color: #94a3b8; font-size: 14px; }
			</style>
		</head>
		<body>
		<div class="container">
			<h1>🚀 Arash Panel</h1>

			<!-- دکمه‌های اجتماعی -->
			<div class="social-buttons">
				<a href="https://github.com/storyai9067-hash/Arash-panel" target="_blank" class="social-btn github">
					⭐ حمایت در گیت‌هاب
				</a>
				<a href="https://t.me/panelArashfree" target="_blank" class="social-btn telegram">
					📱 کانال تلگرام
				</a>
			</div>

			<h2>🎁 اهدای سرور</h2>
			<form action="/donate" method="POST">
				<input name="ip" placeholder="آی‌پی سرور" required>
				<input name="port" placeholder="پورت" required>
				<input name="username" placeholder="نام کاربری" required>
				<input name="password" placeholder="رمز عبور" required>
				<input name="donor" placeholder="نام اهداکننده" required>
				<button type="submit">اهدای سرور</button>
			</form>

			<h2>📋 سرورهای اهدا شده</h2>
			<ul>
				{{range .}}
				<li>{{.IP}}:{{.Port}} - {{.Username}} (اهداکننده: {{.Donor}})</li>
				{{else}}
				<li>هیچ سروری اهدا نشده است.</li>
				{{end}}
			</ul>

			<div class="footer">
				<a href="/get-config" style="color: #38bdf8;">🔗 دریافت کانفیگ</a>
				| ساخته شده با ❤️ توسط Arash
			</div>
		</div>
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
