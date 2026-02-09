import { PanelLeftClose, PanelLeftOpen } from "lucide-react";
import { useState } from "react";
import { PathBar } from "../features/pathbar/PathBar";
import { MarkdownPreview } from "../features/preview/MarkdownPreview";
import { TreeView } from "../features/tree/TreeView";
import { useSelectedPath } from "../shared/hooks/useSelectedPath";

export function HomePage() {
  const [selectedPath, setSelectedPath] = useSelectedPath();
  const [sidebarOpen, setSidebarOpen] = useState(true);

  return (
    <div className="flex h-screen flex-col">
      <PathBar>
        <button
          type="button"
          onClick={() => setSidebarOpen((prev) => !prev)}
          className="shrink-0 rounded p-1 text-gray-500 hover:bg-gray-200"
          aria-label={sidebarOpen ? "Close sidebar" : "Open sidebar"}
        >
          {sidebarOpen ? (
            <PanelLeftClose className="size-4" />
          ) : (
            <PanelLeftOpen className="size-4" />
          )}
        </button>
      </PathBar>
      <div className="flex min-h-0 flex-1">
        {sidebarOpen && (
          <div className="w-64 shrink-0 border-r border-gray-200 bg-white">
            <TreeView
              selectedPath={selectedPath}
              onSelectFile={setSelectedPath}
            />
          </div>
        )}
        <div className="flex-1 overflow-y-auto bg-white">
          <MarkdownPreview
            filePath={selectedPath}
            onNavigate={setSelectedPath}
          />
        </div>
      </div>
    </div>
  );
}
