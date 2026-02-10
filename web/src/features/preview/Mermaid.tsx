import { Maximize2 } from "lucide-react";
import mermaid from "mermaid";
import { memo, useEffect, useId, useRef, useState } from "react";
import { MermaidModal } from "./MermaidModal";

mermaid.initialize({
  startOnLoad: false,
  theme: "default",
});

interface MermaidProps {
  chart: string;
}

export const Mermaid = memo(function Mermaid({ chart }: MermaidProps) {
  const id = useId().replace(/:/g, "m");
  const containerRef = useRef<HTMLDivElement>(null);
  const [error, setError] = useState<string | null>(null);
  const [svgHtml, setSvgHtml] = useState<string>("");
  const [modalOpen, setModalOpen] = useState(false);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const { svg } = await mermaid.render(`mermaid-${id}`, chart);
        if (!cancelled && containerRef.current) {
          containerRef.current.innerHTML = svg;
          setSvgHtml(svg);
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
    <>
      <button
        type="button"
        className="not-prose group relative my-4 flex w-full cursor-pointer justify-center border-none bg-transparent p-0"
        onClick={() => setModalOpen(true)}
      >
        <div ref={containerRef} />
        <div className="absolute top-2 right-2 rounded-lg bg-white/80 p-1.5 text-gray-500 opacity-0 shadow transition-opacity group-hover:opacity-100">
          <Maximize2 size={16} />
        </div>
      </button>
      <MermaidModal
        open={modalOpen}
        onOpenChange={setModalOpen}
        svgHtml={svgHtml}
      />
    </>
  );
});
