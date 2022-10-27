package main

import (
	"flag"
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	markdown "github.com/MichaelMure/go-term-markdown"
	"io/ioutil"
	"net/http"
	"os"
)

type arguments struct {
	weatherCity string
}

func (a *arguments) parseArgs(args []string) error {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.StringVar(&a.weatherCity, "weather", "", "check weather")
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, `flags: %s`, os.Args[0])
		f.PrintDefaults()
		os.Exit(1)
	}
	if err := f.Parse(args[1:]); err != nil {
		return err
	}
	return nil
}

func Execute() error {
	args := &arguments{}
	if err := args.parseArgs(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// weather
	if args.weatherCity != "" {
		weather, err := GetWeather(args.weatherCity)
		if err != nil {
			return err
		}
		fmt.Println(weather)
	}
	return nil
}

func getWeatherData(city string) (string, error) {
	url := "https://wttr.in/" + city
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	weather := string(all)
	return weather, nil
}

func GetWeather(city string) (string, error) {
	weatherData, err := getWeatherData(city)
	if err != nil {
		return "", err
	}
	s := getMD(weatherData)
	result := markdown.Render(s, 280, 6)
	return string(result), nil
}

func getMD(html string) string {
	converter := md.NewConverter("", true, nil)
	res, err := converter.ConvertString(html)
	if err != nil {
		return ""
	}
	return res
}

func main() {
	err := Execute()
	if err != nil {
		fmt.Println(err)
	}
}
