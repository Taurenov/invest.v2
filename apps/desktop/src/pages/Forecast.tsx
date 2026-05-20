import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { Line, LineChart, ReferenceLine, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import Disclaimer from "../components/Disclaimer";
import {
  fetchForecast,
  fetchForecastHistory,
  fetchPriceHistory,
  type Forecast,
  type ForecastHistoryItem,
  type PricePoint,
} from "../api/client";

export default function ForecastPage() {
  const { t, i18n } = useTranslation();
  const [symbol, setSymbol] = useState("SBER");
  const [horizon, setHorizon] = useState(7);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<Forecast | null>(null);
  const [history, setHistory] = useState<ForecastHistoryItem[]>([]);
  const [series, setSeries] = useState<PricePoint[]>([]);

  const run = async () => {
    setLoading(true);
    try {
      const [forecastRes, historyRes, seriesRes] = await Promise.all([
        fetchForecast(symbol, horizon, i18n.language),
        fetchForecastHistory(symbol, 20),
        fetchPriceHistory(symbol, 120),
      ]);
      setResult(forecastRes.data);
      setHistory(historyRes.data);
      setSeries(seriesRes.data.points);
    } catch (e) {
      console.error(e);
      setResult(null);
      setHistory([]);
    } finally {
      setLoading(false);
    }
  };

  const chartData = useMemo(() => {
    if (!result || series.length === 0) return [];
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
      <h1>{t("forecast.title")}</h1>
      <section className="toolbar">
        <label>
          Тикер{" "}
          <select value={symbol} onChange={(e) => setSymbol(e.target.value)}>
            <option value="SBER">SBER</option>
            <option value="GAZP">GAZP</option>
            <option value="LKOH">LKOH</option>
          </select>
        </label>
        <label>
          {t("forecast.horizon")}{" "}
          <input type="number" min={1} max={30} value={horizon} onChange={(e) => setHorizon(Number(e.target.value))} />
        </label>
        <button type="button" onClick={run} disabled={loading}>
          {loading ? t("forecast.loading") : t("forecast.run")}
        </button>
      </section>

      {result && (
        <>
          <article className="card">
            <p className="label">{t("forecast.chart_title")}</p>
            <ResponsiveContainer width="100%" height={260}>
              <LineChart data={chartData}>
                <XAxis dataKey="name" />
                <YAxis domain={["auto", "auto"]} />
                <Tooltip />
                <ReferenceLine
                  y={chartData[Math.max(chartData.length - 2, 0)]?.price}
                  stroke="var(--muted)"
                  strokeDasharray="4 4"
                />
                <Line type="monotone" dataKey="price" stroke="var(--primary)" strokeWidth={2} dot />
              </LineChart>
            </ResponsiveContainer>
            <p style={{ marginTop: "1rem" }}>{result.narrative}</p>
            <p className="label" style={{ marginTop: "0.5rem" }}>
              Модель: {result.model_version} · уверенность: {(result.confidence * 100).toFixed(0)}%
            </p>
          </article>
          <article className="card" style={{ marginTop: "1rem" }}>
            <p className="label">История прогнозов</p>
            <table>
              <thead>
                <tr>
                  <th>Дата</th>
                  <th>Горизонт</th>
                  <th>Δ%</th>
                  <th>Уверенность</th>
                </tr>
              </thead>
              <tbody>
                {history.map((h) => (
                  <tr key={`${h.created_at}-${h.horizon_days}`}>
                    <td>{new Date(h.created_at).toLocaleString()}</td>
                    <td>{h.horizon_days} дн.</td>
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
          <Disclaimer />
        </>
      )}
    </>
  );
}
