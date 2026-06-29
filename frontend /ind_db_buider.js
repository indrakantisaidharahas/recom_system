async function ind_db_builder(dbName, storeName, jsonPath) {
 const response = await fetch(jsonPath);
  
  if (!response.ok) {
    throw new Error(`Failed to load JSON: ${response.status} ${response.statusText}`);
  }
  
  const data = await response.json(); 
  // 2. Open Database
  const dbRequest = indexedDB.open(dbName, 1);
  
  dbRequest.onupgradeneeded = (event) => {
    const db = event.target.result;
    // Create store if it doesn't exist (keyPath depends on your JSON structure)
    if (!db.objectStoreNames.contains(storeName)) {
      db.createObjectStore(storeName, { keyPath: "id" }); 
    }
  };

  dbRequest.onsuccess = (event) => {
    const db = event.target.result;
    const transaction = db.transaction([storeName], "readwrite");
    const store = transaction.objectStore(storeName);

    // 3. Bulk Insert
    data.forEach(item => {
      store.put(item); // 'put' updates if key exists, 'add' throws error if duplicate
    });

    transaction.oncomplete = () => {
      console.log("Import successful!");
      db.close();
    };
  };
}   