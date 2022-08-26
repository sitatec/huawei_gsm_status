package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const searchUrl = "https://www.gsmarena.com/res.php3?sSearch="

var huaweiGSMDeadline = time.Date(2019, time.May, 16, 0, 0, 0, 0, time.UTC)

// 1
var devicesList = []string{
	"Y9 Prime",
	"Y9 2019",
	"P30 Pro",
	"Y7 Pro",
	"nova 7i",
	"P30 Lite",
	"Mate 20 Pro",
	"Honor 8X",
	"nova 3i",
	"Honor 20",
	"nova 3i",
	"Y7 Prime",
	"P30 Lite",
	"P20 Lite",
	"Y6 2019",
	"Mate 10 Lite",
	"P20 Pro",
	"Honor 8A Pro",
	"Y9a",
	"P smart 2021",
	"nova 3",
	"Mate 20",
	"nova 3i",
	"Nova 7 5G Global",
	"Y5 2019",
	"Honor 9X",
	"P30",
	"Mate 10 Pro",
	"P smart",
	"nova 9 SE",
	"Y6p",
	"P10 Lite",
	"P30 Pro",
	"Mate 9",
	"nova 4",
	"Mate 10",
	"nova 7 SE",
	"Honor 20 lite",
	"Mate 30 Pro 5G",
	"Nova 2 Plus Dual SIM TD-LTE",
	"nova 9",
	"P40 Pro",
	"Honor 9X",
	"P Smart S",
	"Honor 8C",
	"Honor 8X",
	"Mate 30",
	"Y7 Prime",
	"P20 Lite",
	"Honor 7X",
	"P10 Lite",
	"Y7p",
	"Honor 9 Lite",
	"Y6 Pro 2019",
	"Honor 6X",
	"Honor 9 Lite",
	"Honor 7C",
	"Honor 10",
	"Honor Play",
	"Honor 10 Lite",
	"Y9 2019",
	"Nova 4e",
	"Honor 7s",
	"Honor 8A Pro",
	"Mate 10 Lite",
	"P10",
	"Y5P",
	"Mate 20 X",
	"Nova 2 Lite",
	"Honor 9i",
	"Nova 4e",
	"Mate 40 Pro",
	"Honor 20",
	"Honor 8X Max",
	"Honor 8X",
	"Mate 20 Pro",
	"Mate 20 X (5G)",
	"Honor 10 Lite",
	"nova 7i",
	"Y7p",
	"Nova 2 Plus",
	"Y7 Pro",
	"nova 3i",
	"Honor 8S",
	"Honor V20",
	"P9 Plus",
	"P10 Plus",
	"P10 Lite",
	"Honor 20 Pro",
	"P8 Lite",
	"Y7p",
	"Honor 7A",
	"Honor 7A Pro",
	"P20 Pro",
	"Honor 10",
	"Y5II",
	"Honor 8",
	"P-Smart+",
	"P50 Pro",
	"Y9 2019",
	"Honor 8X",
	"Nova Plus",
	"P Smart 2019",
	"P Smart 2019",
	"Mate 20 Lite",
	"Mate 20 Lite",
	"P9 Lite",
	"P30 Pro",
	"P30 Pro",
	"Mate 10",
	"nova 3e",
	"nova 8",
	"MATEPAD 2022",
	"P50 Pocket",
	"Honor 8C",
	"Honor V10",
	"Honor 6X",
	"Y6 II",
	"Y5 (2018)",
	"P30",
	"P20",
	"P9 Standard Edition",
	"P smart",
	"Honor 9X Pro",
	"Honor 5X",
	"Y3II",
	"P30 Lite",
	"P30 Lite",
	"Y6 2019",
	"nova 8i",
	"Mate 8",
	"Nova 2",
	"GR3",
	"nova 4",
	"P10 Plus",
	"P30 Pro",
	"Nova Y60",
	"MediaPad T3 10",
	"Mediapad T5",
	"MediaPad T3 10",
	"MatePad T 10s",
	"Y5 2019",
	"P20 Lite",
	"Nova 2 Plus",
	"Mate 10 Pro",
	"Honor 6X",
	"Honor 7X",
	"Y3 (2018)",
	"Nova",
	"MediaPad M5 Pro",
	"Honor Play",
	"Honor 6A Pro",
	"Honor 10X Lite",
	"HONOR 7 PLAY",
	"P30",
	"P40 Pro+",
	"P9 Standard Edition",
	"P smart",
	"Ascend G760-L11",
	"Honor 8A",
	"Honor Magic 4 Pro",
	"Mate 30 Pro",
	"Mate 30 Pro 5G",
	"Y3II",
	"Mate 20 Pro",
	"Mate 40",
	"nova 3",
	"P Smart 2020",
	"Honor 8 Lite",
	"Honor Y6",
	"Honor 9",
	"Honor 9X",
	"Enjoy 5",
	"Enjoy 7 Plus",
	"P9 Lite",
	"P30 Pro",
	"P10 Lite",
	"Honor 7C",
	"Mate 8",
}

