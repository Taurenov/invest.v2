import type { Transaction } from "../api/client";

const TX_CACHE_KEY = "offline_tx_cache_v1";
const TX_QUEUE_KEY = "offline_tx_queue_v1";

type PendingOp =
  | { type: "create"; payload: Partial<Transaction> & { kind: string; amount: number } }
  | { type: "delete"; id: string };

export function getCachedTransactions(): Transaction[] {
  const raw = localStorage.getItem(TX_CACHE_KEY);
  if (!raw) return [];
  try {
    return JSON.parse(raw) as Transaction[];
  } catch {
    return [];
  }
}

export function setCachedTransactions(items: Transaction[]) {
  localStorage.setItem(TX_CACHE_KEY, JSON.stringify(items));
}

export function enqueue(op: PendingOp) {
  const q = getQueue();
  q.push(op);
  localStorage.setItem(TX_QUEUE_KEY, JSON.stringify(q));
}

export function clearQueue() {
  localStorage.removeItem(TX_QUEUE_KEY);
}

export function getQueue(): PendingOp[] {
  const raw = localStorage.getItem(TX_QUEUE_KEY);
  if (!raw) return [];
  try {
    return JSON.parse(raw) as PendingOp[];
  } catch {
    return [];
  }
}

export async function flushQueue(
  apiCreate: (body: Partial<Transaction> & { kind: string; amount: number }) => Promise<unknown>,
  apiDelete: (id: string) => Promise<unknown>
) {
  const q = getQueue();
  if (q.length === 0) return;
  const failed: PendingOp[] = [];
  for (const op of q) {
    try {
      if (op.type === "create") {
        await apiCreate(op.payload);
      } else {
        await apiDelete(op.id);
      }
    } catch {
      failed.push(op);
    }
  }
  if (failed.length === 0) clearQueue();
  else localStorage.setItem(TX_QUEUE_KEY, JSON.stringify(failed));
}
