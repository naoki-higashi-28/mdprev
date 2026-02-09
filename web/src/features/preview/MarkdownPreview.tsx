import { List } from "lucide-react";
import {
  type ComponentPropsWithoutRef,
  type ReactNode,
  useMemo,
  useState,
} from "react";
import Markdown from "react-markdown";
import rehypeKatex from "rehype-katex";
import rehypeRaw from "rehype-raw";
import remarkFrontmatter from "remark-frontmatter";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import { parse as parseYaml } from "yaml";
import { resolveRelativePath } from "../../shared/lib/resolvePath";
import { FrontmatterTable } from "./FrontmatterTable";
import { slugify } from "./slugify";
import { TableOfContents } from "./TableOfContents";
import { useFileContent } from "./useFileContent";

interface MarkdownPreviewProps {
  filePath: string;
  onNavigate: (path: string) => void;
}

function extractText(node: ReactNode): string {
  if (typeof node === "string") return node;
  if (typeof node === "number") return String(node);
  if (!node) return "";
  if (Array.isArray(node)) return node.map(extractText).join("");
  if (typeof node === "object" && "props" in node) {
    const props = node.props as { children?: ReactNode };
    return extractText(props.children);
  }
  return "";
}

type HeadingTag = "h1" | "h2" | "h3" | "h4" | "h5" | "h6";

function createHeadingRenderer(Tag: HeadingTag) {
  return function HeadingComponent({
    children,
    ...props
  }: ComponentPropsWithoutRef<typeof Tag>) {
    const id = slugify(extractText(children));
    return (
      <Tag id={id} {...props}>
        {children}
      </Tag>
    );
  };
}

const headingComponents = {
  h1: createHeadingRenderer("h1"),
  h2: createHeadingRenderer("h2"),
  h3: createHeadingRenderer("h3"),
  h4: createHeadingRenderer("h4"),
  h5: createHeadingRenderer("h5"),
  h6: createHeadingRenderer("h6"),
};

export function MarkdownPreview({
  filePath,
  onNavigate,
}: MarkdownPreviewProps) {
  const { content, error, loading } = useFileContent(filePath);
  const [tocVisible, setTocVisible] = useState(true);

  const parsed = useMemo(() => {
    if (content === null) return null;
    const match = content.match(/^---\r?\n([\s\S]*?)\r?\n---\r?\n([\s\S]*)$/);
    if (!match) return { entries: [], body: content };
    const data = parseYaml(match[1]);
    const entries =
      data && typeof data === "object" ? Object.entries(data) : [];
    return { entries, body: match[2] };
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
        <div className="h-full overflow-y-auto p-6">
          <p className="mb-4 text-sm text-gray-500">{filePath}</p>
          <div className="prose w-full max-w-4xl mx-auto">
            <FrontmatterTable entries={parsed.entries} />
            <Markdown
              remarkPlugins={[remarkFrontmatter, remarkGfm, remarkMath]}
              rehypePlugins={[rehypeRaw, rehypeKatex]}
              components={{
                ...headingComponents,
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

                  const isExternal =
                    href?.startsWith("http://") || href?.startsWith("https://");
                  return (
                    <a
                      href={href}
                      className=""
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
          <TableOfContents markdown={parsed.body} />
        </aside>
      )}
    </div>
  );
}
