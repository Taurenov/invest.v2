import { FormEvent, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { Flash } from "../components/Flash";
import {
  createRecurring,
  createTag,
  deleteBudget,
  deleteRecurring,
  deleteTag,
  fetchBudgets,
  fetchCategories,
  fetchRecurring,
  fetchTags,
  toggleRecurring,
  upsertBudget,
  type BudgetStatus,
  type Category,
  type Recurring,
  type Tag,
} from "../api/client";

export default function ExtrasPage() {
  const { t } = useTranslation();
  const [tags, setTags] = useState<Tag[]>([]);
  const [budgets, setBudgets] = useState<BudgetStatus[]>([]);
  const [rec, setRec] = useState<Recurring[]>([]);
  const [cats, setCats] = useState<Category[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      const [tgs, bgs, rcs, cs] = await Promise.all([
        fetchTags(),
        fetchBudgets(),
        fetchRecurring(),
        fetchCategories(),
      ]);
      setTags(tgs.data);
      setBudgets(bgs.data);
      setRec(rcs.data);
      setCats(cs.data.filter((c) => c.kind === "expense"));
    } catch (e) {
      setError(e instanceof Error ? e.message : "Не удалось загрузить данные");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const categoryMap = useMemo(() => new Map(cats.map((c) => [c.id, `${c.icon ?? ""} ${c.name}`])), [cats]);

  return (
    <>
      <div className="page-header">
        <h1>{t("nav.extras")}</h1>
        <button type="button" className="secondary" onClick={load} disabled={loading}>
          {loading ? "…" : "↻"}
        </button>
      </div>
      {error && <Flash kind="error">{error}</Flash>}

      <section className="grid-2">
        <article className="card">
          <p className="label">{t("extras.tags")}</p>
          <TagForm onCreate={async (name) => { await createTag(name); await load(); }} />
          <table>
            <thead>
              <tr>
                <th>{t("extras.name")}</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {tags.map((x) => (
                <tr key={x.id}>
                  <td>{x.name}</td>
                  <td>
                    <button type="button" className="secondary" onClick={() => deleteTag(x.id).then(load)}>
                      ×
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </article>

        <article className="card">
          <p className="label">{t("extras.budgets")}</p>
          <BudgetForm
            categories={cats}
            onSave={async (categoryId, amount) => {
              await upsertBudget(categoryId, amount);
              await load();
            }}
          />
          <table>
            <thead>
              <tr>
                <th>{t("tx.category")}</th>
                <th>{t("tx.amount")}</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {budgets.map((b) => {
                const pct = Math.min(100, b.percent);
                const barClass = b.percent >= 100 ? "over" : b.percent >= 80 ? "warn" : "ok";
                return (
                  <tr key={b.budget.id}>
                    <td>
                      {categoryMap.get(b.budget.category_id) ?? b.budget.category_id}
                      <div className="budget-bar">
                        <span className={barClass} style={{ width: `${pct}%` }} />
                      </div>
                    </td>
                    <td>
                      <span className={b.percent >= 100 ? "down" : ""}>
                        {b.spent.toFixed(0)} / {b.budget.amount.toFixed(0)} ({b.percent.toFixed(0)}%)
                      </span>
                    </td>
                    <td>
                      <button type="button" className="secondary" onClick={() => deleteBudget(b.budget.id).then(load)}>
                        ×
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </article>
      </section>

      <article className="card" style={{ marginTop: "1rem" }}>
        <p className="label">{t("extras.recurring")}</p>
        <RecurringForm
          categories={cats}
          onCreate={async (body) => {
            await createRecurring(body);
            await load();
          }}
        />
        <table>
          <thead>
            <tr>
              <th>{t("tx.kind")}</th>
              <th>{t("tx.description")}</th>
              <th>{t("tx.amount")}</th>
              <th>{t("extras.next_run")}</th>
              <th>{t("extras.active")}</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {rec.map((r) => (
              <tr key={r.id}>
                <td>{r.kind}</td>
                <td>{r.description}</td>
                <td>{r.amount}</td>
                <td>{r.next_run_at}</td>
                <td>
                  <input
                    type="checkbox"
                    checked={r.is_active}
                    onChange={(e) => toggleRecurring(r.id, e.target.checked).then(load)}
                  />
                </td>
                <td>
                  <button type="button" className="secondary" onClick={() => deleteRecurring(r.id).then(load)}>
                    ×
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </article>
    </>
  );
}

function TagForm({ onCreate }: { onCreate: (name: string) => Promise<void> }) {
  const { t } = useTranslation();
  const [name, setName] = useState("");
  const [err, setErr] = useState("");
  const submit = async (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setErr("");
    try {
      await onCreate(name.trim());
      setName("");
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Ошибка");
    }
  };
  return (
    <>
      {err && <Flash kind="error">{err}</Flash>}
      <form onSubmit={submit} className="toolbar">
        <input value={name} onChange={(e) => setName(e.target.value)} placeholder={t("extras.name")} />
        <button type="submit">{t("extras.add")}</button>
      </form>
    </>
  );
}

function BudgetForm({
  categories,
  onSave,
}: {
  categories: Category[];
  onSave: (categoryId: string, amount: number) => Promise<void>;
}) {
  const { t } = useTranslation();
  const [categoryId, setCategoryId] = useState("");
  const [amount, setAmount] = useState("10000");
  const submit = async (e: FormEvent) => {
    e.preventDefault();
    if (!categoryId) return;
    await onSave(categoryId, Number(amount));
  };
  return (
    <form onSubmit={submit} className="toolbar">
      <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)}>
        <option value="">{t("tx.category")}</option>
        {categories.map((c) => (
          <option key={c.id} value={c.id}>
            {c.icon} {c.name}
          </option>
        ))}
      </select>
      <input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} />
      <button type="submit">{t("extras.save")}</button>
    </form>
  );
}

function RecurringForm({
  categories,
  onCreate,
}: {
  categories: Category[];
  onCreate: (body: Partial<Recurring> & { kind: string; amount: number; schedule: string }) => Promise<void>;
}) {
  const { t } = useTranslation();
  const [kind, setKind] = useState("expense");
  const [amount, setAmount] = useState("199");
  const [desc, setDesc] = useState("Подписка");
  const [schedule, setSchedule] = useState<"daily" | "weekly" | "monthly">("monthly");
  const [next, setNext] = useState(new Date().toISOString().slice(0, 10));
  const [categoryId, setCategoryId] = useState("");

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    await onCreate({
      kind,
      amount: Number(amount),
      schedule,
      description: desc,
      next_run_at: next,
      category_id: categoryId || undefined,
      currency: "RUB",
      is_active: true,
    });
  };

  return (
    <form onSubmit={submit} className="toolbar">
      <select value={kind} onChange={(e) => setKind(e.target.value)}>
        <option value="income">{t("tx.income")}</option>
        <option value="expense">{t("tx.expense")}</option>
      </select>
      <select value={schedule} onChange={(e) => setSchedule(e.target.value as any)}>
        <option value="daily">{t("extras.daily")}</option>
        <option value="weekly">{t("extras.weekly")}</option>
        <option value="monthly">{t("extras.monthly")}</option>
      </select>
      <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)}>
        <option value="">{t("tx.category")}</option>
        {categories.map((c) => (
          <option key={c.id} value={c.id}>
            {c.icon} {c.name}
          </option>
        ))}
      </select>
      <input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} style={{ width: 120 }} />
      <input value={desc} onChange={(e) => setDesc(e.target.value)} placeholder={t("tx.description")} />
      <input type="date" value={next} onChange={(e) => setNext(e.target.value)} />
      <button type="submit">{t("extras.add")}</button>
    </form>
  );
}

