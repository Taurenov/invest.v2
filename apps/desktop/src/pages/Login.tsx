import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useAuth } from "../context/AuthContext";

export default function Login() {
  const { t } = useTranslation();
  const { login } = useAuth();
  const nav = useNavigate();
  const [email, setEmail] = useState("demo@fin-helper.local");
  const [password, setPassword] = useState("demo123");
  const [err, setErr] = useState("");

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setErr("");
    try {
      await login(email, password);
      nav("/");
    } catch {
      setErr(t("auth.error"));
    }
  };

  return (
    <section className="auth-page">
      <article className="card auth-card">
        <h1>{t("auth.login")}</h1>
        <form onSubmit={onSubmit}>
          <label>
            Email
            <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </label>
          <label>
            {t("auth.password")}
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </label>
          {err && <p className="down">{err}</p>}
          <button type="submit">{t("auth.login")}</button>
        </form>
        <p style={{ marginTop: "1rem" }}>
          <Link to="/register">{t("auth.to_register")}</Link>
        </p>
        <p className="label" style={{ marginTop: "0.75rem" }}>
          Dev: token <code>dev-token</code> без входа (memory mode)
        </p>
      </article>
    </section>
  );
}
