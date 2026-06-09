import { FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Category, createTransaction, fetchCategories } from "../api/client";
import { enqueue } from "../offline/store";
import { Flash } from "./Flash";

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
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [offlineSaved, setOfflineSaved] = useState(false);

  useEffect(() => {
    fetchCategories()
      .then((r) => setCategories(r.data))
      .catch((e) => setError(e instanceof Error ? e.message : "Ошибка загрузки категорий"));
  }, []);

  const filtered = categories.filter((c) => c.kind === kind);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setOfflineSaved(false);
    const val = parseFloat(amount);
    if (!val || val <= 0) {
      setError("Укажите сумму больше 0");
      return;
    }
    const body: Record<string, unknown> = {
      kind,
      amount: val,
      currency: "RUB",
      description,
      occurred_at: new Date(`${date}T12:00:00`).toISOString(),
    };
    if (categoryId) body.category_id = categoryId;

    setSaving(true);
    try {
      await createTransaction(body as Parameters<typeof createTransaction>[0]);
      onSaved();
      onClose();
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Не удалось сохранить";
      try {
        enqueue({ type: "create", payload: body as never });
        setOfflineSaved(true);
        onSaved();
        onClose();
      } catch {
        setError(msg);
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="modal-backdrop" onClick={onClose} role="presentation">
      <section className="modal-sheet" onClick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
        <h2>{t("tx.add")}</h2>
        {error && <Flash kind="error">{error}</Flash>}
        {offlineSaved && <Flash kind="info">Сохранено офлайн — синхронизируется при подключении к API</Flash>}
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
            <input
              type="number"
              min="0"
              step="0.01"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              required
              autoFocus
            />
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
          <section className="toolbar" style={{ marginTop: "1rem" }}>
            <button type="submit" disabled={saving}>
              {saving ? "…" : t("tx.save")}
            </button>
            <button type="button" className="secondary" onClick={onClose}>
              {t("tx.cancel")}
            </button>
          </section>
        </form>
      </section>
    </div>
  );
}
