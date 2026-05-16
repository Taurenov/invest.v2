import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Line, LineChart, ReferenceLine, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import Disclaimer from "../components/Disclaimer";
import { fetchForecast, type Forecast } from "../api/client";

export default function ForecastPage() {
  const { t, i18n } = useTranslation();
  const [symbol, setSymbol] = useState("SBER");
  const [horizon, setHorizon] = useState(7);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<Forecast | null>(null);

  const run = async () => {
    setLoading(true);
    try {
      const r = await fetchForecast(symbol, horizon, i18n.language);
      setResult(r.data);
    } catch (e) {
      console.error(e);
      setResult(null);
    } finally {
      setLoading(false);
    }
  };

  const chartData = result
    ? [
        { name: "now", price: result.predicted_value / (1 + result.predicted_change_pct / 100) },
        { name: `+${horizon}d`, price: result.predicted_value },
      ]
    : [];

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
                <ReferenceLine y={chartData[0]?.price} stroke="var(--muted)" strokeDasharray="4 4" />
                <Line type="monotone" dataKey="price" stroke="var(--primary)" strokeWidth={2} dot />
              </LineChart>
            </ResponsiveContainer>
            <p style={{ marginTop: "1rem" }}>{result.narrative}</p>
            <p className="label" style={{ marginTop: "0.5rem" }}>
              Модель: {result.model_version} · уверенность: {(result.confidence * 100).toFixed(0)}%
            </p>
          </article>
          <Disclaimer />
        </>
      )}
    </>
  );
}
