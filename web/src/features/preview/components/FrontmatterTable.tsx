import { ChevronDown, ChevronRight, Info } from "lucide-react";
import { useState } from "react";

interface FrontmatterTableProps {
  entries: [string, unknown][];
}

export function FrontmatterTable({ entries }: FrontmatterTableProps) {
  const [open, setOpen] = useState(false);

  if (entries.length === 0) return null;

  return (
    <div className="not-prose mb-6">
      <button
        type="button"
        onClick={() => setOpen((prev) => !prev)}
        className="flex items-center gap-1 text-gray-400 hover:text-gray-600 cursor-pointer"
        title="Frontmatter"
      >
        <Info className="size-4" />
        {open ? (
          <ChevronDown className="size-3.5" />
        ) : (
          <ChevronRight className="size-3.5" />
        )}
      </button>
      {open && (
        <table className="mt-2 w-full border-collapse text-sm">
          <tbody>
            {entries.map(([key, value]) => (
              <tr key={key} className="border-b border-gray-200">
                <td className="py-1.5 pr-4 font-medium text-gray-500">{key}</td>
                <td className="py-1.5 text-gray-700">{String(value)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
