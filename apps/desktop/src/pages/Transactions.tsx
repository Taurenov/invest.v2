import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import TransactionForm from "../components/TransactionForm";
import { createTransaction, deleteTransaction, formatMoney, searchTransactions, type Transaction } from "../api/client";
import { enqueue, flushQueue, getCachedTransactions, setCachedTransactions } from "../offline/store";

export default function Transactions() {
  const { t } = useTranslation();
  const [tx, setTx] = useState<Transaction[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [offline, setOffline] = useState(false);
  const [q, setQ] = useState("");
  const [kind, setKind] = useState("");

  const reload = async () => {
    try {
      await flushQueue(createTransaction, deleteTransaction);
      const r = q || kind ? await searchTransactions({ q, kind }) : await searchTransactions({});
      setTx(r.data);
      setCachedTransactions(r.data);
      setOffline(false);
    } catch (e) {
      console.error(e);
      setTx(getCachedTransactions());
      setOffline(true);
    }
  };

  useEffect(() => {
    reload();
    const onOnline = () => {
      void reload();
    };
    window.addEventListener("online", onOnline);
    return () => window.removeEventListener("online", onOnline);
  }, []);

  return (
    <>
      <section className="toolbar">
        <h1 style={{ margin: 0 }}>{t("nav.transactions")}</h1>
        <input value={q} onChange={(e) => setQ(e.target.value)} placeholder={t("tx.search")} />
        <select value={kind} onChange={(e) => setKind(e.target.value)}>
          <option value="">{t("tx.all")}</option>
          <option value="income">{t("tx.income")}</option>
          <option value="expense">{t("tx.expense")}</option>
        </select>
        <button type="button" className="secondary" onClick={reload}>
          {t("tx.apply")}
        </button>
        <button type="button" onClick={() => setShowForm(true)}>
          + {t("tx.add")}
        </button>
        <label className="secondary btn" style={{ cursor: "pointer" }}>
          {t("tx.import_csv")}
          <input
            type="file"
            accept=".csv,text/csv"
            style={{ display: "none" }}
            onChange={async (e) => {
              const file = e.target.files?.[0];
              if (!file) return;
              const text = await file.text();
              const rows = text.split(/\r?\n/).filter(Boolean);
              // CSV format: date,kind,amount,description
              for (const line of rows.slice(1)) {
                const [date, kind, amount, description] = line.split(",");
                const payload = {
                  kind: (kind?.trim() || "expense") as "income" | "expense",
                  amount: Number(amount),
                  currency: "RUB",
                  description: (description || "").trim(),
                  occurred_at: new Date(date).toISOString(),
                };
                if (!payload.amount || Number.isNaN(payload.amount)) continue;
                try {
                  await createTransaction(payload);
                } catch {
                  enqueue({ type: "create", payload });
                }
              }
              await reload();
            }}
          />
        </label>
        {offline && <span className="label">offline mode</span>}
      </section>
      {showForm && <TransactionForm onSaved={reload} onClose={() => setShowForm(false)} />}
      <article className="card">
        <table>
          <thead>
            <tr>
              <th>Тип</th>
              <th>Описание</th>
              <th>Дата</th>
              <th>Сумма</th>
            </tr>
          </thead>
          <tbody>
            {tx.map((row) => (
              <tr key={row.id}>
                <td>{row.kind === "income" ? "Доход" : "Расход"}</td>
                <td>{row.description}</td>
                <td>{new Date(row.occurred_at).toLocaleString()}</td>
                <td className={row.kind === "income" ? "up" : "down"}>
                  {formatMoney(row.amount)}
                </td>
                <td>
                  <button
                    type="button"
                    className="secondary"
                    onClick={async () => {
                      try {
                        await deleteTransaction(row.id);
                      } catch {
                        enqueue({ type: "delete", id: row.id });
                      }
                      await reload();
                    }}
                  >
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
