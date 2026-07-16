self.API = (() => {

  function tokenize(text, locale) {
    const words = new Set();
    const segmenter = new Intl.Segmenter(locale, { granularity: 'word' });

    for (let { segment, isWordLike } of segmenter.segment(text)) {
      if (isWordLike) {
        let word = segment.toLowerCase();
        let tw="";
        let ind=0;
        // word = stemmer(word);
       if (word.length >= 3) {
          for (let i = 0; i <= word.length - 3; i++) {
            words.add(word.substring(i, i + 3));
          }
        }
        words.add(word);
      }
    }

    return Array.from(words);
  }
    function t1(text, locale) {
    const words = new Set();
    const segmenter = new Intl.Segmenter(locale, { granularity: 'word' });

    for (let { segment, isWordLike } of segmenter.segment(text)) {
      if (isWordLike) {
        let word = segment.toLowerCase();
        let tw="";
        let ind=0;
        // word = stemmer(word);
      //  if (word.length >= 3) {
      //     for (let i = 0; i <= word.length - 3; i++) {
      //       words.add(word.substring(i, i + 3));
      //     }
      //   }
        words.add(word);
      }
    }

    return Array.from(words);
  }
     function t2(text, locale) {
    const words = new Set();
    const segmenter = new Intl.Segmenter(locale, { granularity: 'word' });

    for (let { segment, isWordLike } of segmenter.segment(text)) {
      if (isWordLike) {
        let word = segment.toLowerCase();
        let tw="";
        let ind=0;
        // word = stemmer(word);
       if (word.length >= 3) {
          for (let i = 0; i <= word.length - 3; i++) {
            words.add(word.substring(i, i + 3));
          }
        }
        //words.add(word);
      }
    }

    return Array.from(words);
  }


  async function search(index, query, locale) {

    const terms = t1(query, locale);
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

     const terms2 = t2(query, locale);
   
    for (const term of terms2) {
      const request = index.openKeyCursor(term);
      await new Promise(resolve => {
        request.onsuccess = () => {
          const cursor = request.result;

          if (!cursor) {
            resolve();
            return;
          }
          const id = cursor.primaryKey;
          map.set(id, (map.get(id) || 0) + 0.5);
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