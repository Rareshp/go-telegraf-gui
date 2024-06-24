package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Config structure to hold form data
type Config struct {
	Name              string
	Endpoint          string
	NumNodes          int
	Nodes             []Node
}

// Node structure to hold node data
type Node struct {
	Name           string
	Namespace      string
	IdentifierType string
	Identifier     string
}

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	numNodes, err := strconv.Atoi(r.FormValue("num_nodes"))
	if err != nil || numNodes <= 0 {
		http.Error(w, "Invalid number of nodes", http.StatusBadRequest)
		return
	}

	config := Config{
		Name:     r.FormValue("name"),
		Endpoint: r.FormValue("endpoint"),
		NumNodes: numNodes,
		Nodes:    make([]Node, numNodes),
	}

	if err := tmpl.ExecuteTemplate(w, "nodes.html", config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	numNodes, err := strconv.Atoi(r.FormValue("num_nodes"))
	if err != nil || numNodes <= 0 {
		http.Error(w, "Invalid number of nodes", http.StatusBadRequest)
		return
	}

	config := Config{
		Name:     r.FormValue("name"),
		Endpoint: r.FormValue("endpoint"),
		NumNodes: numNodes,
	}

	for i := 0; i < numNodes; i++ {
		config.Nodes = append(config.Nodes, Node{
			Name:           r.FormValue("name" + strconv.Itoa(i)),
			Namespace:      r.FormValue("namespace" + strconv.Itoa(i)),
			IdentifierType: r.FormValue("identifier_type" + strconv.Itoa(i)),
			Identifier:     r.FormValue("identifier" + strconv.Itoa(i)),
		})
	}

	configContent := generateConfigContent(config)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(configContent))
}

func generateConfigContent(config Config) string {
	var sb strings.Builder
	sb.WriteString(`## Simply copy paste the code below in your telegraf.conf`)
	sb.WriteString(`[[inputs.opcua]]` + "\n")
	sb.WriteString(`  name = "` + config.Name + `"` + "\n\n")
	sb.WriteString(`  endpoint = "` + config.Endpoint + `"` + "\n\n")
	sb.WriteString(`  ## Maximum time allowed to establish a connect to the endpoint.` + "\n")
	sb.WriteString(`  connect_timeout = "10s"` + "\n\n")
	sb.WriteString(`  ## Maximum time allowed for a request over the estabilished connection.` + "\n")
	sb.WriteString(`  request_timeout = "5s"` + "\n\n")
	sb.WriteString(`  ## Security policy, one of "None", "Basic128Rsa15", "Basic256", "Basic256Sha256", or "auto".` + "\n")
	sb.WriteString(`  security_policy = "None"` + "\n\n")
	sb.WriteString(`  ## Security mode, one of "None", "Sign", "SignAndEncrypt", or "auto".` + "\n")
	sb.WriteString(`  security_mode = "None"` + "\n\n")
	sb.WriteString(`  ## Path to cert.pem. Required when security mode or policy isn't "None".
	## If cert path is not supplied, self-signed cert and key will be generated.
	# certificate = "/etc/telegraf/cert.pem"
	## Path to private key.pem. Required when security mode or policy isn't "None".
	## If key path is not supplied, self-signed cert and key will be generated.
	# private_key = "/etc/telegraf/key.pem"` + "\n\n")
	sb.WriteString(`  ## Authentication Method, one of "Certificate", "UserName", or "Anonymous".` + "\n")
	sb.WriteString(`  auth_method = "Anonymous"` + "\n\n")
	sb.WriteString(`  ## Username. Required for auth_method = "UserName"
	# username = ""
	
	## Password. Required for auth_method = "UserName"
	# password = ""` + "\n\n")
	sb.WriteString(`  nodes = [` + "\n")
	for _, node := range config.Nodes {
		sb.WriteString(`    {name="` + node.Name + `", namespace="` + node.Namespace + `", identifier_type="` + node.IdentifierType + `", identifier="` + node.Identifier + `"},` + "\n")
	}
	sb.WriteString(`  ]` + "\n")
	return sb.String()
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/nodes", nodesHandler)
	http.HandleFunc("/generate", generateHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
