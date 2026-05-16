import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Legend,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { downloadAnalyticsCsv, fetchAnalytics, type AnalyticsReport } from "../api/client";

const COLORS = ["#3b82f6", "#22c55e", "#f59e0b", "#ef4444", "#8b5cf6"];

export default function Analytics() {
  const { t } = useTranslation();
  const [from, setFrom] = useState(() => {
    const d = new Date();
    d.setMonth(d.getMonth() - 11);
    return d.toISOString().slice(0, 10);
  });
  const [to, setTo] = useState(() => new Date().toISOString().slice(0, 10));
  const [report, setReport] = useState<AnalyticsReport | null>(null);

  const load = () => fetchAnalytics(from, to).then((r) => setReport(r.data)).catch(console.error);

  useEffect(() => {
    load();
  }, [from, to]);

  const expenses = report?.by_category.filter((c) => c.kind === "expense") ?? [];

  return (
    <>
      <h1>{t("nav.analytics")}</h1>
      <section className="toolbar">
        <label>
          {t("analytics.from")}
          <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} />
        </label>
        <label>
          {t("analytics.to")}
          <input type="date" value={to} onChange={(e) => setTo(e.target.value)} />
        </label>
        <button type="button" onClick={() => downloadAnalyticsCsv(from, to)}>
          {t("analytics.export")}
        </button>
      </section>

      <section className="grid-2">
        <article className="card">
          <p className="label">{t("analytics.by_month")}</p>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={report?.by_month ?? []}>
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(128,128,128,0.2)" />
              <XAxis dataKey="label" stroke="var(--muted)" fontSize={11} />
              <YAxis stroke="var(--muted)" fontSize={11} />
              <Tooltip />
              <Legend />
              <Bar dataKey="income" fill="#22c55e" name={t("tx.income")} />
              <Bar dataKey="expense" fill="#ef4444" name={t("tx.expense")} />
            </BarChart>
          </ResponsiveContainer>
        </article>
        <article className="card">
          <p className="label">{t("analytics.by_category")}</p>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie data={expenses} dataKey="total" nameKey="name" cx="50%" cy="50%" outerRadius={100} label>
                {expenses.map((_, i) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </article>
      </section>
    </>
  );
}
