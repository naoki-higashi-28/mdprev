import { Check, ClipboardCopy, List } from "lucide-react";
import {
  type ComponentPropsWithoutRef,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import Markdown from "react-markdown";
import rehypeKatex from "rehype-katex";
import rehypeRaw from "rehype-raw";
import remarkFrontmatter from "remark-frontmatter";
import { remarkAlert } from "remark-github-blockquote-alert";
import "remark-github-blockquote-alert/alert.css";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import { parse as parseYaml } from "yaml";
import { resolveRelativePath } from "../../shared/lib/resolvePath";
import { CodeBlock } from "./CodeBlock";
import { FrontmatterTable } from "./FrontmatterTable";
import { Mermaid } from "./Mermaid";
import { rehypeHeadingIds } from "./slugify";
import type { TocEntry } from "./TableOfContents";
import { TableOfContents } from "./TableOfContents";
import { useFileContent } from "./useFileContent";

interface MarkdownPreviewProps {
  filePath: string;
  onNavigate: (path: string) => void;
}

export function MarkdownPreview({
  filePath,
  onNavigate,
}: MarkdownPreviewProps) {
  const { content, error, loading } = useFileContent(filePath);
  const [tocVisible, setTocVisible] = useState(true);
  const [copied, setCopied] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const [headings, setHeadings] = useState<TocEntry[]>([]);

  const handleCopyMarkdown = async () => {
    if (content === null) return;
    await navigator.clipboard.writeText(content);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const parsed = useMemo(() => {
    if (content === null) return null;
    const match = content.match(/^---\r?\n([\s\S]*?)\r?\n---\r?\n([\s\S]*)$/);
    if (!match) return { entries: [], body: content };
    const data = parseYaml(match[1]);
    const entries =
      data && typeof data === "object" ? Object.entries(data) : [];
    return { entries, body: match[2] };
  }, [content]);

  useLayoutEffect(() => {
    if (!content || !contentRef.current) return;
    const els = contentRef.current.querySelectorAll("h1, h2, h3, h4, h5, h6");
    const entries: TocEntry[] = [];
    for (const el of els) {
      if (!el.id) continue;
      entries.push({
        level: Number.parseInt(el.tagName[1], 10),
        text: el.textContent ?? "",
        id: el.id,
      });
    }
    setHeadings(entries);
  }, [content]);

  if (!filePath) {
    return (
      <div className="flex h-full items-center justify-center text-gray-400">
        Select a file
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center text-gray-500">
        Loading...
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-full items-center justify-center text-red-600">
        {error}
      </div>
    );
  }

  if (parsed === null) {
    return null;
  }

  return (
    <div className="relative flex h-full">
      <div className="relative flex-1 min-w-0">
        <button
          type="button"
          onClick={() => setTocVisible((prev) => !prev)}
          className="absolute right-4 top-4 z-10 rounded p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 cursor-pointer"
          title={tocVisible ? "Hide TOC" : "Show TOC"}
        >
          <List className="size-5" />
        </button>
        <div ref={scrollRef} className="h-full overflow-y-auto p-6">
          <div className="mb-4 flex items-center gap-1.5">
            <p className="text-sm text-gray-500">{filePath}</p>
            <button
              type="button"
              onClick={handleCopyMarkdown}
              className="rounded p-1 text-gray-400 hover:text-gray-600 hover:bg-gray-100 cursor-pointer"
              title="Copy as Markdown"
            >
              {copied ? (
                <Check className="size-3.5" />
              ) : (
                <ClipboardCopy className="size-3.5" />
              )}
            </button>
          </div>
          <div ref={contentRef} className="prose w-full max-w-4xl mx-auto">
            <FrontmatterTable entries={parsed.entries} />
            <Markdown
              remarkPlugins={[
                remarkFrontmatter,
                remarkGfm,
                remarkMath,
                remarkAlert,
              ]}
              rehypePlugins={[rehypeRaw, rehypeKatex, rehypeHeadingIds]}
              components={{
                pre: ({ children }: ComponentPropsWithoutRef<"pre">) => {
                  if (
                    children &&
                    typeof children === "object" &&
                    "type" in children &&
                    children.type === "code"
                  ) {
                    const codeProps = children.props as {
                      className?: string;
                      children?: string;
                    };
                    const match = codeProps.className?.match(/language-(\w+)/);
                    const lang = match?.[1];
                    const code = String(codeProps.children ?? "").replace(
                      /\n$/,
                      "",
                    );

                    if (lang === "mermaid") {
                      return <Mermaid chart={code} />;
                    }

                    return <CodeBlock code={code} language={lang} />;
                  }
                  return <pre>{children}</pre>;
                },
                img: ({
                  src,
                  alt,
                  ...props
                }: ComponentPropsWithoutRef<"img">) => {
                  const resolvedSrc =
                    src &&
                    !src.startsWith("http://") &&
                    !src.startsWith("https://")
                      ? `/raw/${resolveRelativePath(filePath, src)}`
                      : src;
                  return (
                    <img
                      src={resolvedSrc}
                      alt={alt ?? ""}
                      className="max-w-full"
                      {...props}
                    />
                  );
                },
                a: ({
                  href,
                  children,
                  ...props
                }: ComponentPropsWithoutRef<"a">) => {
                  if (
                    href &&
                    !href.startsWith("http://") &&
                    !href.startsWith("https://") &&
                    !href.startsWith("#") &&
                    (href.endsWith(".md") || href.endsWith(".markdown"))
                  ) {
                    const resolvedPath = resolveRelativePath(filePath, href);
                    return (
                      <a
                        href={`/${resolvedPath}`}
                        onClick={(e) => {
                          e.preventDefault();
                          onNavigate(resolvedPath);
                        }}
                        {...props}
                      >
                        {children}
                      </a>
                    );
                  }

                  if (href?.startsWith("#")) {
                    return (
                      <a
                        href={href}
                        onClick={(e) => {
                          e.preventDefault();
                          const targetId = href.slice(1);
                          const el = document.getElementById(targetId);
                          if (el && scrollRef.current) {
                            const containerRect =
                              scrollRef.current.getBoundingClientRect();
                            const elementRect = el.getBoundingClientRect();
                            const offsetTop =
                              elementRect.top -
                              containerRect.top +
                              scrollRef.current.scrollTop;
                            scrollRef.current.scrollTo({
                              top: offsetTop,
                              behavior: "smooth",
                            });
                          }
                        }}
                        {...props}
                      >
                        {children}
                      </a>
                    );
                  }

                  const isExternal =
                    href?.startsWith("http://") || href?.startsWith("https://");
                  return (
                    <a
                      href={href}
                      {...(isExternal
                        ? { target: "_blank", rel: "noopener noreferrer" }
                        : {})}
                      {...props}
                    >
                      {children}
                    </a>
                  );
                },
              }}
            >
              {parsed.body}
            </Markdown>
          </div>
        </div>
      </div>
      {tocVisible && (
        <aside className="w-64 shrink-0 border-l border-gray-200 overflow-y-auto p-4">
          <TableOfContents
            headings={headings}
            filePath={filePath}
            scrollRef={scrollRef}
          />
        </aside>
      )}
    </div>
  );
}
