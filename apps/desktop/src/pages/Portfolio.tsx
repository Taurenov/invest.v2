import { FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  addHolding,
  fetchPortfolio,
  formatMoney,
  removeHolding,
  type Portfolio,
} from "../api/client";

export default function PortfolioPage() {
  const { t } = useTranslation();
  const [p, setP] = useState<Portfolio | null>(null);
  const [symbol, setSymbol] = useState("SBER");
  const [qty, setQty] = useState("10");
  const [cost, setCost] = useState("250");

  const load = () => fetchPortfolio().then((r) => setP(r.data)).catch(console.error);
  useEffect(() => {
    load();
  }, []);

  const onAdd = async (e: FormEvent) => {
    e.preventDefault();
    await addHolding({ symbol, quantity: parseFloat(qty), avg_cost: parseFloat(cost) });
    load();
  };

  if (!p) return <p>{t("common.loading")}</p>;

  return (
    <>
      <h1>{t("portfolio.title")}</h1>
      <section className="cards">
        <article className="card">
          <p className="label">{t("portfolio.value")}</p>
          <p className="value">{formatMoney(p.total_value)}</p>
        </article>
        <article className="card">
          <p className="label">{t("portfolio.cost")}</p>
          <p className="value">{formatMoney(p.total_cost)}</p>
        </article>
        <article className="card">
          <p className="label">P&amp;L</p>
          <p className={`value ${p.total_pnl >= 0 ? "up" : "down"}`}>{formatMoney(p.total_pnl)}</p>
        </article>
      </section>

      <article className="card" style={{ marginBottom: "1rem" }}>
        <p className="label">{t("portfolio.add")}</p>
        <form onSubmit={onAdd} className="toolbar">
          <input value={symbol} onChange={(e) => setSymbol(e.target.value.toUpperCase())} placeholder="SBER" />
          <input type="number" value={qty} onChange={(e) => setQty(e.target.value)} placeholder="Qty" />
          <input type="number" value={cost} onChange={(e) => setCost(e.target.value)} placeholder="Avg" />
          <button type="submit">{t("portfolio.add_btn")}</button>
        </form>
      </article>

      <article className="card">
        <table>
          <thead>
            <tr>
              <th>{t("portfolio.ticker")}</th>
              <th>{t("portfolio.qty")}</th>
              <th>{t("portfolio.price")}</th>
              <th>P&amp;L</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {p.holdings.map((h) => (
              <tr key={h.id}>
                <td>{h.symbol}</td>
                <td>{h.quantity}</td>
                <td>{h.current_price?.toFixed(2) ?? "—"}</td>
                <td className={(h.pnl ?? 0) >= 0 ? "up" : "down"}>
                  {h.pnl != null ? `${formatMoney(h.pnl)} (${(h.pnl_percent ?? 0).toFixed(1)}%)` : "—"}
                </td>
                <td>
                  <button type="button" className="secondary" onClick={() => removeHolding(h.id).then(load)}>
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
