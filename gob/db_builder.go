package main

import (
	"bytes"
	
	"fmt"
	
	"log"
	"os"
     "strings"
	
	"github.com/sbinet/npyio"
	  "gonum.org/v1/gonum/mat"
	"github.com/redis/go-redis/v9"

	"math"
	"encoding/binary"
	"context"
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
func Float32SliceToBytes(vec []float32) []byte {
    buf := new(bytes.Buffer)

    // write as raw binary (little endian is standard for Redis vectors)
    _ = binary.Write(buf, binary.LittleEndian, vec)

    return buf.Bytes()
}


/*------------------------api strcutres----------------------------*/

func main() {
	 godotenv.Load("../.env")
/*-----------------initialising redis vectror databse----------------------*/
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Username: "default",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

_, err:=rdb.FTCreate(ctx,
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


/*---------------loading the aime vectors---------------*/
	win, err := os.Open(os.Getenv("HPATH"))
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



/*------------------inserting into the hnsw storage----------------------*/
pipe := rdb.Pipeline()
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
	pipe.HSet(ctx,
    key,
    "embedding", Float32SliceToBytes(vec),
)
}




_, err= pipe.Exec(ctx)




}


