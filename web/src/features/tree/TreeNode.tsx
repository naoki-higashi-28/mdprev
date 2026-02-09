import { ChevronDown, ChevronRight, File, Folder } from "lucide-react";
import { useState } from "react";
import type { TreeEntry } from "../../shared/types";
import { fetchTree } from "./api";

interface TreeNodeProps {
  entry: TreeEntry;
  depth: number;
  selectedPath: string;
  onSelectFile: (path: string) => void;
}

export function TreeNode({
  entry,
  depth,
  selectedPath,
  onSelectFile,
}: TreeNodeProps) {
  const [expanded, setExpanded] = useState(false);
  const [children, setChildren] = useState<TreeEntry[]>([]);
  const [loaded, setLoaded] = useState(false);

  const paddingLeft = `${depth * 16 + 8}px`;

  if (entry.type === "dir") {
    const handleToggle = async () => {
      if (!loaded) {
        try {
          const data = await fetchTree(entry.path);
          setChildren(data.entries ?? []);
          setLoaded(true);
        } catch {
          // Silently handle error - directory will appear empty
        }
      }
      setExpanded((prev) => !prev);
    };

    return (
      <div>
        <button
          type="button"
          onClick={handleToggle}
          className="flex w-full items-center gap-1 py-0.5 text-left text-sm hover:bg-gray-100"
          style={{ paddingLeft }}
        >
          {expanded ? (
            <ChevronDown className="size-4 shrink-0 text-gray-500" />
          ) : (
            <ChevronRight className="size-4 shrink-0 text-gray-500" />
          )}
          <Folder className="size-4 shrink-0 text-gray-500" />
          <span className="truncate text-gray-700">{entry.name}</span>
        </button>
        {expanded && (
          <div>
            {children.map((child) => (
              <TreeNode
                key={child.path}
                entry={child}
                depth={depth + 1}
                selectedPath={selectedPath}
                onSelectFile={onSelectFile}
              />
            ))}
          </div>
        )}
      </div>
    );
  }

  const isMarkdown = entry.ext === "md" || entry.ext === "markdown";
  if (!isMarkdown) {
    return null;
  }

  const isSelected = entry.path === selectedPath;

  return (
    <button
      type="button"
      onClick={() => onSelectFile(entry.path)}
      className={`flex w-full items-center gap-1 py-0.5 text-left text-sm ${
        isSelected
          ? "bg-blue-100 font-medium text-blue-800"
          : "text-gray-700 hover:bg-gray-100"
      }`}
      style={{ paddingLeft }}
    >
      <File className="size-4 shrink-0 text-gray-400" />
      <span className="truncate">{entry.name}</span>
    </button>
  );
}
