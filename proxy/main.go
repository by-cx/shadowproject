package main

type TaskHost struct {
	UUID               string
	Image              string   // Docker image this task uses
	RegisteredBackends []string // List of URLs for backend servers
}

//var TaskHosts []TaskHost

func main() {
	//rpURL, err := url.ParseRequestURI("http://rosti.cz")
	//if err != nil {
	//	log.Fatalln(err)
	//}
}
