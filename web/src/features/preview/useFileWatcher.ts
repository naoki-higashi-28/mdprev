import { useEffect, useRef, useState } from "react";

export function useFileWatcher(path: string): number {
  const [revision, setRevision] = useState(0);
  const pathRef = useRef(path);

  useEffect(() => {
    pathRef.current = path;
  }, [path]);

  useEffect(() => {
    const eventSource = new EventSource("/api/watch");

    eventSource.onmessage = (event) => {
      if (event.data === pathRef.current) {
        setRevision((r) => r + 1);
      }
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return revision;
}
