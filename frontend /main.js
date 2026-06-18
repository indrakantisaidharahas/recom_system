const url=document.getElementById("AURL")
const val=String(url.val)
const button=document.getElementById("gert")
const res=document.getElementById("results")

button.onclick=()=>{
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




}