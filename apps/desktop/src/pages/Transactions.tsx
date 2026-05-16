import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import TransactionForm from "../components/TransactionForm";
import { deleteTransaction, fetchTransactions, formatMoney, type Transaction } from "../api/client";

export default function Transactions() {
  const { t } = useTranslation();
  const [tx, setTx] = useState<Transaction[]>([]);
  const [showForm, setShowForm] = useState(false);

  const reload = () => fetchTransactions().then((r) => setTx(r.data)).catch(console.error);

  useEffect(() => {
    reload();
  }, []);

  return (
    <>
      <section className="toolbar">
        <h1 style={{ margin: 0 }}>{t("nav.transactions")}</h1>
        <button type="button" onClick={() => setShowForm(true)}>
          + {t("tx.add")}
        </button>
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
                  <button type="button" className="secondary" onClick={() => deleteTransaction(row.id).then(reload)}>
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
