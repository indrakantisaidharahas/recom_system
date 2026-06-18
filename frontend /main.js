const url=document.getElementById("AURL")
const val=url.val
const button=document.getElementById("gert")
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
    data.forEach(item => console.log(item));
    
    // If data is an object containing an array (e.g., { "results": [...] }):
    data.results.forEach(item => console.log(item));
  })




}