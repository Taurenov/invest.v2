const API_BASE = import.meta.env.VITE_API_URL ?? "http://127.0.0.1:8080";
const DEV_TOKEN = import.meta.env.VITE_API_TOKEN ?? "dev-token";

const TOKEN_KEY = "fin_token";
const USER_KEY = "fin_user";

export type User = { id: string; email: string; display_name: string };

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function getUser(): User | null {
  const raw = localStorage.getItem(USER_KEY);
  return raw ? (JSON.parse(raw) as User) : null;
}

export function setSession(token: string, user: User) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}

function authHeader(): Record<string, string> {
  const t = getToken();
  const token = t || DEV_TOKEN;
  return { Authorization: `Bearer ${token}`, "Content-Type": "application/json" };
}

async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: { ...authHeader(), ...init?.headers },
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `API ${res.status}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export type Transaction = {
  id: string;
  kind: "income" | "expense";
  amount: number;
  currency: string;
  description: string;
  occurred_at: string;
  category_id?: string;
};

export type Category = {
  id: string;
  name: string;
  kind: "income" | "expense";
  icon?: string;
  color?: string;
};

export type Goal = {
  id: string;
  title: string;
  goal_type: string;
  target_amount: number;
  current_amount: number;
  currency: string;
};

export type Quote = {
  symbol: string;
  exchange: string;
  price: number;
  change_pct: number;
};

export type Forecast = {
  symbol: string;
  predicted_change_pct: number;
  predicted_value: number;
  confidence: number;
  narrative: string;
  disclaimer: string;
  model_version: string;
};

export type ForecastHistoryItem = {
  created_at: string;
  horizon_days: number;
  predicted_change_pct: number;
  confidence: number;
  model_version: string;
};

export type PricePoint = { time: string; close: number };

export async function login(email: string, password: string) {
  const r = await api<{ token: string; user: User }>("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
  setSession(r.token, r.user);
  return r;
}

export async function register(email: string, password: string, display_name: string) {
  const r = await api<{ token: string; user: User }>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify({ email, password, display_name }),
  });
  setSession(r.token, r.user);
  return r;
}

export function fetchMe() {
  return api<{ data: User }>("/api/v1/me");
}

export function fetchTransactions() {
  return api<{ data: Transaction[] }>("/api/v1/me/transactions");
}

export function searchTransactions(params: { q?: string; kind?: string; from?: string; to?: string }) {
  const qs = new URLSearchParams();
  if (params.q) qs.set("q", params.q);
  if (params.kind) qs.set("kind", params.kind);
  if (params.from) qs.set("from", params.from);
  if (params.to) qs.set("to", params.to);
  return api<{ data: Transaction[] }>(`/api/v1/me/transactions?${qs.toString()}`);
}

export function createTransaction(body: Partial<Transaction> & { kind: string; amount: number }) {
  return api<{ data: Transaction }>("/api/v1/me/transactions", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export function deleteTransaction(id: string) {
  return api<void>(`/api/v1/me/transactions/${id}`, { method: "DELETE" });
}

export function fetchCategories() {
  return api<{ data: Category[] }>("/api/v1/me/categories");
}

export function createCategory(body: { name: string; kind: string; icon?: string; color?: string }) {
  return api<{ data: Category }>("/api/v1/me/categories", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export function fetchGoals() {
  return api<{ data: Goal[] }>("/api/v1/me/goals");
}

export function createGoal(body: {
  title: string;
  goal_type: string;
  target_amount: number;
  currency?: string;
}) {
  return api<{ data: Goal }>("/api/v1/me/goals", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export function contributeGoal(goalId: string, amount: number, note?: string) {
  return api<{ data: Goal }>(`/api/v1/me/goals/${goalId}/contribute`, {
    method: "POST",
    body: JSON.stringify({ amount, note }),
  });
}

export function calcROI(initial: number, current: number) {
  return api<{ data: { roi_percent: number } }>("/api/v1/calculator/roi", {
    method: "POST",
    body: JSON.stringify({ initial, current }),
  });
}

export function calcCAGR(initial: number, final: number, years: number) {
  return api<{ data: { cagr_percent: number } }>("/api/v1/calculator/cagr", {
    method: "POST",
    body: JSON.stringify({ initial, final, years }),
  });
}

export function calcSavings(monthly: number, annual_rate_pct: number, months: number, initial_balance = 0) {
  return api<{ data: { future_value: number } }>("/api/v1/calculator/savings", {
    method: "POST",
    body: JSON.stringify({ monthly, annual_rate_pct, months, initial_balance }),
  });
}

export function fetchQuote(symbol: string) {
  return api<{ data: Quote }>(`/api/v1/markets/${symbol}/quote?exchange=MOEX`);
}

export function fetchForecast(symbol: string, horizonDays: number, locale: string) {
  return api<{ data: Forecast }>(
    `/api/v1/markets/${symbol}/forecast?exchange=MOEX&horizon_days=${horizonDays}&locale=${locale}`
  );
}

export function fetchForecastHistory(symbol: string, limit = 20) {
  return api<{ data: ForecastHistoryItem[] }>(
    `/api/v1/markets/${symbol}/forecast/history?exchange=MOEX&limit=${limit}`
  );
}

export function fetchPriceHistory(symbol: string, points = 120) {
  return api<{ data: { symbol: string; exchange: string; points: PricePoint[] } }>(
    `/api/v1/markets/${symbol}/history?exchange=MOEX&points=${points}`
  );
}

export type Holding = {
  id: string;
  symbol: string;
  exchange: string;
  name: string;
  quantity: number;
  avg_cost: number;
  current_price?: number;
  market_value?: number;
  pnl?: number;
  pnl_percent?: number;
};

export type Portfolio = {
  id: string;
  name: string;
  holdings: Holding[];
  total_cost: number;
  total_value: number;
  total_pnl: number;
};

export type WatchlistItem = { instrument_id: string; symbol: string; exchange: string; name: string };

export type UserSettings = {
  locale: string;
  base_currency: string;
  theme: "dark" | "light" | "system";
  timezone: string;
};

export type Tag = { id: string; name: string; color?: string; created_at: string };
export type Budget = { id: string; category_id: string; amount: number; currency: string; period: string };
export type BudgetStatus = { budget: Budget; spent: number; remaining: number; percent: number };
export type Recurring = {
  id: string;
  kind: string;
  category_id?: string;
  amount: number;
  currency: string;
  description?: string;
  schedule: "daily" | "weekly" | "monthly";
  day_of_month?: number;
  day_of_week?: number;
  next_run_at: string;
  is_active: boolean;
};

export function fetchTags() {
  return api<{ data: Tag[] }>("/api/v1/me/tags");
}
export function createTag(name: string, color = "#3b82f6") {
  return api<{ data: Tag }>("/api/v1/me/tags", { method: "POST", body: JSON.stringify({ name, color }) });
}
export function deleteTag(id: string) {
  return api<void>(`/api/v1/me/tags/${id}`, { method: "DELETE" });
}
export function setTransactionTags(txID: string, tagIDs: string[]) {
  return api<void>(`/api/v1/me/transactions/${txID}/tags`, { method: "PUT", body: JSON.stringify({ tag_ids: tagIDs }) });
}

export function fetchBudgets() {
  return api<{ data: BudgetStatus[] }>("/api/v1/me/budgets");
}
export function upsertBudget(category_id: string, amount: number, currency = "RUB") {
  return api<{ data: Budget }>("/api/v1/me/budgets", { method: "POST", body: JSON.stringify({ category_id, amount, currency }) });
}
export function deleteBudget(id: string) {
  return api<void>(`/api/v1/me/budgets/${id}`, { method: "DELETE" });
}

export function fetchRecurring() {
  return api<{ data: Recurring[] }>("/api/v1/me/recurring");
}
export function createRecurring(body: Partial<Recurring> & { kind: string; amount: number; schedule: string }) {
  return api<{ data: Recurring }>("/api/v1/me/recurring", { method: "POST", body: JSON.stringify(body) });
}
export function toggleRecurring(id: string, active: boolean) {
  return api<void>(`/api/v1/me/recurring/${id}/toggle?active=${active}`, { method: "POST" });
}
export function deleteRecurring(id: string) {
  return api<void>(`/api/v1/me/recurring/${id}`, { method: "DELETE" });
}

export type Alert = { id: string; type: string; title: string; message: string; read: boolean };

export type AnalyticsReport = {
  from: string;
  to: string;
  by_month: { label: string; income: number; expense: number }[];
  by_category: { name: string; kind: string; total: number }[];
};

export type CompanySummary = {
  symbol: string;
  exchange: string;
  summary_text: string;
  key_metrics: Record<string, unknown>;
};

export function fetchPortfolio() {
  return api<{ data: Portfolio }>("/api/v1/me/portfolio");
}

export function addHolding(body: { symbol: string; exchange?: string; quantity: number; avg_cost: number }) {
  return api<{ data: Holding }>("/api/v1/me/portfolio/holdings", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export function removeHolding(id: string) {
  return api<void>(`/api/v1/me/portfolio/holdings/${id}`, { method: "DELETE" });
}

export function fetchWatchlist() {
  return api<{ data: WatchlistItem[] }>("/api/v1/me/watchlist");
}

export function addWatchlistItem(symbol: string, exchange = "MOEX") {
  return api<{ data: WatchlistItem }>("/api/v1/me/watchlist/items", {
    method: "POST",
    body: JSON.stringify({ symbol, exchange }),
  });
}

export function removeWatchlistItem(instrumentId: string) {
  return api<void>(`/api/v1/me/watchlist/items/${instrumentId}`, { method: "DELETE" });
}

export function fetchSettings() {
  return api<{ data: UserSettings }>("/api/v1/me/settings");
}

export function updateSettings(s: Partial<UserSettings>) {
  return api<{ data: UserSettings }>("/api/v1/me/settings", {
    method: "PATCH",
    body: JSON.stringify(s),
  });
}

export function fetchAnalytics(from: string, to: string) {
  return api<{ data: AnalyticsReport }>(`/api/v1/me/analytics?from=${from}&to=${to}`);
}

export async function downloadAnalyticsCsv(from: string, to: string) {
  const res = await fetch(
    `${API_BASE}/api/v1/me/analytics/export.csv?from=${from}&to=${to}`,
    { headers: authHeader() }
  );
  if (!res.ok) throw new Error("export failed");
  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "analytics.csv";
  a.click();
  URL.revokeObjectURL(url);
}

export function fetchSummary(symbol: string) {
  return api<{ data: CompanySummary }>(`/api/v1/markets/${symbol}/summary?exchange=MOEX`);
}

export function fetchAlerts() {
  return api<{ data: Alert[] }>("/api/v1/me/alerts");
}

export function markAlertsRead() {
  return api<void>("/api/v1/me/alerts/read", { method: "POST" });
}

export function quotesWebSocket(onMessage: (data: unknown) => void): WebSocket {
  const base = API_BASE.replace(/^http/, "ws");
  const token = getToken() || DEV_TOKEN;
  const ws = new WebSocket(`${base}/ws/quotes?token=${encodeURIComponent(token)}`);
  ws.onmessage = (ev) => {
    try {
      onMessage(JSON.parse(ev.data));
    } catch {
      /* ignore */
    }
  };
  return ws;
}

export function formatMoney(amount: number, currency = "RUB") {
  return new Intl.NumberFormat(currency === "RUB" ? "ru-RU" : "en-US", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(amount);
}

export function exportLocalBackup() {
  const payload = {
    token: localStorage.getItem(TOKEN_KEY),
    user: localStorage.getItem(USER_KEY),
    settings: localStorage.getItem("fin_settings_cache"),
    txCache: localStorage.getItem("offline_tx_cache_v1"),
    txQueue: localStorage.getItem("offline_tx_queue_v1"),
    exported_at: new Date().toISOString(),
  };
  const blob = new Blob([JSON.stringify(payload, null, 2)], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `fin-helper-backup-${new Date().toISOString().slice(0, 10)}.json`;
  a.click();
  URL.revokeObjectURL(url);
}

export async function restoreLocalBackup(file: File) {
  const text = await file.text();
  const data = JSON.parse(text) as Record<string, string | null>;
  if (typeof data.token === "string") localStorage.setItem(TOKEN_KEY, data.token);
  if (typeof data.user === "string") localStorage.setItem(USER_KEY, data.user);
  if (typeof data.settings === "string") localStorage.setItem("fin_settings_cache", data.settings);
  if (typeof data.txCache === "string") localStorage.setItem("offline_tx_cache_v1", data.txCache);
  if (typeof data.txQueue === "string") localStorage.setItem("offline_tx_queue_v1", data.txQueue);
}
