
const button=document.getElementById("gert")
const res=document.getElementById("results")
const res2=document.getElementById("search_results")
const db= await ind_db_builder("wizard","gringots","./names_to_id.json")
const storeName="gringots"


button.onclick=()=>{
const url=document.getElementById("AURL")
const val=String(url.value)

    
   
res.innerHTML="your recommendations are "+"<br>"

	fetch('https://recom-system-bm2h.onrender.com/recommendations', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    "username":val
  })

}) .then(response => response.json())
  .then(data => {
    // If data is an array:
    
  data.forEach(item => {
    res.innerHTML += item + "<br>";
});
    

  })
console.log("the values is ")
console.log(val)


}
const button2=document.getElementById("search")
button2.onclick=()=>{
  const inp=document.getElementById("user_search")
  let num=30;
  const val=String(inp.value)
    const tx = db.transaction(storeName, "readonly");
  const readStore = tx.objectStore(storeName);
  const index = readStore.index("accio");
  res2.innerHTML="your search results  are "+"<br>"
  API.search(index, val, "en").then(async (res) => {
    console.log(res);

    for (const [id, score] of res) {
      const req = readStore.get(id);

      const anime = await new Promise((resolve, reject) => {
        req.onsuccess = () => resolve(req.result);
        req.onerror = () => reject(req.error);
      });

      console.log(anime.text, score);
      if(num>0){
      res2.innerHTML+="<br>" +anime.text;
      num=num-1;
    }
    }


});

}