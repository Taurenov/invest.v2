import { FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Flash } from "../components/Flash";
import NumberField from "../components/NumberField";
import { contributeGoal, createGoal, fetchGoals, formatMoney, type Goal } from "../api/client";

export default function GoalsPage() {
  const { t } = useTranslation();
  const [goals, setGoals] = useState<Goal[]>([]);
  const [title, setTitle] = useState("");
  const [target, setTarget] = useState(100000);
  const [contribId, setContribId] = useState("");
  const [contribAmt, setContribAmt] = useState(5000);
  const [error, setError] = useState("");
  const [ok, setOk] = useState("");

  const load = () =>
    fetchGoals()
      .then((r) => setGoals(r.data))
      .catch((e) => setError(e instanceof Error ? e.message : "Ошибка загрузки"));

  useEffect(() => {
    load();
  }, []);

  const onCreate = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setOk("");
    if (!title.trim() || target <= 0) {
      setError("Укажите название и сумму цели");
      return;
    }
    try {
      await createGoal({ title: title.trim(), goal_type: "savings", target_amount: target, currency: "RUB" });
      setTitle("");
      setTarget(100000);
      setOk("Цель создана");
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось создать цель");
    }
  };

  const onContrib = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setOk("");
    if (!contribId || contribAmt <= 0) return;
    try {
      await contributeGoal(contribId, contribAmt);
      setContribAmt(5000);
      setOk("Взнос добавлен");
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось внести вклад");
    }
  };

  return (
    <>
      <div className="page-header">
        <h1>{t("goals.title")}</h1>
      </div>
      {error && <Flash kind="error">{error}</Flash>}
      {ok && <Flash kind="success">{ok}</Flash>}

      <section className="grid-2">
        <article className="card">
          <p className="label">{t("goals.create")}</p>
          <form onSubmit={onCreate}>
            <label>
              <span className="field-label">{t("goals.name")}</span>
              <input value={title} onChange={(e) => setTitle(e.target.value)} required />
            </label>
            <NumberField label={t("goals.target")} value={target} onChange={setTarget} min={1000} step={1000} suffix="₽" />
            <button type="submit" style={{ marginTop: "1rem" }}>
              {t("goals.create_btn")}
            </button>
          </form>
        </article>
        <article className="card">
          <p className="label">{t("goals.contribute")}</p>
          <form onSubmit={onContrib}>
            <label>
              <span className="field-label">{t("goals.pick")}</span>
              <select value={contribId} onChange={(e) => setContribId(e.target.value)} required>
                <option value="">—</option>
                {goals.map((g) => (
                  <option key={g.id} value={g.id}>
                    {g.title}
                  </option>
                ))}
              </select>
            </label>
            <NumberField label={t("goals.amount")} value={contribAmt} onChange={setContribAmt} min={100} step={500} suffix="₽" />
            <button type="submit" style={{ marginTop: "1rem" }}>
              {t("goals.contribute_btn")}
            </button>
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
                <span style={{ width: `${pct}%` }} />
              </p>
            </article>
          );
        })}
      </section>
    </>
  );
}
