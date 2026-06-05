package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"os"
     "strings"
	"github.com/sbinet/npyio/npy"
	"encoding/csv"
	"strconv"
	"github.com/sbinet/npyio"
	  "gonum.org/v1/gonum/mat"
	  "github.com/knights-analytics/hugot"
	"github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
)


/*-----------------------cleaners---------------------*/
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

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
/*-----------------initialising redis vectror databse----------------------*/
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password docs
		DB:       0,  // use default DB
		Protocol: 2,
	})

	rdb.FTDropIndexWithArgs(ctx,
		"anime_idx",
		&redis.FTDropIndexOptions{
			DeleteDocs: true,
		},
	)

_, err := rdb.FTCreate(ctx,
    "anime_idx",
    &redis.FTCreateOptions{
        OnHash: true,
        Prefix: []any{"anime:"},
    },
    &redis.FieldSchema{
        FieldName: "embedding",
        FieldType: redis.SearchFieldTypeVector,
        VectorArgs: &redis.FTVectorArgs{
            HNSWOptions: &redis.FTHNSWOptions{
                Type:           "FLOAT32",
                Dim:            384,
                DistanceMetric: "COSINE",
            },
        },
    },
).Result()


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
        "userName": os.Getenv("USER"),
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




titles := []string{}

for _, list := range result.Data.MediaListCollection.Lists {
	for _, entry := range list.Entries {

		if entry.Media.Title.English != "" {
			titles = append(titles, entry.Media.Title.English)
		} else {
			titles = append(titles, entry.Media.Title.Romaji)
		}
	}
}







/*---------------------getting anime in databse -----------------------------*/

	in, err := os.Open("/home/saidharahas/jupyter_projects/anime_reco/anime_names_clean.npy")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer in.Close()

	
	file ,err:=os.Open("/home/saidharahas/jupyter_projects/anime_reco/names_to_id.csv")
	defer file.Close()
	reader:=csv.NewReader(file)

	headers,err:=reader.Read()
	fmt.Println(headers[0])
	nameid:=make(map[string]int)

	for{
		record,err:=reader.Read()
		if err!=nil{
break
		}
		v:=cleanNpyString(record[0])
    id, err := strconv.Atoi(record[1])
		nameid[normalize(v)]=id

		fmt.Println(normalize(v))
fmt.Println(record[1])
	}


    var data []string
	err = npy.Read(in, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
dataMap := make(map[string]bool)
		for _, v := range data {
	    v1:=cleanNpyString(v)		
		dataMap[normalize(v1)]=true
		
	}
	
indata:=[]int{}

a:=0

for _,i:= range titles{
	_,ok:=dataMap[normalize(i)]
	if ok {
		x,_:=nameid[normalize(i)]
     indata=append(indata,x)
	 a++
	}
}
fmt.Println(a)
calc()



}


func calc(){

	in, err := os.Open("/home/saidharahas/jupyter_projects/anime_reco/W.npy")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer in.Close()

	 var m mat.Dense

    
    if err := npyio.Read(in, &m); err != nil {
        log.Fatal(err)
    }
     rows, cols := m.Dims()
    fmt.Printf("Loaded matrix: %d x %d\n", rows, cols)
    fmt.Printf("Element at (0,0): %f\n", m.At(0, 0))


}