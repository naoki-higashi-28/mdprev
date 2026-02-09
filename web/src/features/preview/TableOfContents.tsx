import { useMemo } from "react";
import { slugify, stripTags } from "./slugify";

interface TocEntry {
  level: number;
  text: string;
  id: string;
}

interface TableOfContentsProps {
  markdown: string;
}

export function TableOfContents({ markdown }: TableOfContentsProps) {
  const headings = useMemo(() => {
    // Remove fenced code blocks before extracting headings
    const stripped = markdown
      .replace(/^(`{3,}).*\n[\s\S]*?^\1\s*$/gm, "")
      .replace(/^(~{3,}).*\n[\s\S]*?^\1\s*$/gm, "");
    const results: TocEntry[] = [];
    const regex = /^(#{1,6})\s+(.+)$/gm;
    let match = regex.exec(stripped);
    while (match) {
      results.push({
        level: match[1].length,
        text: stripTags(match[2]),
        id: slugify(match[2]),
      });
      match = regex.exec(stripped);
    }
    return results;
  }, [markdown]);

  if (headings.length === 0) return null;

  const minLevel = Math.min(...headings.map((h) => h.level));

  return (
    <nav className="text-sm">
      <ul className="space-y-1">
        {headings.map((heading, index) => (
          <li
            key={`${heading.id}-${index}`}
            style={{ paddingLeft: `${(heading.level - minLevel) * 0.75}rem` }}
          >
            <a
              href={`#${heading.id}`}
              onClick={(e) => {
                e.preventDefault();
                const el = document.getElementById(heading.id);
                el?.scrollIntoView({ behavior: "smooth" });
              }}
              className="block py-0.5 text-gray-600 hover:text-gray-900 truncate"
            >
              {heading.text}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
}
