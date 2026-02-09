/** Strip inline markdown syntax and HTML tags to get plain text. */
function stripMarkup(text: string): string {
  return (
    text
      // HTML tags: <br>, <br/>, <span class="x">...</span>, etc.
      .replace(/<[^>]+>/g, " ")
      // images: ![alt](url) -> alt
      .replace(/!\[([^\]]*)\]\([^)]*\)/g, "$1")
      // links: [text](url) -> text
      .replace(/\[([^\]]*)\]\([^)]*\)/g, "$1")
      // inline code: `code` -> code
      .replace(/`([^`]+)`/g, "$1")
      // bold/italic: ***text***, **text**, *text*, ___text___, __text__, _text_
      .replace(/\*{1,3}([^*]+)\*{1,3}/g, "$1")
      .replace(/_{1,3}([^_]+)_{1,3}/g, "$1")
      // strikethrough: ~~text~~ -> text
      .replace(/~~([^~]+)~~/g, "$1")
  );
}

/** Strip HTML tags and inline markdown for display. */
export function stripTags(text: string): string {
  return stripMarkup(text).replace(/\s+/g, " ").trim();
}

export function slugify(text: string): string {
  return stripMarkup(text)
    .toLowerCase()
    .replace(/[^\p{L}\p{N}\s_-]/gu, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");
}
