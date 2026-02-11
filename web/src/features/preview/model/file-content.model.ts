export class FileContent {
  readonly value: string;

  constructor(value: string) {
    this.value = value;
  }
}

export class FileNotFoundError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "FileNotFoundError";
  }
}
