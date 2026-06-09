import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { Flash } from "../components/Flash";
import {
  addWatchlistItem,
  fetchMarketCatalog,
  fetchSummary,
  fetchWatchlist,
  removeWatchlistItem,
  type CompanySummary,
  type MoexInstrument,
  type WatchlistItem,
} from "../api/client";

export default function Markets() {
  const { t } = useTranslation();
  const [watchlist, setWatchlist] = useState<WatchlistItem[]>([]);
  const [catalog, setCatalog] = useState<MoexInstrument[]>([]);
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [summary, setSummary] = useState<CompanySummary | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);

  const loadWatchlist = () => fetchWatchlist().then((r) => setWatchlist(r.data)).catch(console.error);

  const loadCatalog = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const r = await fetchMarketCatalog(query, 300);
      setCatalog(r.data);
      setLastUpdate(new Date());
    } catch (e) {
      setError(e instanceof Error ? e.message : "Не удалось загрузить котировки MOEX");
    } finally {
      setLoading(false);
    }
  }, [query]);

  useEffect(() => {
    loadWatchlist();
  }, []);

  useEffect(() => {
    loadCatalog();
    const id = setInterval(loadCatalog, 15000);
    return () => clearInterval(id);
  }, [loadCatalog]);

  const onAddWatch = async (symbol: string) => {
    try {
      await addWatchlistItem(symbol);
      loadWatchlist();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Не удалось добавить в список");
    }
  };

  const openSummary = async (symbol: string) => {
    const r = await fetchSummary(symbol);
    setSummary(r.data);
  };

  const watchSet = new Set(watchlist.map((w) => w.symbol));

  return (
    <>
      <div className="page-header">
        <h1>{t("markets.title")}</h1>
        <span className="label">
          {loading ? t("markets.updating") : lastUpdate ? `${t("markets.updated")} ${lastUpdate.toLocaleTimeString()}` : ""}
        </span>
      </div>

      {error && <Flash kind="error">{error}</Flash>}

      {watchlist.length > 0 && (
        <article className="card" style={{ marginBottom: "1rem" }}>
          <p className="label">{t("markets.watchlist")}</p>
          <div className="watchlist-chips">
            {watchlist.map((it) => {
              const live = catalog.find((c) => c.symbol === it.symbol);
              return (
                <div key={it.instrument_id} className="watch-chip">
                  <strong>{it.symbol}</strong>
                  <span>{live?.price ? live.price.toFixed(2) : "—"}</span>
                  <span className={live && live.change_pct >= 0 ? "up" : "down"}>
                    {live ? `${live.change_pct >= 0 ? "+" : ""}${live.change_pct.toFixed(2)}%` : ""}
                  </span>
                  <button type="button" className="secondary" onClick={() => removeWatchlistItem(it.instrument_id).then(loadWatchlist)}>
                    ×
                  </button>
                </div>
              );
            })}
          </div>
        </article>
      )}

      <section className="toolbar">
        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder={t("markets.search")}
          style={{ flex: 1, minWidth: 200 }}
        />
        <button type="button" className="secondary" onClick={loadCatalog}>
          ↻
        </button>
      </section>

      <article className="card markets-table-wrap">
        <table>
          <thead>
            <tr>
              <th>{t("portfolio.ticker")}</th>
              <th>{t("markets.company")}</th>
              <th>{t("markets.sector")}</th>
              <th>{t("markets.price")}</th>
              <th>Δ%</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {catalog.map((it) => (
              <tr key={it.symbol}>
                <td>
                  <strong>{it.symbol}</strong>
                </td>
                <td style={{ color: "var(--muted)", maxWidth: 220 }}>{it.short_name || it.name}</td>
                <td>
                  <span className="moex-sector">{it.sector}</span>
                </td>
                <td>{it.price ? it.price.toFixed(2) : "—"}</td>
                <td className={it.change_pct >= 0 ? "up" : "down"}>
                  {it.price ? `${it.change_pct >= 0 ? "+" : ""}${it.change_pct.toFixed(2)}%` : "—"}
                </td>
                <td className="markets-actions">
                  <button type="button" className="secondary" onClick={() => openSummary(it.symbol)}>
                    {t("markets.summary")}
                  </button>
                  <Link to={`/forecast`} state={{ symbol: it.symbol }} className="btn secondary" style={{ textDecoration: "none" }}>
                    AI
                  </Link>
                  {!watchSet.has(it.symbol) && (
                    <button type="button" onClick={() => onAddWatch(it.symbol)}>
                      +
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        <p className="label" style={{ marginTop: "1rem", textAlign: "center" }}>
          {t("markets.count", { count: catalog.length })}
        </p>
      </article>

      {summary && (
        <article className="card" style={{ marginTop: "1rem" }}>
          <section className="toolbar">
            <h2 style={{ margin: 0 }}>
              {summary.name || summary.symbol} — {t("markets.summary")}
            </h2>
            <button type="button" className="secondary" onClick={() => setSummary(null)}>
              ×
            </button>
          </section>
          {summary.sector && (
            <p className="label" style={{ marginTop: "0.5rem" }}>
              {summary.symbol} · MOEX · {summary.sector}
            </p>
          )}
          <p style={{ marginTop: "0.75rem", lineHeight: 1.6 }}>{summary.summary_text}</p>
        </article>
      )}
    </>
  );
}
