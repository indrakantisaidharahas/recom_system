self.API = (() => {

  function tokenize(text, locale) {
    const words = new Set();
    const segmenter = new Intl.Segmenter(locale, { granularity: 'word' });

    for (let { segment, isWordLike } of segmenter.segment(text)) {
      if (isWordLike) {
        let word = segment.toLowerCase();
        // word = stemmer(word);
        words.add(word);
      }
    }

    return Array.from(words);
  }

  async function search(index, query, locale) {

    const terms = tokenize(query, locale);
    const map = new Map();
    for (const term of terms) {
      const request = index.openKeyCursor(term);
      await new Promise(resolve => {
        request.onsuccess = () => {
          const cursor = request.result;

          if (!cursor) {
            resolve();
            return;
          }
          const id = cursor.primaryKey;
          map.set(id, (map.get(id) || 0) + 1);
          cursor.continue();
        };

      });

    }

    return [...map.entries()].sort((a, b) => b[1] - a[1]);

  }

  return {
    tokenize,
    search
  };

})();