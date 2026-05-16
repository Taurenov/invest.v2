import { useEffect, useRef } from "react";
import { fetchAlerts, markAlertsRead } from "../api/client";
import { showNotification } from "../utils/notify";

export function useAlertsPoll(intervalMs = 15000) {
  const seen = useRef<Set<string>>(new Set());

  useEffect(() => {
    const tick = async () => {
      try {
        const r = await fetchAlerts();
        for (const a of r.data) {
          if (a.read || seen.current.has(a.id)) continue;
          seen.current.add(a.id);
          await showNotification(a.title, a.message);
        }
        if (r.data.some((a) => !a.read)) {
          await markAlertsRead();
        }
      } catch {
        /* api down */
      }
    };
    tick();
    const id = setInterval(tick, intervalMs);
    return () => clearInterval(id);
  }, [intervalMs]);
}
