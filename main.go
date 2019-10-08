package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Profile struct {
	IsCustomImage bool   `json:"is_custom_image"`
	Image512      string `json:"image_512"`
}

type User struct {
	Name    string  `json:"name"`
	Deleted bool    `json:"deleted"`
	Profile Profile `json:"profile"`
}

type UserList struct {
	OK      bool   `json:"ok"`
	Members []User `json:"members"`
}

func main() {

	token := flag.String("t", "", "token")
	flag.Parse()

	resp, err := http.Get(
		fmt.Sprintf("https://slack.com/api/users.list?token=%s", *token))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var ul UserList
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&ul); err != nil {
		log.Fatal(err)
	}

	var names []string
	if !ul.OK {
		log.Fatal("Not OK")
	}
	for _, u := range ul.Members {
		if u.Deleted {
			continue
		}
		if u.Profile.IsCustomImage {
			continue
		}

		r, err := url.Parse(u.Profile.Image512)
		if err != nil {
			log.Fatal(err)
		}
		v, err := url.ParseQuery(r.RawQuery)
		if err != nil {
			log.Fatal(err)
		}
		v["d"] = []string{"404"}
		r.RawQuery = v.Encode()
		resp, err := http.Get(r.String())
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode == 404 {
			names = append(names, "@"+u.Name)
		}
	}

	fmt.Printf("People without photo in their profiles: %v\n",
		strings.Join(names, ", "))
	for _, n := range names {
		fmt.Printf("/remind %s tomorrow set a profile photo, please\n", n)
	}
}
