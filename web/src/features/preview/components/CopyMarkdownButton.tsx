import { Check, ClipboardCopy } from "lucide-react";
import { memo, useState } from "react";

interface CopyMarkdownButtonProps {
  content: string | null;
}

export const CopyMarkdownButton = memo(function CopyMarkdownButton({
  content,
}: CopyMarkdownButtonProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    if (content === null) return;
    await navigator.clipboard.writeText(content);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <button
      type="button"
      onClick={handleCopy}
      className="rounded p-1 text-gray-400 hover:text-gray-600 hover:bg-gray-100 cursor-pointer"
      title="Copy as Markdown"
    >
      {copied ? (
        <Check className="size-3.5" />
      ) : (
        <ClipboardCopy className="size-3.5" />
      )}
    </button>
  );
});
