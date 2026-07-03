
const button=document.getElementById("gert")
const res=document.getElementById("results")
ind_db_builder("wizard","gringots","./names_to_id.json")
button.onclick=()=>{
const url=document.getElementById("AURL")
const val=String(url.value)

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
