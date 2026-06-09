import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useLocation } from "react-router-dom";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import Disclaimer from "../components/Disclaimer";
import { Flash } from "../components/Flash";
import RangeField from "../components/RangeField";
import {
  fetchForecast,
  fetchForecastHistory,
  fetchMarketCatalog,
  fetchPriceHistory,
  fetchSummary,
  type Forecast,
  type ForecastHistoryItem,
  type MoexInstrument,
  type PricePoint,
} from "../api/client";

const CHART_GRADIENT_ID = "forecastGradient";

export default function ForecastPage() {
  const { t, i18n } = useTranslation();
  const location = useLocation();
  const routeSymbol = (location.state as { symbol?: string } | null)?.symbol;
  const [instruments, setInstruments] = useState<MoexInstrument[]>([]);
  const [symbol, setSymbol] = useState(routeSymbol || "SBER");
  const [horizon, setHorizon] = useState(7);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState<Forecast | null>(null);
  const [history, setHistory] = useState<ForecastHistoryItem[]>([]);
  const [series, setSeries] = useState<PricePoint[]>([]);
  const [about, setAbout] = useState("");

  const selected = instruments.find((i) => i.symbol === symbol);

  useEffect(() => {
    fetchMarketCatalog("", 80)
      .then((r) => setInstruments(r.data))
      .catch(() => setInstruments([]));
  }, []);

  useEffect(() => {
    if (routeSymbol) setSymbol(routeSymbol);
  }, [routeSymbol]);

  useEffect(() => {
    if (!symbol) return;
    fetchSummary(symbol)
      .then((r) => setAbout(r.data.summary_text))
      .catch(() => setAbout(selected?.description ?? ""));
  }, [symbol, selected]);

  const run = async () => {
    setLoading(true);
    setError("");
    try {
      const forecastRes = await fetchForecast(symbol, horizon, i18n.language);
      setResult(forecastRes.data);
      try {
        const [historyRes, seriesRes] = await Promise.all([
          fetchForecastHistory(symbol, 20),
          fetchPriceHistory(symbol, 120),
        ]);
        setHistory(historyRes.data);
        setSeries(seriesRes.data.points);
      } catch {
        setHistory([]);
        setSeries([]);
      }
    } catch (e) {
      setResult(null);
      setHistory([]);
      setSeries([]);
      setError(e instanceof Error ? e.message : "Не удалось построить прогноз");
    } finally {
      setLoading(false);
    }
  };

  const chartData = useMemo(() => {
    if (!result) return [];
    const base = series.map((p) => ({
      name: new Date(p.time).toLocaleDateString(i18n.language === "ru" ? "ru-RU" : "en-US", {
        day: "2-digit",
        month: "2-digit",
      }),
      price: p.close,
    }));
    base.push({ name: `+${horizon}d`, price: result.predicted_value });
    return base;
  }, [result, series, horizon, i18n.language]);

  return (
    <>
      <div className="page-header">
        <h1>{t("forecast.title")}</h1>
      </div>

      <section className="moex-grid">
        {instruments.map((inst) => (
          <button
            key={inst.symbol}
            type="button"
            className={`moex-chip ${symbol === inst.symbol ? "active" : ""}`}
            onClick={() => setSymbol(inst.symbol)}
          >
            <strong>{inst.symbol}</strong>
            <small>{inst.name}</small>
            <span className="moex-sector">{inst.sector}</span>
          </button>
        ))}
      </section>

      {selected && (
        <article className="card company-card" style={{ marginBottom: "1rem" }}>
          <div className="company-card-head">
            <div>
              <h2 style={{ fontSize: "1.25rem", margin: 0 }}>{selected.name}</h2>
              <p className="label" style={{ marginTop: "0.35rem" }}>
                {selected.symbol} · MOEX · {selected.sector}
              </p>
            </div>
          </div>
          <p style={{ marginTop: "0.85rem", lineHeight: 1.6, color: "var(--muted)" }}>{about || selected.description}</p>
        </article>
      )}

      <article className="card">
        <section className="toolbar" style={{ alignItems: "flex-end" }}>
          <RangeField label={t("forecast.horizon")} value={horizon} onChange={setHorizon} min={1} max={30} unit={t("forecast.days")} />
          <button type="button" onClick={run} disabled={loading}>
            {loading ? t("forecast.loading") : t("forecast.run")}
          </button>
        </section>
      </article>

      {error && <Flash kind="error">{error}</Flash>}

      {result && (
        <>
          <article className="card">
            <p className="label">{t("forecast.chart_title")}</p>
            <div className="chart-wrap">
              <ResponsiveContainer width="100%" height={300}>
                <AreaChart data={chartData}>
                  <defs>
                    <linearGradient id={CHART_GRADIENT_ID} x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stopColor="#5b6ef5" stopOpacity={0.45} />
                      <stop offset="100%" stopColor="#5b6ef5" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.06)" />
                  <XAxis dataKey="name" stroke="var(--muted)" fontSize={11} tickLine={false} />
                  <YAxis stroke="var(--muted)" fontSize={11} tickLine={false} domain={["auto", "auto"]} />
                  <Tooltip contentStyle={{ background: "var(--bg-elevated)", border: "1px solid var(--border)", borderRadius: 10 }} />
                  <Area
                    type="monotone"
                    dataKey="price"
                    stroke="#5b6ef5"
                    strokeWidth={2.5}
                    fill={`url(#${CHART_GRADIENT_ID})`}
                    dot={{ r: 3, fill: "#5b6ef5" }}
                    animationDuration={800}
                  />
                  {chartData.length > 1 && (
                    <ReferenceLine
                      y={chartData[Math.max(chartData.length - 2, 0)]?.price}
                      stroke="var(--muted)"
                      strokeDasharray="4 4"
                    />
                  )}
                </AreaChart>
              </ResponsiveContainer>
            </div>
            <div className="forecast-result">
              <span className={result.predicted_change_pct >= 0 ? "up" : "down"} style={{ fontSize: "1.5rem", fontWeight: 700 }}>
                {result.predicted_change_pct >= 0 ? "+" : ""}
                {result.predicted_change_pct.toFixed(2)}%
              </span>
              <p style={{ marginTop: "0.75rem", lineHeight: 1.55 }}>{result.narrative}</p>
              {(result.company_name || result.sector) && (
                <p className="label" style={{ marginTop: "0.5rem" }}>
                  {result.company_name} · {result.sector}
                </p>
              )}
              <p className="label" style={{ marginTop: "0.35rem" }}>
                {t("forecast.model")}: {result.model_version} · {t("forecast.confidence")}: {(result.confidence * 100).toFixed(0)}%
              </p>
            </div>
          </article>
          {history.length > 0 && (
            <article className="card" style={{ marginTop: "1rem" }}>
              <p className="label">{t("forecast.history")}</p>
              <table>
                <thead>
                  <tr>
                    <th>{t("forecast.date")}</th>
                    <th>{t("forecast.horizon")}</th>
                    <th>Δ%</th>
                    <th>{t("forecast.confidence")}</th>
                  </tr>
                </thead>
                <tbody>
                  {history.map((h) => (
                    <tr key={`${h.created_at}-${h.horizon_days}`}>
                      <td>{new Date(h.created_at).toLocaleString()}</td>
                      <td>{h.horizon_days} {t("forecast.days")}</td>
                      <td className={h.predicted_change_pct >= 0 ? "up" : "down"}>
                        {h.predicted_change_pct >= 0 ? "+" : ""}
                        {h.predicted_change_pct.toFixed(2)}%
                      </td>
                      <td>{(h.confidence * 100).toFixed(0)}%</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </article>
          )}
          <Disclaimer />
        </>
      )}
    </>
  );
}
