
async function ind_db_builder(dbName, storeName, jsonPath) {
 const response = await fetch(jsonPath);

  if (!response.ok) {
    throw new Error(`Failed to load JSON: ${response.status} ${response.statusText}`);
  }

  const data = await response.json();

  const dbRequest = indexedDB.open(dbName, 1);

  dbRequest.onupgradeneeded = (event) => {
    const db = event.target.result;
    const store = db.createObjectStore(storeName, { keyPath: 'docid' });
    store.createIndex('accio', 'terms', { multiEntry: true });
  };

  dbRequest.onsuccess = (event) => {
    console.log("connection to db opened");

    const db = event.target.result;
    const transaction = db.transaction([storeName], "readwrite");
    const store = transaction.objectStore(storeName);

    let ind = 0;

    for (const value of Object.values(data.name)) {
      store.put({
        docid: ind,
        text: value,
        terms: API.tokenize(value, 'en')
      });
      ind++;
    }

   transaction.oncomplete = () => {
  console.log("Import successful!");

  const tx = db.transaction(storeName, "readonly");
  const readStore = tx.objectStore(storeName);
  const index = readStore.index("accio");

  API.search(index, "My Hero Academia", "en").then(async (res) => {
    console.log(res);

    for (const [id, score] of res) {
      const req = readStore.get(id);

      const anime = await new Promise((resolve, reject) => {
        req.onsuccess = () => resolve(req.result);
        req.onerror = () => reject(req.error);
      });

      console.log(anime.text, score);
    }

    db.close();
  });
};
  };
}