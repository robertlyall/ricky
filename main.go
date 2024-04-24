package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

type Character struct {
	Gender  string `json:"gender"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
	Status  string `json:"status"`
	Type    string `json:"type"`
}

type Episode struct {
	Characters []string `json:"characters"`
	ID 		int      `json:"id"`
	Name 	string   `json:"name"`
}

type Filter struct {
	Key   string
	Value string
}

type Location struct {
	Dimension string `json:"dimension"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

func getCharacters(filter Filter) []Character {
	var characters []Character

	for page := 1; ; page++ {
		url := ""

		if filter.Key == "" {
			url = fmt.Sprintf("https://rickandmortyapi.com/api/character/?page=%d", page)
		} else {
			url = fmt.Sprintf("https://rickandmortyapi.com/api/character/?page=%d&%s=%s", page, filter.Key, filter.Value)
		}

		resp, err := http.Get(url)

		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		defer resp.Body.Close()

		var data struct {
			Results []Character `json:"results"`
		}

		err = json.NewDecoder(resp.Body).Decode(&data)

		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		characters = append(characters, data.Results...)

		if len(data.Results) == 0 {
			break
		}
	}

	return characters
}

func getEpisodes() []Episode {
	var episodes []Episode

	for page := 1; ; page++ {
		resp, err := http.Get(fmt.Sprintf("https://rickandmortyapi.com/api/episode/?page=%d", page))

		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		defer resp.Body.Close()

		var data struct {
			Results []Episode `json:"results"`
		}

		err = json.NewDecoder(resp.Body).Decode(&data)

		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		episodes = append(episodes, data.Results...)

		if len(data.Results) == 0 {
			break
		}
	}

	return episodes
}

func getLocations(filter Filter) []Location {
	var locations []Location

	for page := 1; ; page++ {
		url := ""

		if filter.Key == "" {
			url = fmt.Sprintf("https://rickandmortyapi.com/api/location/?page=%d", page)
		} else {
			url = fmt.Sprintf("https://rickandmortyapi.com/api/location/?page=%d&%s=%s", page, filter.Key, filter.Value)
		}

		resp, err := http.Get(url)

		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		defer resp.Body.Close()

		var data struct {
			Results []Location `json:"results"`
		}

		err = json.NewDecoder(resp.Body).Decode(&data)

		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		locations = append(locations, data.Results...)

		if len(data.Results) == 0 {
			break
		}
	}

	return locations
}

func main() {
	var count bool
	var filter string

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Action: func(cCtx *cli.Context) error {
					var characters []Character

					if filter == "" {
						fmt.Println("Please provide a filter")
						return nil
					} else {
						parts := strings.Split(filter, "=")

						key := parts[0]
						value := parts[1]

						characters = getCharacters(Filter{Key: key, Value: value})
					}

					episodes := getEpisodes()

					fmt.Println(len(episodes))

					count := 0

					for _, episode := range episodes {
						found := false

						for _, url := range episode.Characters {
							for _, character := range characters {
								if strings.Contains(url, fmt.Sprintf("/character/%d", character.ID)) {
									found = true
									break
								}
							}

							if found {
								break;
							}
						}

						if found {
							count++
						}
					}

					fmt.Println("Count:", count)

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Destination: &filter,
						Name:        "filter",
						Usage:       "Filters the API response, for example --filter species=Robot",
					},
				},
				Name:  "appearances",
				Usage: "The number of unique episodes that all characters matching a specific --filter have appeared in. For example --filter species=Robot would display the total number of unique episodes where a character of the Robot species appears.",
			},
			{
				Action: func(cCtx *cli.Context) error {
					var characters []Character

					if filter == "" {
						characters = getCharacters(Filter{})
					} else {
						parts := strings.Split(filter, "=")

						key := parts[0]
						value := parts[1]

						characters = getCharacters(Filter{Key: key, Value: value})
					}

					if count {
						fmt.Println("Count:", len(characters))
						return nil
					}

					tbl := table.New("ID", "Name", "Species", "Status", "Gender", "Type")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

					for _, character := range characters {
						tbl.AddRow(character.ID, character.Name, character.Species, character.Status, character.Gender, character.Type)
					}

					tbl.Print()

					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Destination: &count,
						Name:        "count",
						Usage:       "Just displays the total number matching characters",
					},
					&cli.StringFlag{
						Destination: &filter,
						Name:        "filter",
						Usage:       "Filters the API response, for example --filter species=Robot",
					},
				},
				Name:  "characters",
				Usage: "Displays characters and displays basic details",
			},
			{
				Action: func(cCtx *cli.Context) error {
					var locations []Location

					if filter == "" {
						locations = getLocations(Filter{})
					} else {
						parts := strings.Split(filter, "=")

						key := parts[0]
						value := parts[1]

						locations = getLocations(Filter{Key: key, Value: value})
					}

					if count {
						fmt.Println("Count:", len(locations))
						return nil
					}

					tbl := table.New("ID", "Name")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

					for _, location := range locations {
						tbl.AddRow(location.ID, location.Name)
					}

					tbl.Print()

					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Destination: &count,
						Name:        "count",
						Usage:       "Just displays the total number matching locations",
					},
					&cli.StringFlag{
						Destination: &filter,
						Name:        "filter",
						Usage:       "Filters the API response, for example --filter type=Planet",
					},
				},
				Name:  "locations",
				Usage: "Displays locations and displays basic details",
			},
		},
		Name:  "ricky",
		Usage: "Displays characters and displays basic details",
		Action: func(c *cli.Context) error {
			fmt.Println("Use --help to see available commands")
			return nil
		},
	}

	app.Run(os.Args)
}
