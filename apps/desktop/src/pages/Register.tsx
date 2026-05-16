import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useAuth } from "../context/AuthContext";

export default function Register() {
  const { t } = useTranslation();
  const { register } = useAuth();
  const nav = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [err, setErr] = useState("");

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setErr("");
    try {
      await register(email, password, name);
      nav("/");
    } catch {
      setErr(t("auth.error"));
    }
  };

  return (
    <section className="auth-page">
      <article className="card auth-card">
        <h1>{t("auth.register")}</h1>
        <form onSubmit={onSubmit}>
          <label>
            {t("auth.name")}
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Email
            <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </label>
          <label>
            {t("auth.password")}
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} minLength={6} required />
          </label>
          {err && <p className="down">{err}</p>}
          <button type="submit">{t("auth.register")}</button>
        </form>
        <p style={{ marginTop: "1rem" }}>
          <Link to="/login">{t("auth.to_login")}</Link>
        </p>
      </article>
    </section>
  );
}
