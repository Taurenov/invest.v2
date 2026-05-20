import { NavLink } from "react-router-dom";
import { useTranslation } from "react-i18next";
import i18n from "../i18n";
import { useAuth } from "../context/AuthContext";

export default function Layout({ children }: { children: React.ReactNode }) {
  const { t } = useTranslation();
  const { user, logout } = useAuth();

  const setLocale = (lng: string) => {
    i18n.changeLanguage(lng);
    localStorage.setItem("locale", lng);
  };

  return (
    <section className="app-shell">
      <aside className="sidebar">
        <p style={{ padding: "0 1rem 1rem", fontWeight: 700 }} className="nav-label">
          {t("app_name")}
        </p>
        <NavLink to="/" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`} end>
          <span>⌂</span>
          <span className="nav-label">{t("nav.home")}</span>
        </NavLink>
        <NavLink to="/transactions" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>⇄</span>
          <span className="nav-label">{t("nav.transactions")}</span>
        </NavLink>
        <NavLink to="/analytics" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>◔</span>
          <span className="nav-label">{t("nav.analytics")}</span>
        </NavLink>
        <NavLink to="/markets" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>📈</span>
          <span className="nav-label">{t("nav.markets")}</span>
        </NavLink>
        <NavLink to="/portfolio" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>💼</span>
          <span className="nav-label">{t("nav.portfolio")}</span>
        </NavLink>
        <NavLink to="/goals" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>🎯</span>
          <span className="nav-label">{t("nav.goals")}</span>
        </NavLink>
        <NavLink to="/calculator" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>🧮</span>
          <span className="nav-label">{t("nav.calculator")}</span>
        </NavLink>
        <NavLink to="/forecast" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>✦</span>
          <span className="nav-label">{t("nav.forecast")}</span>
        </NavLink>
        <NavLink to="/settings" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>⚙</span>
          <span className="nav-label">{t("nav.settings")}</span>
        </NavLink>
        <NavLink to="/extras" className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}>
          <span>✨</span>
          <span className="nav-label">{t("nav.extras")}</span>
        </NavLink>
        <p style={{ marginTop: "auto", padding: "1rem" }} className="nav-label">
          <small>{user?.display_name || user?.email}</small>
          <br />
          <button type="button" className="secondary" onClick={() => setLocale("ru")}>
            RU
          </button>{" "}
          <button type="button" className="secondary" onClick={() => setLocale("en")}>
            EN
          </button>
          <br />
          <button type="button" className="secondary" style={{ marginTop: "0.5rem" }} onClick={logout}>
            {t("auth.logout")}
          </button>
        </p>
      </aside>
      <main className="main">{children}</main>
    </section>
  );
}
