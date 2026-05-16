import { FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  addWatchlistItem,
  fetchQuote,
  fetchSummary,
  fetchWatchlist,
  quotesWebSocket,
  removeWatchlistItem,
  type CompanySummary,
  type WatchlistItem,
} from "../api/client";

type LiveQuote = { symbol: string; price: number; change_pct: number };

export default function Markets() {
  const { t } = useTranslation();
  const [items, setItems] = useState<WatchlistItem[]>([]);
  const [quotes, setQuotes] = useState<Record<string, LiveQuote>>({});
  const [newSymbol, setNewSymbol] = useState("");
  const [summary, setSummary] = useState<CompanySummary | null>(null);

  const load = () => fetchWatchlist().then((r) => setItems(r.data)).catch(console.error);

  useEffect(() => {
    load();
  }, []);

  useEffect(() => {
    items.forEach((it) => {
      fetchQuote(it.symbol)
        .then((r) =>
          setQuotes((q) => ({
            ...q,
            [it.symbol]: { symbol: it.symbol, price: r.data.price, change_pct: r.data.change_pct },
          }))
        )
        .catch(console.error);
    });
    const ws = quotesWebSocket((msg: unknown) => {
      const m = msg as { type?: string; data?: LiveQuote[] };
      if (m.type !== "quotes" || !m.data) return;
      setQuotes((prev) => {
        const next = { ...prev };
        for (const q of m.data!) next[q.symbol] = q;
        return next;
      });
    });
    return () => ws.close();
  }, [items]);

  const onAdd = async (e: FormEvent) => {
    e.preventDefault();
    if (!newSymbol) return;
    await addWatchlistItem(newSymbol.toUpperCase());
    setNewSymbol("");
    load();
  };

  const openSummary = async (symbol: string) => {
    const r = await fetchSummary(symbol);
    setSummary(r.data);
  };

  return (
    <>
      <h1>{t("markets.title")}</h1>
      <form onSubmit={onAdd} className="toolbar">
        <input
          value={newSymbol}
          onChange={(e) => setNewSymbol(e.target.value)}
          placeholder={t("markets.add_ticker")}
        />
        <button type="submit">{t("markets.add")}</button>
      </form>

      <article className="card">
        <table>
          <thead>
            <tr>
              <th>{t("portfolio.ticker")}</th>
              <th>{t("markets.price")}</th>
              <th>Δ%</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {items.map((it) => {
              const q = quotes[it.symbol];
              return (
                <tr key={it.instrument_id}>
                  <td>{it.symbol}</td>
                  <td>{q ? q.price.toFixed(2) : "—"}</td>
                  <td className={q && q.change_pct >= 0 ? "up" : "down"}>
                    {q ? `${q.change_pct >= 0 ? "+" : ""}${q.change_pct.toFixed(2)}%` : "—"}
                  </td>
                  <td>
                    <button type="button" className="secondary" onClick={() => openSummary(it.symbol)}>
                      {t("markets.summary")}
                    </button>{" "}
                    <button
                      type="button"
                      className="secondary"
                      onClick={() => removeWatchlistItem(it.instrument_id).then(load)}
                    >
                      ×
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </article>

      {summary && (
        <article className="card" style={{ marginTop: "1rem" }}>
          <section className="toolbar">
            <h2 style={{ margin: 0 }}>
              {summary.symbol} — {t("markets.summary")}
            </h2>
            <button type="button" className="secondary" onClick={() => setSummary(null)}>
              ×
            </button>
          </section>
          <p style={{ marginTop: "0.75rem", lineHeight: 1.6 }}>{summary.summary_text}</p>
          <pre style={{ marginTop: "0.75rem", fontSize: "0.8rem", color: "var(--muted)" }}>
            {JSON.stringify(summary.key_metrics, null, 2)}
          </pre>
        </article>
      )}
    </>
  );
}
