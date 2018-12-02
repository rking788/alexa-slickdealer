package main

import (
	"os"

	"github.com/mikeflynn/go-alexa/skillserver"
)

const (
	// Built-in intents
	helpIntent   = "AMAZON.HelpIntent"
	cancelIntent = "AMAZON.CancelIntent"
	stopIntent   = "AMAZON.StopIntent"

	// Custom Intents
	frontpageDealIntent = "FrontpageDealIntent"
	popularDealIntent   = "PopularDealIntent"
	aboutIntent         = "AboutIntent"
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

	var response *skillserver.EchoResponse

	switch request.GetIntentName() {
	case frontpageDealIntent:
		response = handleFrontPageDealIntent()
	case popularDealIntent:
		response = handlePopularDealIntent()
	case helpIntent:
		response = handleHelpIntent()
	case aboutIntent:
		fallthrough
	default:
		response = handleAboutIntent()
	}

	if response == nil {
		response = skillserver.NewEchoResponse()
		response.OutputSpeech("Sorry, something went wrong loading the deals. Please try again later.")
	}

	*echoResponse = *response
}

func handleFrontPageDealIntent() *skillserver.EchoResponse {
	response := skillserver.NewEchoResponse()
	response.OutputSpeech("Frontpage deal data here")
	response.SimpleCard("Frontpage Deals", "Frontpage deal data here")

	return response
}

func handlePopularDealIntent() *skillserver.EchoResponse {
	response := skillserver.NewEchoResponse()
	response.OutputSpeech("Popular deal data here")
	response.SimpleCard("Popular Deals", "Popular deal data here")

	return response
}

func handleAboutIntent() *skillserver.EchoResponse {
	response := skillserver.NewEchoResponse()
	response.OutputSpeech("Slick Dealer was created by Rob in New Hampshire as an unofficial Slick Deals application.")
	response.SimpleCard("About", "Slick Dealer was created by Rob in New Hampshire as an unofficial Slick Deals application.")

	return response
}

func handleHelpIntent() *skillserver.EchoResponse {
	response := skillserver.NewEchoResponse()
	response.OutputSpeech("Help regarding the available commands here.")
	response.SimpleCard("Frontpage Deals", "Help regarding the available commands here.")

	return response
}
