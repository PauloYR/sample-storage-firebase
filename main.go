package main

const (
	bucket = "palco-planner-test.appspot.com"
)

func main() {
	AppConfig := NewAppConfig()

	AppConfig.server.Start()
}
