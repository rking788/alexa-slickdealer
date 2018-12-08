package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

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
		"/health": skillserver.StdApplication{
			Methods: "GET",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Ok!"))
			},
		},
	}
)

type feedResponse struct {
	Channel struct {
		Item []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}

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

	feedResponse, _ := requestFeed("frontpage")
	builder := skillserver.NewSSMLTextBuilder()
	cardBody := strings.Builder{}

	builder.AppendSentence("Here are the current frontpage deals:")
	cardBody.WriteString("Here are the current frontpage deals:")
	for _, item := range feedResponse.Channel.Item[:3] {
		builder.AppendSentence(item.Title)
		cardBody.WriteString(item.Title)
	}

	response := skillserver.NewEchoResponse()
	response.OutputSpeechSSML(builder.Build())
	response.SimpleCard("Frontpage Deals", cardBody.String())
	return response
}

func handlePopularDealIntent() *skillserver.EchoResponse {

	feedResponse, _ := requestFeed("popdeals")
	builder := skillserver.NewSSMLTextBuilder()
	cardBody := strings.Builder{}

	builder.AppendSentence("Here are the current popular deals:")
	cardBody.WriteString("Here are the current popular deals:")
	for _, item := range feedResponse.Channel.Item[:3] {
		builder.AppendSentence(item.Title)
		cardBody.WriteString(item.Title)
	}

	response := skillserver.NewEchoResponse()
	response.OutputSpeechSSML(builder.Build())
	response.SimpleCard("Popular Deals", cardBody.String())
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
	builder := skillserver.NewSSMLTextBuilder()

	builder.AppendSentence("Here are some things you can ask: ")
	builder.AppendSentence("Give me the frontpage deals.")
	builder.AppendSentence("Give me the popular deals.")

	return response.OutputSpeechSSML(builder.Build())
}

func requestFeed(mode string) (*feedResponse, error) {

	endpoint, _ := url.Parse("https://slickdeals.net/newsearch.php")
	queryParams := endpoint.Query()
	queryParams.Set("mode", mode)
	queryParams.Set("searcharea", "deals")
	queryParams.Set("searchin", "first")
	queryParams.Set("rss", "1")

	endpoint.RawQuery = queryParams.Encode()
	response, err := http.Get(endpoint.String())
	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(response.Body)
	feed := &feedResponse{}
	xml.Unmarshal(data, &feed)

	return feed, nil
}
