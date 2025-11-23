export type KeyValue = {
  key: string;
  value: string;
};

export function mapToKeyValueArray(map: Record<string, string>): KeyValue[] {
  let array: KeyValue[] = [];

  for (const [key, value] of Object.entries(map)) {
    array.push({ key, value });
  }

  array = array.sort((a, b) => a.key.localeCompare(b.key));

  return array;
}