var proxiesList = []string{
	"203.30.191.185:80",
}

func main() {
	rand.Seed(time.Now().Unix())
	gsmSupportedDevicesCount := 0
	gsmUnsupportedDevicesCount := 0

	for _, deviceName := range devicesList {
		if isGSMSupported(deviceName) {
			gsmSupportedDevicesCount++
		} else {
			gsmUnsupportedDevicesCount++
		}
		seconds := 10
		if seconds%2 == 0 {
			seconds = 12
		} else if seconds%3 == 0 {
			seconds = 6
		} else {
			seconds = 9
		}

		time.Sleep(time.Second * time.Duration(seconds))
	}
	fmt.Println("supported devices", gsmSupportedDevicesCount)
	fmt.Println("unsupported devices", gsmUnsupportedDevicesCount)
}

func logScrapperRequests(scrapper *colly.Collector) {
	scrapper.OnRequest(func(request *colly.Request) {
		fmt.Println("Loading", request.URL)
	})
	scrapper.OnResponse(func(response *colly.Response) {
		fmt.Println("Processing response from", response.Request.URL)
	})
}

func isGSMSupported(deviceName string) bool {
	fmt.Println("isGSMSupported")
	scrapper, err := getScrapperWithProxy()
	if err != nil {
		fmt.Println(err)
	}
	var gsmSupported bool
	deviceDetailsURL := getDeviceDetailsURL(deviceName)
	fmt.Println("Device url ==", deviceDetailsURL)
	scrapper.OnHTML("#specs-list table tbody tr td[data-spec=status]", func(element *colly.HTMLElement) {
		fmt.Println("OnHTML")
		deviceReleseDateString := strings.Split(element.Text, "Released")[1]
		parts := strings.Split(deviceReleseDateString, ", ")
		if !strings.Contains(parts[1], " ") {
			parts[1] = parts[1] + " 01"
		}
		parts[0], parts[1] = parts[1], parts[0]
		dateString := strings.Join(parts, ", ")
		deviceReleaseDate, err := time.Parse("January 2, 2006", dateString)
		fmt.Println(deviceName, "Released on", deviceReleaseDate)
		if err == nil {
			if deviceReleaseDate.Before(huaweiGSMDeadline) {
				gsmSupported = true
			}
		} else {
			log.Fatalln(err)
		}
	})
	scrapper.Visit("https://www.gsmarena.com/" + deviceDetailsURL)
	return gsmSupported
}

func getDeviceDetailsURL(deviceName string) (deviceDetailsURL string) {
	fmt.Println("getDeviceDetailsURL")
	_scrapper, err := getScrapperWithProxy()
	if err != nil {
		fmt.Println("Error from get proxy", err)
	}
	_scrapper.OnHTML("#review-body li:first-child a:first-child", func(element *colly.HTMLElement) {
		fmt.Println("getDeviceDetailsURL -> OnHTML")
		deviceDetailsURL = element.Attr("href")
	})

	_scrapper.Visit(searchUrl + url.QueryEscape(deviceName))
	return
}

func getScrapperWithProxy() (*colly.Collector, error) {
	scrapper := colly.NewCollector()
	// proxyUrl := proxiesList[rand.Intn(len(proxiesList))]
	// err := setProxy("http://143.110.151.242:3128", scrapper)
	// fmt.Println("Error ==", err)
	scrapper.OnError(func(r *colly.Response, err error) {
		fmt.Println("ON_ORROR")
		fmt.Println("Err ==", err)
		fmt.Println(r.StatusCode)
	})
	return scrapper, nil
}

func setProxy(proxyURL string, scrapper *colly.Collector) error {
	proxyParsed, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}

	scrapper.SetProxyFunc(http.ProxyURL(proxyParsed))

	return nil
}
