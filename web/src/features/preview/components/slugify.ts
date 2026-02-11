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
    .replace(/[^\p{L}\p{N}\p{Emoji}\s_-]/gu, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");
}

export function createSlugifier(): (text: string) => string {
  const counts = new Map<string, number>();
  return (text: string) => {
    const base = slugify(text);
    const count = counts.get(base) ?? 0;
    counts.set(base, count + 1);
    return count === 0 ? base : `${base}-${count}`;
  };
}

interface HastNode {
  type: string;
  value?: string;
  tagName?: string;
  properties?: Record<string, unknown>;
  children?: HastNode[];
}

function getHastText(node: HastNode): string {
  if (node.type === "text") return node.value ?? "";
  if (node.children) return node.children.map(getHastText).join("");
  return "";
}

export function rehypeHeadingIds() {
  return (tree: HastNode) => {
    const slugify = createSlugifier();

    function walk(node: HastNode) {
      if (node.type === "element" && /^h[1-6]$/.test(node.tagName ?? "")) {
        if (!node.properties) node.properties = {};
        node.properties.id = slugify(getHastText(node));
      }
      if (node.children) {
        for (const child of node.children) {
          walk(child);
        }
      }
    }

    walk(tree);
  };
}
