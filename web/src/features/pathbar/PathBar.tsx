import { SiGithub } from "@icons-pack/react-simple-icons";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";

interface PathBarProps {
  children?: ReactNode;
}

export function PathBar({ children }: PathBarProps) {
  const [name, setName] = useState("");

  useEffect(() => {
    fetch("/api/info")
      .then((res) => res.json())
      .then((data: { name: string }) => setName(data.name))
      .catch(() => {});
  }, []);

  return (
    <div className="flex items-center gap-2 border-b border-gray-200 bg-gray-50 px-4 py-2">
      {children}
      <span className="min-w-0 flex-1 truncate text-sm font-medium text-gray-800">
        {name}
      </span>
      {import.meta.env.VITE_GITHUB_URL && (
        <a
          href={import.meta.env.VITE_GITHUB_URL}
          target="_blank"
          rel="noopener noreferrer"
          className="shrink-0 rounded p-1 text-gray-500 hover:bg-gray-200"
          aria-label="GitHub"
        >
          <SiGithub className="size-4" />
        </a>
      )}
    </div>
  );
}
