export async function fetchFile(path: string): Promise<string> {
  const params = new URLSearchParams({ path });
  const res = await fetch(`/api/file?${params.toString()}`);
  if (!res.ok) {
    if (res.status === 404) {
      throw new FileNotFoundError("File not found");
    }
    throw new Error(`Failed to fetch file: ${res.status}`);
  }
  return res.text();
}

export class FileNotFoundError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "FileNotFoundError";
  }
}
