import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import TransactionForm from "../components/TransactionForm";
import { fetchGoals, fetchTransactions, formatMoney, type Goal, type Transaction } from "../api/client";
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

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

  const chartData = tx.slice(0, 8).map((x, i) => ({
    name: new Date(x.occurred_at).toLocaleDateString("ru-RU", { day: "2-digit", month: "short" }),
    v: x.kind === "income" ? x.amount : -x.amount,
    i,
  }));

  return (
    <>
      <section className="toolbar">
        <h1 style={{ margin: 0 }}>{t("home.greeting")}</h1>
        <button type="button" onClick={() => setShowForm(true)}>
          + {t("tx.add")}
        </button>
      </section>
      {showForm && <TransactionForm onSaved={reload} onClose={() => setShowForm(false)} />}
      <section className="cards">
        <article className="card">
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
          <p className="label">{t("home.recent")} — поток</p>
          <ResponsiveContainer width="100%" height={220}>
            <LineChart data={chartData}>
              <XAxis dataKey="name" stroke="var(--muted)" fontSize={12} />
              <YAxis stroke="var(--muted)" fontSize={12} />
              <Tooltip />
              <Line type="monotone" dataKey="v" stroke="var(--primary)" dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </article>
        <article className="card">
          <p className="label">{t("home.goals")}</p>
          {goals.map((g) => {
            const pct = Math.min(100, (g.current_amount / g.target_amount) * 100);
            return (
              <section key={g.id} style={{ marginTop: "1rem" }}>
                <p>{g.title}</p>
                <p className="value" style={{ fontSize: "1rem" }}>
                  {pct.toFixed(0)}% · {formatMoney(g.current_amount)} / {formatMoney(g.target_amount)}
                </p>
                <p className="progress" style={{ display: "block" }}>
                  <span style={{ display: "block", width: `${pct}%`, height: 8, background: "var(--primary)", borderRadius: 4 }} />
                </p>
              </section>
            );
          })}
        </article>
      </section>

      <article className="card" style={{ marginTop: "1rem" }}>
        <p className="label">{t("home.recent")}</p>
        <table>
          <thead>
            <tr>
              <th>Описание</th>
              <th>Дата</th>
              <th>Сумма</th>
            </tr>
          </thead>
          <tbody>
            {tx.map((row) => (
              <tr key={row.id}>
                <td>{row.description}</td>
                <td>{new Date(row.occurred_at).toLocaleDateString()}</td>
                <td className={row.kind === "income" ? "up" : "down"}>
                  {row.kind === "income" ? "+" : "−"}
                  {formatMoney(row.amount)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </article>
    </>
  );
}
