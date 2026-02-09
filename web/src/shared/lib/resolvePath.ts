export function resolveRelativePath(
  currentFilePath: string,
  relativePath: string,
): string {
  if (
    relativePath.startsWith("http://") ||
    relativePath.startsWith("https://")
  ) {
    return relativePath;
  }

  const dir = currentFilePath.includes("/")
    ? currentFilePath.substring(0, currentFilePath.lastIndexOf("/"))
    : "";

  const parts = dir ? dir.split("/") : [];
  for (const segment of relativePath.split("/")) {
    if (segment === "..") {
      parts.pop();
    } else if (segment !== "." && segment !== "") {
      parts.push(segment);
    }
  }

  return parts.join("/");
}
