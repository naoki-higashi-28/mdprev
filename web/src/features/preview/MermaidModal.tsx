import * as Dialog from "@radix-ui/react-dialog";
import {
  Maximize2,
  Minimize2,
  RotateCcw,
  X,
  ZoomIn,
  ZoomOut,
} from "lucide-react";
import {
  type PointerEvent as ReactPointerEvent,
  type WheelEvent as ReactWheelEvent,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";

interface MermaidModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  svgHtml: string;
}

const MIN_SCALE = 0.25;
const MAX_SCALE = 5;
const ZOOM_STEP = 0.25;

export function MermaidModal({
  open,
  onOpenChange,
  svgHtml,
}: MermaidModalProps) {
  const [scale, setScale] = useState(1);
  const [translate, setTranslate] = useState({ x: 0, y: 0 });
  const [maximized, setMaximized] = useState(false);
  const dragging = useRef(false);
  const lastPointer = useRef({ x: 0, y: 0 });
  const areaRef = useRef<HTMLDivElement>(null);
  const diagramRef = useRef<HTMLDivElement>(null);

  const clampScale = useCallback(
    (s: number) => Math.min(MAX_SCALE, Math.max(MIN_SCALE, s)),
    [],
  );

  const fitToView = useCallback(
    (isMaximized: boolean) => {
      setTranslate({ x: 0, y: 0 });
      requestAnimationFrame(() => {
        const svg = diagramRef.current?.querySelector("svg");
        const area = areaRef.current;
        if (!svg || !area) {
          setScale(1);
          return;
        }
        const areaRect = area.getBoundingClientRect();
        const svgW = svg.width.baseVal?.value || svg.clientWidth;
        const svgH = svg.height.baseVal?.value || svg.clientHeight;
        if (svgW > 0 && svgH > 0) {
          const fit = Math.min(areaRect.width / svgW, areaRect.height / svgH);
          setScale(clampScale(isMaximized ? fit : fit * 0.8));
        } else {
          setScale(1);
        }
      });
    },
    [clampScale],
  );

  useEffect(() => {
    if (open) {
      fitToView(false);
    } else {
      setMaximized(false);
    }
  }, [open, fitToView]);

  const reset = useCallback(() => {
    fitToView(maximized);
  }, [fitToView, maximized]);

  const zoomIn = () => setScale((s) => clampScale(s + ZOOM_STEP));
  const zoomOut = () => setScale((s) => clampScale(s - ZOOM_STEP));

  const handleWheel = (e: ReactWheelEvent) => {
    e.stopPropagation();
    const delta = e.deltaY > 0 ? -ZOOM_STEP : ZOOM_STEP;
    setScale((s) => clampScale(s + delta));
  };

  const handlePointerDown = (e: ReactPointerEvent) => {
    if (e.button !== 0) return;
    dragging.current = true;
    lastPointer.current = { x: e.clientX, y: e.clientY };
    (e.target as HTMLElement).setPointerCapture(e.pointerId);
  };

  const handlePointerMove = (e: ReactPointerEvent) => {
    if (!dragging.current) return;
    const dx = e.clientX - lastPointer.current.x;
    const dy = e.clientY - lastPointer.current.y;
    lastPointer.current = { x: e.clientX, y: e.clientY };
    setTranslate((t) => ({ x: t.x + dx, y: t.y + dy }));
  };

  const handlePointerUp = () => {
    dragging.current = false;
  };

  const handleDoubleClick = () => {
    reset();
  };

  const controlBtnClass =
    "rounded-md p-1.5 text-gray-500 hover:bg-gray-100 hover:text-gray-700 transition-colors";

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-50 bg-black/40" />
        <Dialog.Content
          className={`fixed top-1/2 left-1/2 z-50 flex -translate-x-1/2 -translate-y-1/2 flex-col bg-white shadow-xl outline-none transition-all ${
            maximized
              ? "h-dvh w-dvw"
              : "h-[70dvh] w-[70dvw] rounded-xl border border-gray-200"
          }`}
          aria-describedby={undefined}
        >
          {/* Header */}
          <div className="flex items-center justify-between border-b border-gray-200 px-4 py-2">
            <Dialog.Title className="text-sm font-medium text-gray-600">
              Mermaid Diagram
            </Dialog.Title>
            <div className="flex gap-1">
              <button
                type="button"
                onClick={zoomOut}
                className={controlBtnClass}
                title="Zoom out"
              >
                <ZoomOut size={16} />
              </button>
              <button
                type="button"
                onClick={zoomIn}
                className={controlBtnClass}
                title="Zoom in"
              >
                <ZoomIn size={16} />
              </button>
              <button
                type="button"
                onClick={reset}
                className={controlBtnClass}
                title="Reset"
              >
                <RotateCcw size={16} />
              </button>
              <button
                type="button"
                onClick={() => setMaximized((v) => !v)}
                className={controlBtnClass}
                title={maximized ? "Restore" : "Maximize"}
              >
                {maximized ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
              </button>
              <Dialog.Close asChild>
                <button type="button" className={controlBtnClass} title="Close">
                  <X size={16} />
                </button>
              </Dialog.Close>
            </div>
          </div>

          {/* Diagram area */}
          <div
            ref={areaRef}
            className="flex flex-1 items-center justify-center overflow-hidden"
          >
            <div
              ref={diagramRef}
              role="application"
              className="cursor-grab select-none active:cursor-grabbing"
              style={{
                transform: `translate(${translate.x}px, ${translate.y}px) scale(${scale})`,
                transformOrigin: "center center",
              }}
              onWheel={handleWheel}
              onPointerDown={handlePointerDown}
              onPointerMove={handlePointerMove}
              onPointerUp={handlePointerUp}
              onDoubleClick={handleDoubleClick}
              // biome-ignore lint/security/noDangerouslySetInnerHtml: SVG rendered by mermaid
              dangerouslySetInnerHTML={{ __html: svgHtml }}
            />
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
