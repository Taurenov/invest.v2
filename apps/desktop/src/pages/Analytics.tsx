import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Legend,
  Line,
  ComposedChart,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { downloadAnalyticsCsv, fetchAnalytics, formatMoney, type AnalyticsReport } from "../api/client";

const COLORS = ["#5b6ef5", "#00d4aa", "#8b5cf6", "#ff6b6b", "#fbbf24", "#38bdf8", "#f472b6", "#a78bfa"];

const tooltipStyle = {
  background: "var(--bg-elevated)",
  border: "1px solid var(--border)",
  borderRadius: 10,
  fontSize: 13,
};

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

  const stats = useMemo(() => {
    const months = report?.by_month ?? [];
    const income = months.reduce((s, m) => s + m.income, 0);
    const expense = months.reduce((s, m) => s + m.expense, 0);
    const net = income - expense;
    const rate = income > 0 ? (net / income) * 100 : 0;
    const avgExpense = months.length ? expense / months.length : 0;
    return { income, expense, net, rate, avgExpense, months: months.length };
  }, [report]);

  const expenses = report?.by_category.filter((c) => c.kind === "expense") ?? [];
  const incomes = report?.by_category.filter((c) => c.kind === "income") ?? [];
  const expenseTotal = expenses.reduce((s, c) => s + c.total, 0);

  const trendData = useMemo(() => {
    let running = 0;
    return (report?.by_month ?? []).map((m) => {
      running += m.income - m.expense;
      return { ...m, net: m.income - m.expense, balance: running };
    });
  }, [report]);

  const pieData = expenses.map((c) => ({
    name: c.name,
    value: c.total,
    pct: expenseTotal > 0 ? (c.total / expenseTotal) * 100 : 0,
  }));

  return (
    <>
      <div className="page-header">
        <h1>{t("nav.analytics")}</h1>
        <button type="button" onClick={() => downloadAnalyticsCsv(from, to)}>
          {t("analytics.export")}
        </button>
      </div>

      <section className="toolbar analytics-filters">
        <label className="date-chip">
          {t("analytics.from")}
          <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} />
        </label>
        <label className="date-chip">
          {t("analytics.to")}
          <input type="date" value={to} onChange={(e) => setTo(e.target.value)} />
        </label>
      </section>

      <section className="cards">
        <article className="card">
          <p className="label">{t("analytics.total_income")}</p>
          <p className="value up">+{formatMoney(stats.income)}</p>
        </article>
        <article className="card">
          <p className="label">{t("analytics.total_expense")}</p>
          <p className="value down">−{formatMoney(stats.expense)}</p>
        </article>
        <article className="card card-hero">
          <p className="label">{t("analytics.net")}</p>
          <p className="value">{formatMoney(stats.net)}</p>
        </article>
        <article className="card">
          <p className="label">{t("analytics.savings_rate")}</p>
          <p className={`value ${stats.rate >= 0 ? "up" : "down"}`}>{stats.rate.toFixed(1)}%</p>
        </article>
      </section>

      <section className="grid-2">
        <article className="card">
          <p className="label">{t("analytics.by_month")}</p>
          <div className="chart-wrap">
            <ResponsiveContainer width="100%" height={320}>
              <BarChart data={report?.by_month ?? []} barGap={4} barCategoryGap="18%">
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.06)" vertical={false} />
                <XAxis dataKey="label" stroke="var(--muted)" fontSize={11} tickLine={false} axisLine={false} />
                <YAxis stroke="var(--muted)" fontSize={11} tickLine={false} axisLine={false} tickFormatter={(v) => `${(v / 1000).toFixed(0)}k`} />
                <Tooltip contentStyle={tooltipStyle} formatter={(v: number) => formatMoney(v)} />
                <Legend />
                <Bar dataKey="income" name={t("tx.income")} fill="#00d4aa" radius={[6, 6, 0, 0]} animationDuration={700} />
                <Bar dataKey="expense" name={t("tx.expense")} fill="#ff6b6b" radius={[6, 6, 0, 0]} animationDuration={700} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </article>

        <article className="card">
          <p className="label">{t("analytics.by_category")}</p>
          <div className="chart-wrap">
            <ResponsiveContainer width="100%" height={320}>
              <PieChart>
                <Pie
                  data={pieData}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  innerRadius={70}
                  outerRadius={110}
                  paddingAngle={3}
                  animationDuration={800}
                  label={({ name, pct }) => (pct > 8 ? `${name} ${pct.toFixed(0)}%` : "")}
                >
                  {pieData.map((_, i) => (
                    <Cell key={i} fill={COLORS[i % COLORS.length]} stroke="transparent" />
                  ))}
                </Pie>
                <Tooltip contentStyle={tooltipStyle} formatter={(v: number) => formatMoney(v)} />
              </PieChart>
            </ResponsiveContainer>
          </div>
          <ul className="legend-list">
            {pieData.map((c, i) => (
              <li key={c.name}>
                <span className="legend-dot" style={{ background: COLORS[i % COLORS.length] }} />
                <span>{c.name}</span>
                <strong>{formatMoney(c.value)}</strong>
                <span className="muted-pct">{c.pct.toFixed(1)}%</span>
              </li>
            ))}
          </ul>
        </article>
      </section>

      <article className="card" style={{ marginTop: "1.15rem" }}>
        <p className="label">{t("analytics.net_trend")}</p>
        <div className="chart-wrap">
          <ResponsiveContainer width="100%" height={240}>
            <ComposedChart data={trendData}>
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.06)" vertical={false} />
              <XAxis dataKey="label" stroke="var(--muted)" fontSize={11} tickLine={false} axisLine={false} />
              <YAxis stroke="var(--muted)" fontSize={11} tickLine={false} axisLine={false} />
              <Tooltip contentStyle={tooltipStyle} />
              <Bar dataKey="net" name={t("analytics.net")} fill="#5b6ef5" radius={[4, 4, 0, 0]} opacity={0.85} />
              <Line type="monotone" dataKey="balance" name={t("analytics.balance_trend")} stroke="#00d4aa" strokeWidth={2.5} dot={false} />
            </ComposedChart>
          </ResponsiveContainer>
        </div>
      </article>

      <section className="grid-2" style={{ marginTop: "1.15rem" }}>
        <article className="card">
          <p className="label">{t("analytics.expense_detail")}</p>
          <CategoryTable items={expenses} total={expenseTotal} />
        </article>
        <article className="card">
          <p className="label">{t("analytics.income_detail")}</p>
          <CategoryTable
            items={incomes}
            total={incomes.reduce((s, c) => s + c.total, 0)}
            accent="#00d4aa"
          />
        </article>
      </section>

      <p className="label" style={{ marginTop: "1rem", textAlign: "center" }}>
        {t("analytics.avg_monthly")}: {formatMoney(stats.avgExpense)} · {stats.months} {t("analytics.months")}
      </p>
    </>
  );
}

function CategoryTable({
  items,
  total,
  accent = "#5b6ef5",
}: {
  items: { name: string; total: number }[];
  total: number;
  accent?: string;
}) {
  if (!items.length) return <p className="label" style={{ marginTop: "1rem" }}>—</p>;
  return (
    <table>
      <thead>
        <tr>
          <th>Категория</th>
          <th>Сумма</th>
          <th>%</th>
        </tr>
      </thead>
      <tbody>
        {items
          .sort((a, b) => b.total - a.total)
          .map((c) => {
            const pct = total > 0 ? (c.total / total) * 100 : 0;
            return (
              <tr key={c.name}>
                <td>
                  {c.name}
                  <div className="budget-bar">
                    <span className="ok" style={{ width: `${pct}%`, background: accent }} />
                  </div>
                </td>
                <td>{formatMoney(c.total)}</td>
                <td>{pct.toFixed(1)}%</td>
              </tr>
            );
          })}
      </tbody>
    </table>
  );
}
