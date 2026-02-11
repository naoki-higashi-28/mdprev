export class TreeEntry {
  readonly type: "dir" | "file";
  readonly name: string;
  readonly path: string;
  readonly ext?: string;
  readonly size?: number;
  readonly mtime?: number;

  constructor(params: {
    type: "dir" | "file";
    name: string;
    path: string;
    ext?: string;
    size?: number;
    mtime?: number;
  }) {
    this.type = params.type;
    this.name = params.name;
    this.path = params.path;
    this.ext = params.ext;
    this.size = params.size;
    this.mtime = params.mtime;
  }
}

export class Tree {
  readonly path: string;
  readonly entries: TreeEntry[];

  constructor(params: { path: string; entries: TreeEntry[] }) {
    this.path = params.path;
    this.entries = params.entries;
  }
}

export class SearchResult {
  readonly query: string;
  readonly entries: TreeEntry[];

  constructor(params: { query: string; entries: TreeEntry[] }) {
    this.query = params.query;
    this.entries = params.entries;
  }
}
