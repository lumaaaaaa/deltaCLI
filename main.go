// devin's midflight speedcode challenge
// WiFiOnboard/Delta CLI flight tracker
//

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	isAirborne bool
)

func status() {
	// TODO: add proper animation for fun :)
	animationCharacter := "🌎"

	resp, err := http.Get("https://wifi.inflightinternet.com/abp/v2/statusTray?fig2=true")
	if err != nil {
		fmt.Println("⚠️failed to contact the flight status server!")
		os.Exit(3)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	isAirborne = strings.Split(strings.Split(string(body), "\"flight_state\":\"")[1], "\"")[0] == "IN_AIR"
	airline := strings.Split(strings.Split(string(body), "\"name\":\"")[1], "\"")[0]
	flightNumber := strings.Split(strings.Split(string(body), "\"flightNumberInfo\":\"")[1], "\"")[0]
	manufacturer := strings.Split(strings.Split(string(body), "\"manufacturer\":\"")[1], "\"")[0]
	model := strings.Split(strings.Split(string(body), "\"make\":\"")[1], "\"")[0]
	departureAirportCode := strings.Split(strings.Split(string(body), "\"departureAirportCodeIata\":\"")[1], "\"")[0]
	destinationAirportCode := strings.Split(strings.Split(string(body), "\"destinationAirportCodeIata\":\"")[1], "\"")[0]
	gate := strings.Split(strings.Split(strings.Split(string(body), "\"arrival\":{")[1], "\"gate\":\"")[1], "\"")[0]

	fmt.Printf("✈️ %s %s | %s %s | %s -> %s | Gate %s\n", airline, flightNumber, manufacturer, model, departureAirportCode, destinationAirportCode, gate)

	for isAirborne {
		resp, err = http.Get("https://wifi.inflightinternet.com/abp/v2/statusTray?fig2=true")
		if err != nil {
			fmt.Println("⚠️failed to contact the flight status server!")
			os.Exit(3)
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		isAirborne = strings.Split(strings.Split(string(body), "\"flight_state\":\"")[1], "\"")[0] == "IN_AIR"
		latitude := strings.Split(strings.Split(string(body), "\"latitude\":")[1], ",")[0]
		longitude := strings.Split(strings.Split(string(body), "\"longitude\":")[1], ",")[0]
		altitude := strings.Split(strings.Split(string(body), "\"altitude\":")[1], ",")[0]
		airSpeed := strings.Split(strings.Split(string(body), "\"horizontalVelocity\":\"")[1], "\"")[0]
		remainingTime := strings.Split(strings.Split(string(body), "\"timeToLand\":\"")[1], "\"")[0]
		minutesToLand, err := strconv.Atoi(remainingTime)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("\r%s ETA: %s | Position: %s, %s | Altitude: %sft | Air Speed: %smph      ", animationCharacter, time.Until(time.Now().Add(time.Duration(minutesToLand*60000000000))), latitude, longitude, altitude, airSpeed)

		time.Sleep(250 * time.Millisecond)

		// redo this please this is so awful
		// my laptop is dying and I don't have the charger, so it's gotta stay like this for now :(
		if animationCharacter == "🌎" {
			animationCharacter = "🌍"
		} else if animationCharacter == "🌍" {
			animationCharacter = "🌏"
		} else {
			animationCharacter = "🌎"
		}
	}

	fmt.Println("🎯 your flight is no longer airborne.")
	os.Exit(0)
}

func main() {
	fmt.Println("🎯 deltaCLI - flight status")

	go status()

	select {}
}
