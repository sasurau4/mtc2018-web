/**
 * 文字を渡された文字数で省略する
 */
export const omitText = (text: string, maxLength: number = 20): string => {
  if (text.length > maxLength) {
    text = `${text.substr(0, maxLength)}...`;
  }
  return text;
};

/**
 * 渡されたものを指定された数だけ
 * 複製して配列で返却
 */
export function duplicate<T>(n: number) {
  return (item: T): T[] => {
    const list: T[] = [];
    for (let i = 0; i < n; i++) {
      list.push(item);
    }
    return list;
  };
}

/**
 * 渡された回数だけ関数を実行する
 */
export const times = (n: number) => {
  return (fn: (index: number) => void) => {
    for (let i = 0; i < n; i++) {
      fn(i);
    }
  };
};

/**
 * string[]を<br />でつなげる
 */
export function joinWithBr(texts: string[]): React.ReactNode[] {
  const result: React.ReactNode[] = [];
  for (let i = 0; i < texts.length; i++) {
    result.push(texts[i]);
    result.push(<br key={i} />);
  }

  result.splice(result.length - 1, 1);
  return result;
}

export function isJapan(lang: string): boolean {
  return lang === 'ja-JP' || lang === 'ja';
}

// https://stackoverflow.com/questions/105034/create-guid-uuid-in-javascript
export function generateUUID(): string {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
      .toString(16)
      .substring(1);
  }
  return (
    s4() +
    s4() +
    '-' +
    s4() +
    '-' +
    s4() +
    '-' +
    s4() +
    '-' +
    s4() +
    s4() +
    s4()
  );
}
