import { Search, X } from "lucide-react";

interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
}

export function SearchInput({ value, onChange }: SearchInputProps) {
  return (
    <div className="relative flex items-center">
      <Search className="absolute left-2 size-4 text-gray-400" />
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Search files..."
        className="w-full rounded border border-gray-200 bg-gray-50 py-1 pr-7 pl-8 text-sm text-gray-700 placeholder-gray-400 outline-none focus:border-blue-300 focus:bg-white"
      />
      {value && (
        <button
          type="button"
          onClick={() => onChange("")}
          className="absolute right-1.5 rounded p-0.5 text-gray-400 hover:text-gray-600"
        >
          <X className="size-3.5" />
        </button>
      )}
    </div>
  );
}
