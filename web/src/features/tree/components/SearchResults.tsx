import { File } from "lucide-react";
import type { TreeEntry } from "../model/tree.model";

interface SearchResultsProps {
  entries: TreeEntry[];
  selectedPath: string;
  onSelectFile: (path: string) => void;
}

export function SearchResults({
  entries,
  selectedPath,
  onSelectFile,
}: SearchResultsProps) {
  if (entries.length === 0) {
    return <div className="p-4 text-sm text-gray-500">No results found</div>;
  }

  return (
    <div className="py-2">
      {entries.map((entry) => {
        const isSelected = entry.path === selectedPath;
        const dir = entry.path.includes("/")
          ? entry.path.slice(0, entry.path.lastIndexOf("/"))
          : "";

        return (
          <button
            type="button"
            key={entry.path}
            onClick={() => onSelectFile(entry.path)}
            className={`flex w-full flex-col px-3 py-1 text-left text-sm ${
              isSelected
                ? "bg-blue-100 text-blue-800"
                : "text-gray-700 hover:bg-gray-100"
            }`}
          >
            <span className="flex items-center gap-1">
              <File className="size-4 shrink-0 text-gray-400" />
              <span className={`truncate ${isSelected ? "font-medium" : ""}`}>
                {entry.name}
              </span>
            </span>
            {dir && (
              <span className="ml-5 truncate text-xs text-gray-400">{dir}</span>
            )}
          </button>
        );
      })}
    </div>
  );
}
