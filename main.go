package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	uuid := flag.String("uuid", "", "SKILLS_PASSPORT_UUID")
	secondsWait := flag.Int("sleep", 60, "SECONDS_TO_WAIT_BETWEEN_REQUESTS")
	authorization := flag.String("auth", "", "bearer \"JWT_TOKEN\"")
	flag.Parse()

	if *uuid == "" {
		log.Fatal("uuid flag was not provided (-help)")
	}

	if *authorization == "" {
		log.Fatal("auth flag was not provided (-help)")
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	log.Printf("Start waiting for demoexam results with uuid: %s\n Wait between requests: %d seconds", *uuid, *secondsWait)
	// I think, we will panic only if we could not parse json or get 403 code
	defer func() {
        if r := recover(); r != nil {
            log.Println("Demoexam results might be ready! Exiting program... (After panic :( )")
        }
    }()

	for {
		req, err := http.NewRequest(http.MethodGet, "https://api.dp.worldskills.ru/api/de/skill-passport/" + *uuid, nil)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		
		req.Header.Set("authorization", *authorization)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		// Somehow trying to gather more info
		res := map[string]map[string]interface{}{}
		json.Unmarshal(body, &res)
		e := res["error"]
		code := e["code"]
		if code.(float64) != http.StatusForbidden {
			log.Println("Demoexam results are ready! Exiting program...")
			break
		}
		
		log.Println("Not ready")
		time.Sleep(time.Duration(*secondsWait) * time.Second)
	}
}
