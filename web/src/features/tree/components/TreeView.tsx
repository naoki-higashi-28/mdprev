import { useState } from "react";
import { useSearch } from "../hooks/useSearch";
import { useTree } from "../hooks/useTree";
import { SearchInput } from "./SearchInput";
import { SearchResults } from "./SearchResults";
import { TreeNode } from "./TreeNode";

interface TreeViewProps {
  selectedPath: string;
  onSelectFile: (path: string) => void;
}

export function TreeView({ selectedPath, onSelectFile }: TreeViewProps) {
  const { tree, error, loading } = useTree();
  const [searchQuery, setSearchQuery] = useState("");
  const searchResult = useSearch(searchQuery);

  return (
    <div className="flex h-full flex-col">
      <div className="shrink-0 border-b border-gray-200 p-2">
        <SearchInput value={searchQuery} onChange={setSearchQuery} />
      </div>
      <div className="flex-1 overflow-y-auto">
        {searchQuery ? (
          searchResult ? (
            <SearchResults
              entries={searchResult.entries}
              selectedPath={selectedPath}
              onSelectFile={onSelectFile}
            />
          ) : (
            <div className="p-4 text-sm text-gray-500">Searching...</div>
          )
        ) : loading ? (
          <div className="p-4 text-sm text-gray-500">Loading...</div>
        ) : error ? (
          <div className="p-4 text-sm text-red-600">{error}</div>
        ) : !tree || tree.entries.length === 0 ? (
          <div className="p-4 text-sm text-gray-500">No files found</div>
        ) : (
          <div className="py-2">
            {tree.entries.map((entry) => (
              <TreeNode
                key={entry.path}
                entry={entry}
                depth={0}
                selectedPath={selectedPath}
                onSelectFile={onSelectFile}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
