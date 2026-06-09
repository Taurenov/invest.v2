import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import TransactionForm from "../components/TransactionForm";
import { fetchGoals, fetchTransactions, formatMoney, type Goal, type Transaction } from "../api/client";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

export default function Home() {
  const { t } = useTranslation();
  const [tx, setTx] = useState<Transaction[]>([]);
  const [goals, setGoals] = useState<Goal[]>([]);
  const [showForm, setShowForm] = useState(false);

  const reload = () => {
    fetchTransactions().then((r) => setTx(r.data)).catch(console.error);
    fetchGoals().then((r) => setGoals(r.data)).catch(console.error);
  };

  useEffect(() => {
    reload();
  }, []);

  const income = tx.filter((x) => x.kind === "income").reduce((s, x) => s + x.amount, 0);
  const expense = tx.filter((x) => x.kind === "expense").reduce((s, x) => s + x.amount, 0);
  const balance = income - expense;

  const chartData = tx
    .slice()
    .sort((a, b) => new Date(a.occurred_at).getTime() - new Date(b.occurred_at).getTime())
    .slice(-12)
    .map((x) => ({
      name: new Date(x.occurred_at).toLocaleDateString("ru-RU", { day: "2-digit", month: "short" }),
      v: x.kind === "income" ? x.amount : -x.amount,
    }));

  return (
    <>
      <div className="page-header">
        <h1>{t("home.greeting")}</h1>
        <button type="button" onClick={() => setShowForm(true)}>
          + {t("tx.add")}
        </button>
      </div>
      {showForm && <TransactionForm onSaved={reload} onClose={() => setShowForm(false)} />}

      <section className="cards">
        <article className="card card-hero" style={{ gridColumn: "span 1" }}>
          <p className="label">{t("home.balance")}</p>
          <p className="value">{formatMoney(balance)}</p>
        </article>
        <article className="card">
          <p className="label">{t("home.income_month")}</p>
          <p className="value up">+{formatMoney(income)}</p>
        </article>
        <article className="card">
          <p className="label">{t("home.expense_month")}</p>
          <p className="value down">−{formatMoney(expense)}</p>
        </article>
      </section>

      <section className="grid-2">
        <article className="card">
          <p className="label">{t("home.recent")}</p>
          <div className="chart-wrap">
            <ResponsiveContainer width="100%" height={240}>
              <AreaChart data={chartData}>
                <defs>
                  <linearGradient id="homeFlow" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#5b6ef5" stopOpacity={0.5} />
                    <stop offset="100%" stopColor="#5b6ef5" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.06)" />
                <XAxis dataKey="name" stroke="var(--muted)" fontSize={11} tickLine={false} />
                <YAxis stroke="var(--muted)" fontSize={11} tickLine={false} />
                <Tooltip />
                <Area
                  type="monotone"
                  dataKey="v"
                  stroke="#5b6ef5"
                  strokeWidth={2}
                  fill="url(#homeFlow)"
                  animationDuration={700}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </article>
        <article className="card">
          <p className="label">{t("home.goals")}</p>
          {goals.map((g) => {
            const pct = Math.min(100, (g.current_amount / g.target_amount) * 100);
            return (
              <section key={g.id} style={{ marginTop: "1.1rem" }}>
                <p style={{ fontWeight: 600 }}>{g.title}</p>
                <p className="value" style={{ fontSize: "1rem", marginTop: "0.25rem" }}>
                  {pct.toFixed(0)}% · {formatMoney(g.current_amount)} / {formatMoney(g.target_amount)}
                </p>
                <p className="progress">
                  <span style={{ width: `${pct}%` }} />
                </p>
              </section>
            );
          })}
          {goals.length === 0 && <p className="label" style={{ marginTop: "1rem" }}>—</p>}
        </article>
      </section>

      <article className="card" style={{ marginTop: "1.15rem" }}>
        <p className="label">{t("home.recent")}</p>
        <table>
          <tbody>
            {tx.slice(0, 6).map((x) => (
              <tr key={x.id}>
                <td>{new Date(x.occurred_at).toLocaleDateString()}</td>
                <td>{x.description || "—"}</td>
                <td className={x.kind === "income" ? "up" : "down"}>
                  {x.kind === "income" ? "+" : "−"}
                  {formatMoney(x.amount)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </article>
    </>
  );
}
