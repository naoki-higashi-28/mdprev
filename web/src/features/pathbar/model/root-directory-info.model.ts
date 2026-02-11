export class RootDirectoryInfo {
  readonly rootDirectoryName: string;

  constructor(params: { rootDirectoryName: string }) {
    this.rootDirectoryName = params.rootDirectoryName;
  }
}
