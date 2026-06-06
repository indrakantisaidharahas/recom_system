package main

import (
	"bytes"
	
	"fmt"
	
	"log"
	"os"
     "strings"
	"github.com/sbinet/npyio/npy"
	"encoding/csv"
	"strconv"
	"github.com/sbinet/npyio"
	  "gonum.org/v1/gonum/mat"
	"github.com/redis/go-redis/v9"
	"bufio"
	"math"
	"encoding/binary"
	"context"
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

func main() {
	// err := godotenv.Load()
    // if err != nil {
    //     log.Fatal("Error loading .env file")
    // }
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

_, err:= rdb.FTCreate(ctx,
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
                Dim:            50,
                DistanceMetric: "COSINE",
            },
        },
    },
).Result()

filet, err := os.Open("output.txt")
	if err != nil {
		panic(err)
	}
	defer filet.Close()

	var titles []string

	scanner := bufio.NewScanner(filet)
	for scanner.Scan() {
		titles= append(titles, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println("total lines:", len(titles))
	






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
	idname:=make(map[int]string)

	for{
		record,err:=reader.Read()
		if err!=nil{
break
		}
		v:=cleanNpyString(record[0])
    id, err := strconv.Atoi(record[1])
		nameid[normalize(v)]=id
		  idname[id]=normalize(v)

		fmt.Println(normalize(v))
//fmt.Println(record[1])
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

	win, err := os.Open("/home/saidharahas/jupyter_projects/anime_reco/H.npy")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer win.Close()

	 var m mat.Dense

    
    if err := npyio.Read(win, &m); err != nil {
        log.Fatal(err)
    }
     rows, cols := m.Dims()
    fmt.Printf("Loaded matrix: %d x %d\n", rows, cols)
    fmt.Printf("Element at (0,0): %f\n", m.At(0, 0))


for i:=0;i<cols;i++{
	colView:=m.ColView(i)
	n := colView.Len()

// 2. Allocate float32 slice
vec := make([]float32, n)

// 3. Iterate and cast
key := fmt.Sprintf("anime:%d", i)
for j := 0; j < rows; j++ {
    vec[j] = float32(colView.AtVec(j))
}
	rdb.HSet(ctx,
    key,
    "embedding", Float32SliceToBytes(vec),
)
}

// for i:=range titles{
// 	id:=nameid[normalize(i)]
	
// }
// 1. get query anime embedding
vecBytes, _ := rdb.HGet(ctx, "anime:5114", "embedding").Bytes()

// 2. search
results, err := rdb.FTSearchWithArgs(ctx,
    "anime_idx",
    "*=>[KNN 5 @embedding $vec AS score]",
    &redis.FTSearchOptions{
        Return: []redis.FTSearchReturn{
            {FieldName: "score"},
        },
        DialectVersion: 2,
        Params: map[string]any{
            "vec": vecBytes,
        },
    },
).Result()
if err!=nil{
	fmt.Println("somethings wrng")
}
for _, doc := range results.Docs {
  parts := strings.Split(doc.ID, ":")
    if len(parts) != 2 {
		fmt.Println("wrng")
        continue
    }

    numID, err := strconv.Atoi(parts[1])
    if err != nil {
        continue
    }

    fmt.Println("Numeric ID:", numID)
    fmt.Println(idname[numID])
    fmt.Println("Score:", doc.Fields["score"])
}


}


// func calc(){


// }