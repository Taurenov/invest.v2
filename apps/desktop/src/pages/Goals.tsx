import { FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { contributeGoal, createGoal, fetchGoals, formatMoney, type Goal } from "../api/client";

export default function GoalsPage() {
  const { t } = useTranslation();
  const [goals, setGoals] = useState<Goal[]>([]);
  const [title, setTitle] = useState("");
  const [target, setTarget] = useState("");
  const [contribId, setContribId] = useState("");
  const [contribAmt, setContribAmt] = useState("");

  const load = () => fetchGoals().then((r) => setGoals(r.data)).catch(console.error);
  useEffect(() => {
    load();
  }, []);

  const onCreate = async (e: FormEvent) => {
    e.preventDefault();
    const val = parseFloat(target);
    if (!title || !val) return;
    await createGoal({ title, goal_type: "savings", target_amount: val, currency: "RUB" });
    setTitle("");
    setTarget("");
    load();
  };

  const onContrib = async (e: FormEvent) => {
    e.preventDefault();
    const val = parseFloat(contribAmt);
    if (!contribId || !val) return;
    await contributeGoal(contribId, val);
    setContribAmt("");
    load();
  };

  return (
    <>
      <h1>{t("goals.title")}</h1>
      <section className="grid-2">
        <article className="card">
          <p className="label">{t("goals.create")}</p>
          <form onSubmit={onCreate}>
            <label>
              {t("goals.name")}
              <input value={title} onChange={(e) => setTitle(e.target.value)} required />
            </label>
            <label>
              {t("goals.target")}
              <input type="number" value={target} onChange={(e) => setTarget(e.target.value)} required />
            </label>
            <button type="submit">{t("goals.create_btn")}</button>
          </form>
        </article>
        <article className="card">
          <p className="label">{t("goals.contribute")}</p>
          <form onSubmit={onContrib}>
            <label>
              {t("goals.pick")}
              <select value={contribId} onChange={(e) => setContribId(e.target.value)} required>
                <option value="">—</option>
                {goals.map((g) => (
                  <option key={g.id} value={g.id}>
                    {g.title}
                  </option>
                ))}
              </select>
            </label>
            <label>
              {t("goals.amount")}
              <input type="number" value={contribAmt} onChange={(e) => setContribAmt(e.target.value)} required />
            </label>
            <button type="submit">{t("goals.contribute_btn")}</button>
          </form>
        </article>
      </section>

      <section className="cards" style={{ marginTop: "1rem" }}>
        {goals.map((g) => {
          const pct = Math.min(100, (g.current_amount / g.target_amount) * 100);
          return (
            <article className="card" key={g.id}>
              <p className="label">{g.title}</p>
              <p className="value" style={{ fontSize: "1.1rem" }}>
                {pct.toFixed(0)}% · {formatMoney(g.current_amount)} / {formatMoney(g.target_amount)}
              </p>
              <p className="progress">
                <span style={{ display: "block", width: `${pct}%`, height: 8, background: "var(--primary)", borderRadius: 4 }} />
              </p>
            </article>
          );
        })}
      </section>
    </>
  );
}
