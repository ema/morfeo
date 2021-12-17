package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	secret := os.Getenv("MORFEO_SECRET")
	if secret == "" {
		log.Fatal("MORFEO_SECRET must be set")
	}

	addr := os.Getenv("MORFEO_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
        <html>
            <head><title>Morfeo</title></head>
            <body>
                <h1>Morfeo</h1>
                <form method="post" action="/suspend">
                    <input type="text" name="secret" />
                    <input type="submit" />
                </form>
            </body>
        </html>
        `)
	})

	http.HandleFunc("/suspend", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Only POST requests are allowed")
			return
		}

		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		if r.FormValue("secret") != os.Getenv("MORFEO_SECRET") {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Bad secret")
		} else {
			fmt.Fprintf(w, "OK boss, going to sleep")

			// Suspend the system
			cmd := exec.Command("systemctl", "suspend")
			_, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			} else {
				log.Printf("suspending")
			}
		}
	})

	fmt.Printf("Starting server on %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
