export interface TocEntry {
  level: number;
  text: string;
  id: string;
}

interface TableOfContentsProps {
  headings: TocEntry[];
  filePath: string;
  scrollRef: React.RefObject<HTMLDivElement | null>;
}

export function TableOfContents({
  headings,
  filePath,
  scrollRef,
}: TableOfContentsProps) {
  const fileName = filePath.split("/").pop() ?? filePath;
  const minLevel = Math.min(...headings.map((h) => h.level));

  return (
    <nav className="text-sm">
      <a
        href="#top"
        onClick={(e) => {
          e.preventDefault();
          scrollRef.current?.scrollTo({ top: 0, behavior: "smooth" });
        }}
        className="block py-0.5 font-semibold text-gray-800 hover:text-gray-600 truncate mb-1"
        title={filePath}
      >
        {fileName}
      </a>
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
