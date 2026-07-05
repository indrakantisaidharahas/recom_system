async function ind_db_builder(dbName, storeName, jsonPath) {
  const response = await fetch(jsonPath);

  if (!response.ok) {
    throw new Error(`Failed to load JSON: ${response.status} ${response.statusText}`);
  }

  const data = await response.json();

  return new Promise((resolve, reject) => {

    const dbRequest = indexedDB.open(dbName, 1);

    dbRequest.onupgradeneeded = (event) => {
      const db = event.target.result;
      const store = db.createObjectStore(storeName, { keyPath: "docid" });
      store.createIndex("accio", "terms", { multiEntry: true });
    };

    dbRequest.onerror = () => reject(dbRequest.error);

    dbRequest.onsuccess = (event) => {
      const db = event.target.result;

      const transaction = db.transaction(storeName, "readwrite");
      const store = transaction.objectStore(storeName);

      let ind = 0;

      for (const value of Object.values(data.name)) {
        store.put({
          docid: ind++,
          text: value,
          terms: API.tokenize(value, "en")
        });
      }

      transaction.oncomplete = () => {
        console.log("Import successful!");
        resolve(db);
      };
    };

  });
}