import mermaid from "mermaid";
import { useEffect, useId, useRef, useState } from "react";

mermaid.initialize({
  startOnLoad: false,
  theme: "default",
});

interface MermaidProps {
  chart: string;
}

export function Mermaid({ chart }: MermaidProps) {
  const id = useId().replace(/:/g, "m");
  const containerRef = useRef<HTMLDivElement>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const { svg } = await mermaid.render(`mermaid-${id}`, chart);
        if (!cancelled && containerRef.current) {
          containerRef.current.innerHTML = svg;
          setError(null);
        }
      } catch (e) {
        if (!cancelled) {
          setError(e instanceof Error ? e.message : String(e));
        }
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [chart, id]);

  if (error) {
    return (
      <pre className="not-prose rounded bg-red-50 p-4 text-sm text-red-700 overflow-x-auto">
        <code>{chart}</code>
        <p className="mt-2 font-semibold">Mermaid error: {error}</p>
      </pre>
    );
  }

  return (
    <div ref={containerRef} className="not-prose my-4 flex justify-center" />
  );
}
