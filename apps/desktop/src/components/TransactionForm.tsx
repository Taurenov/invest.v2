import { FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Category, createTransaction, fetchCategories } from "../api/client";

type Props = {
  onSaved: () => void;
  onClose: () => void;
};

export default function TransactionForm({ onSaved, onClose }: Props) {
  const { t } = useTranslation();
  const [categories, setCategories] = useState<Category[]>([]);
  const [kind, setKind] = useState<"income" | "expense">("expense");
  const [amount, setAmount] = useState("");
  const [categoryId, setCategoryId] = useState("");
  const [description, setDescription] = useState("");
  const [date, setDate] = useState(new Date().toISOString().slice(0, 10));

  useEffect(() => {
    fetchCategories().then((r) => setCategories(r.data)).catch(console.error);
  }, []);

  const filtered = categories.filter((c) => c.kind === kind);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    const val = parseFloat(amount);
    if (!val || val <= 0) return;
    await createTransaction({
      kind,
      amount: val,
      currency: "RUB",
      description,
      category_id: categoryId || undefined,
      occurred_at: new Date(date).toISOString(),
    } as never);
    onSaved();
    onClose();
  };

  return (
    <section className="drawer">
      <h2>{t("tx.add")}</h2>
      <form onSubmit={submit}>
        <label>
          {t("tx.kind")}
          <select value={kind} onChange={(e) => setKind(e.target.value as "income" | "expense")}>
            <option value="income">{t("tx.income")}</option>
            <option value="expense">{t("tx.expense")}</option>
          </select>
        </label>
        <label>
          {t("tx.amount")}
          <input type="number" min="0" step="0.01" value={amount} onChange={(e) => setAmount(e.target.value)} required />
        </label>
        <label>
          {t("tx.category")}
          <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)}>
            <option value="">—</option>
            {filtered.map((c) => (
              <option key={c.id} value={c.id}>
                {c.icon} {c.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          {t("tx.description")}
          <input value={description} onChange={(e) => setDescription(e.target.value)} />
        </label>
        <label>
          {t("tx.date")}
          <input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
        </label>
        <section className="toolbar">
          <button type="submit">{t("tx.save")}</button>
          <button type="button" className="secondary" onClick={onClose}>
            {t("tx.cancel")}
          </button>
        </section>
      </form>
    </section>
  );
}
