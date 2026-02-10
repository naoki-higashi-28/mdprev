import { Check, Copy } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { type BundledLanguage, bundledLanguages, codeToHtml } from "shiki";

interface CodeBlockProps {
  code: string;
  language?: string;
}

function isSupportedLanguage(lang: string): lang is BundledLanguage {
  return lang in bundledLanguages;
}

export function CodeBlock({ code, language }: CodeBlockProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [copied, setCopied] = useState(false);
  const [highlighted, setHighlighted] = useState(false);

  useEffect(() => {
    let cancelled = false;

    const lang = language && isSupportedLanguage(language) ? language : "text";

    codeToHtml(code, {
      lang,
      theme: "github-light",
    }).then((html) => {
      if (!cancelled && containerRef.current) {
        containerRef.current.innerHTML = html;
        setHighlighted(true);
      }
    });

    return () => {
      cancelled = true;
    };
  }, [code, language]);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="not-prose group relative my-4">
      {language && (
        <span className="absolute left-3 top-2 text-xs text-gray-400 select-none">
          {language}
        </span>
      )}
      <button
        type="button"
        onClick={handleCopy}
        className="absolute right-2 top-2 rounded p-1.5 text-gray-400 opacity-0 transition-opacity hover:text-gray-600 hover:bg-gray-200 group-hover:opacity-100 cursor-pointer"
        title="Copy"
      >
        {copied ? <Check className="size-4" /> : <Copy className="size-4" />}
      </button>
      {!highlighted && (
        <pre className="rounded-lg bg-gray-50 p-4 pt-8 text-sm overflow-x-auto">
          <code>{code}</code>
        </pre>
      )}
      <div
        ref={containerRef}
        className={`[&_pre]:rounded-lg [&_pre]:p-4 [&_pre]:pt-8 [&_pre]:text-sm [&_pre]:overflow-x-auto ${highlighted ? "" : "hidden"}`}
      />
    </div>
  );
}
