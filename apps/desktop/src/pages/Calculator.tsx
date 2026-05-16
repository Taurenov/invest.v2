import { FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";
import { calcCAGR, calcROI, calcSavings, formatMoney } from "../api/client";

export default function CalculatorPage() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<"roi" | "cagr" | "savings">("roi");
  const [result, setResult] = useState<string>("");

  const onROI = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const fd = new FormData(e.currentTarget);
    const r = await calcROI(Number(fd.get("initial")), Number(fd.get("current")));
    setResult(`${r.data.roi_percent.toFixed(2)}%`);
  };

  const onCAGR = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const fd = new FormData(e.currentTarget);
    const r = await calcCAGR(Number(fd.get("initial")), Number(fd.get("final")), Number(fd.get("years")));
    setResult(`${r.data.cagr_percent.toFixed(2)}% ${t("calc.per_year")}`);
  };

  const onSavings = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const fd = new FormData(e.currentTarget);
    const r = await calcSavings(
      Number(fd.get("monthly")),
      Number(fd.get("rate")),
      Number(fd.get("months")),
      Number(fd.get("initial") || 0)
    );
    setResult(formatMoney(r.data.future_value));
  };

  return (
    <>
      <h1>{t("calc.title")}</h1>
      <section className="toolbar">
        <button type="button" className={tab === "roi" ? "" : "secondary"} onClick={() => setTab("roi")}>
          ROI
        </button>
        <button type="button" className={tab === "cagr" ? "" : "secondary"} onClick={() => setTab("cagr")}>
          CAGR
        </button>
        <button type="button" className={tab === "savings" ? "" : "secondary"} onClick={() => setTab("savings")}>
          {t("calc.savings")}
        </button>
      </section>

      <article className="card">
        {tab === "roi" && (
          <form onSubmit={onROI}>
            <label>
              {t("calc.initial")}
              <input name="initial" type="number" defaultValue={100000} required />
            </label>
            <label>
              {t("calc.current")}
              <input name="current" type="number" defaultValue={125000} required />
            </label>
            <button type="submit">{t("calc.run")}</button>
          </form>
        )}
        {tab === "cagr" && (
          <form onSubmit={onCAGR}>
            <label>
              {t("calc.initial")}
              <input name="initial" type="number" defaultValue={100000} required />
            </label>
            <label>
              {t("calc.final")}
              <input name="final" type="number" defaultValue={133100} required />
            </label>
            <label>
              {t("calc.years")}
              <input name="years" type="number" defaultValue={3} step="0.1" required />
            </label>
            <button type="submit">{t("calc.run")}</button>
          </form>
        )}
        {tab === "savings" && (
          <form onSubmit={onSavings}>
            <label>
              {t("calc.monthly")}
              <input name="monthly" type="number" defaultValue={15000} required />
            </label>
            <label>
              {t("calc.rate")}
              <input name="rate" type="number" defaultValue={12} step="0.1" required />
            </label>
            <label>
              {t("calc.months")}
              <input name="months" type="number" defaultValue={24} required />
            </label>
            <label>
              {t("calc.initial_balance")}
              <input name="initial" type="number" defaultValue={0} />
            </label>
            <button type="submit">{t("calc.run")}</button>
          </form>
        )}
        {result && (
          <p className="value" style={{ marginTop: "1.25rem" }}>
            {t("calc.result")}: {result}
          </p>
        )}
      </article>
    </>
  );
}
