import { useCallback } from "react";
import { useLocation, useNavigate } from "react-router";

export function useSelectedPath(): [string, (path: string) => void] {
  const location = useLocation();
  const selectedPath =
    location.pathname === "/"
      ? ""
      : decodeURIComponent(location.pathname.substring(1));

  const navigate = useNavigate();
  const setSelectedPath = useCallback(
    (path: string) => {
      navigate(path ? `/${path}` : "/");
    },
    [navigate],
  );

  return [selectedPath, setSelectedPath];
}
