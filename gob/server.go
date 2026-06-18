package main

import (
    
    "fmt"
   
    "net/http"
    "os"
    "strings"
   
    "encoding/csv"
    "strconv"
    
    "sort"  
    "github.com/redis/go-redis/v9"
   
   
   
    "encoding/json"
    "context"

    "bytes"
    "log"
 "github.com/joho/godotenv"

)
/*-----------------anilist querry--------------------*/
const query = `
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
var rdb *redis.Client
var ctx = context.Background()
func init() {

 /*-----------------setting up datasbase connection ---------------------*/
 godotenv.Load("../.env")


rdb = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URL"),
        Username: "default",
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })
}

var name_id map[string]int
var id_name map[int]string 
 /*-----------------setting up datasbase connection ---------------------*/



type pair struct {
    id    int
    score float64
}
// func getPage(w http.ResponseWriter ,r http.Request){

// }
/*------------function that  handles recommendation requests---------------*/
func getRecom(w http.ResponseWriter ,r *http.Request){
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    fmt.Println("processing")
    if r.Method != http.MethodPost {
         fmt.Println("error")
         fmt.Println("error")
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var data map[string]string
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        fmt.Println("error")
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    username,ok:= data["username"]
    if(!ok){
         fmt.Println("error")
    }
    fmt.Sprintf("user name which sen the request is %s",username)


var titles []string// arrays to store comon anime with user and databse 
var idarr []int

/*-------------------strting to process user data---------------*/
payload := map[string]interface{}{
    "query": query,
    "variables": map[string]string{
        "userName": username,
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
err = json.NewDecoder(resp.Body).Decode(&result)
if err != nil {
    log.Fatal(err)
}

for _, list := range result.Data.MediaListCollection.Lists {
    for _, entry := range list.Entries {

        if entry.Media.Title.English != "" {
            li:=normalize(entry.Media.Title.English)
            ind,ok:=name_id[li]
            if(ok){
            titles = append(titles, li)
            idarr=append(idarr,ind)
                }
        } else {
            
             li:=normalize( entry.Media.Title.Romaji)
            ind,ok:=name_id[li]
            if(ok){
            titles = append(titles, li)
            idarr=append(idarr,ind)
                }
        }
    }
}
/*-----------------ending of procesing user data---------------*/



/*-------------------------scoring---------------------*/

set := make(map[int]float64)
for _,ind:=range idarr{
    key := fmt.Sprintf("anime:%d", ind)
  vecBytes, _ := rdb.HGet(ctx, key, "embedding").Bytes()
   result, _ := rdb.FTSearchWithArgs(
    ctx,
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

for _, doc := range result.Docs {
  parts := strings.Split(doc.ID, ":")
    if len(parts) != 2 {
        fmt.Println("wrng")
        continue
    }

    numID, err := strconv.Atoi(parts[1])
    if err != nil {
        continue
    }

    // fmt.Println("Numeric ID:", numID)
    // fmt.Println(idname[numID])
    // fmt.Println("Score:", doc.Fields["score"])
   raw := doc.Fields["score"]

sim, err := strconv.ParseFloat(fmt.Sprintf("%v", raw), 64)
if err != nil {
    continue
}

set[numID] += sim


}



}


var arr []pair

for id, sc := range set {
    arr = append(arr, pair{id, sc})
}

// sort descending
sort.Slice(arr, func(i, j int) bool {
    return arr[i].score > arr[j].score
})

var reply []string

for i := 0; i < len(arr) && i < 20; i++ {
    name, ok := id_name[arr[i].id]
    if ok {
        reply = append(reply, 
            name)
    }
}

// for key, value := range set {
//     if value {
//         reply = append(reply, key)
//     }
// }




/*----------------------------------------------------*/

w.Header().Set("Content-Type", "application/json")

json.NewEncoder(w).Encode(reply)

}






func main(){
   
    
    /*-----------------setting anime id and name converisions----------------*/
     name_id,id_name=init_maps()
     

    ///http.HandleFunc("/",getPage)
    http.HandleFunc("/recommendations",getRecom)
  port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}

log.Fatal(http.ListenAndServe(":" + port, nil))
    
}






func init_maps() (map[string]int,map[int]string){
    file ,err:=os.Open(os.Getenv("CSV"))
    defer file.Close()
       if err!=nil{

     fmt.Println("error in making the maps ")
     return nil,nil
        }
    defer file.Close()
    reader:=csv.NewReader(file)

    headers,err:=reader.Read()
       if err!=nil{

     fmt.Println("error in making the maps ")
     return nil,nil
        }
    fmt.Println(headers[0])
    nameid:=make(map[string]int)
    idname:=make(map[int]string)

    for{
        record,err:=reader.Read()


        if err!=nil{

     // fmt.Println("error in making the maps ")
     // return nil,nil
     break
        }



        v:=cleanNpyString(record[0])

    id, err := strconv.Atoi(record[1])


     if err!=nil{

     fmt.Println("error in making the maps ")
     return nil,nil
        }


        nameid[normalize(v)]=id
          idname[id]=normalize(v)

        //fmt.Println(normalize(v))
//fmt.Println(record[1])
    }
fmt.Println("succesfully intialise the maps ")



    return nameid,idname
}

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