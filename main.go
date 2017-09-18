/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// itfgolang is a web server that feeds back information about the app version
// and kubernetes pod ID to demonstrate distributed microservice deployment
package main

import (
	"expvar"
	"flag"
	"html/template"
	"log"
	"net/http"
	"sync"
	"os/exec"
)

// Command-line flags.
var (
	httpAddr   = flag.String("http", ":8080", "Listen address")
	version    = flag.String("version", "1.0", "itfgolang version")
)

func main() {
	flag.Parse()
	hostname := GetHostname()
	http.Handle("/", NewServer(*version, hostname))
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

// Exported variables for monitoring the server.
// These are exported via HTTP as a JSON object at /debug/vars.
var (
	hitCount       = expvar.NewInt("hitCount")
)

// Server implements the itfgolang server.
// It serves the user interface (it's an http.Handler)
// takes the version var from this code and grabs
// the local container hostname.
type Server struct {
	version string
	hostname string

	mu  sync.RWMutex // protects the read/write
}

// GetHost returns the hostname of the container where the code is running
func GetHostname() string{
	out, err := exec.Command("hostname").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

// NewServer returns an initialized itfgolang server.
func NewServer(version string, hostname string) *Server {
	s := &Server{version: version, hostname: hostname}
	return s
}

// ServeHTTP implements the HTTP user interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hitCount.Add(1)
	s.mu.RLock()
	data := struct {
		Version string
		Hostname string
	}{
		s.version,
		s.hostname,
	}
	s.mu.RUnlock()
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Print(err)
	}
}

// tmpl is the HTML template that drives the user interface.
var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html><html><body><center>
	<h2>Innovation and Tech Forum Golang Demo</h2>
	<h1>
	{{if .Version}}
		ITFGolang Version: {{.Version}} </br>
	{{else}}
		Can't find version :-( </br>
	{{end}}
	{{if .Myhostname}}
		Running on container: {{.Hostname}} </br>
	{{else}}
		Can't find container ID :O/ </br>
	{{end}}
	</h1>
</center></body></html>
`))
