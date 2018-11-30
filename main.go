package main

import (
	"fmt"
	"os"

	"github.com/mikeflynn/go-alexa/skillserver"
)

var (
	slickDealsAppID = os.Getenv("SLICK_DEALS_APP_ID")
	applications    = map[string]interface{}{
		"/echo/slickdeals": skillserver.EchoApplication{
			AppID:          slickDealsAppID,
			OnIntent:       intentHandler,
			OnLaunch:       launchHandler,
			OnSessionEnded: sessionEndedHandler,
		},
	}
)

func main() {

	port := os.Getenv("PORT")

	skillserver.Run(applications, port)
}

func launchHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
	echoResponse.OutputSpeech("You have successfully launched a new session.")
	echoResponse.EndSession(false)
}

func sessionEndedHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
	echoResponse.OutputSpeech("Session ended.")
}

func intentHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
	echoResponse.OutputSpeech(fmt.Sprintf("You have invoked the %s intent.", request.GetIntentName()))
}
