import type { ReactNode } from "react";

export function Flash({ kind, children }: { kind: "error" | "success" | "info"; children: ReactNode }) {
  return <div className={`flash flash-${kind}`}>{children}</div>;
}
