package main

import (
	"bytes"
	"encoding/json"
	
	"net/http"
	"log"
	"os"
     "strings"

	"fmt"
	

	"github.com/joho/godotenv"
	"math"
	"encoding/binary"
	
)
func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	
	replacer := strings.NewReplacer(
		":", "",
		"·", "",
		"★", "",
		"!", "",
		"'", "",
		"’", "",
		"-", "",
		"  ", " ",
	)

	s = replacer.Replace(s)

	return s
}
func cleanNpyString(s string) string {
    return strings.ReplaceAll(s, "\x00", "")
}
/*-----------redis vector db helper fiction---------------*/
func floatsToBytes(fs []float32) []byte {
	buf := make([]byte, len(fs)*4)

	for i, f := range fs {
		u := math.Float32bits(f)
		binary.NativeEndian.PutUint32(buf[i*4:], u)
	}

	return buf
}
func Float32SliceToBytes(vec []float32) []byte {
    buf := new(bytes.Buffer)

    // write as raw binary (little endian is standard for Redis vectors)
    _ = binary.Write(buf, binary.LittleEndian, vec)

    return buf.Bytes()
}


/*------------------------api strcutres----------------------------*/

type Response struct {
	Data struct {
		MediaListCollection struct {
			Lists []struct {
				Name    string `json:"name"`
				Entries []struct {
					Media struct {
						ID    int `json:"id"`
						Title struct {
							Romaji  string `json:"romaji"`
							English string `json:"english"`
						} `json:"title"`
					} `json:"media"`
					Score  float64 `json:"score"`
					Status string  `json:"status"`
				} `json:"entries"`
			} `json:"lists"`
		} `json:"MediaListCollection"`
	} `json:"data"`
}


func main(){
	titles := []string{}

		err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

query := `
query ($userName: String) {
  MediaListCollection(userName: $userName, type: ANIME, status: COMPLETED) {
    lists {
      name
      entries {
        media {
          id
          title {
            romaji
            english
          }
        }
        score
        status
      }
    }
  }
}`
payload := map[string]interface{}{
    "query": query,
    "variables": map[string]string{
        "userName": os.Getenv("ANILIST_USER"),
    },
}
	postBody, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		"https://graphql.anilist.co",
		"application/json",
		bytes.NewBuffer(postBody),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result Response

/*---------------------decoodug result in to format of go object from json ----------------------------*/
err = json.NewDecoder(resp.Body).Decode(&result)
if err != nil {
	log.Fatal(err)
}






for _, list := range result.Data.MediaListCollection.Lists {
	for _, entry := range list.Entries {

		if entry.Media.Title.English != "" {
			titles = append(titles, entry.Media.Title.English)
		} else {
			titles = append(titles, entry.Media.Title.Romaji)
		}
	}
}
file, err := os.Create("output.txt")
if err != nil {
	log.Fatal(err)
}
defer file.Close()
fmt.Println(len(titles))
for _, title := range titles {
	_, err := file.WriteString(title + "\n")
	fmt.Println(title)
	if err != nil {
		log.Fatal(err)
	}
}

log.Println("TXT file created successfully")
}