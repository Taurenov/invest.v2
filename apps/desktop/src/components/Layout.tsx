import { NavLink } from "react-router-dom";
import { useTranslation } from "react-i18next";
import i18n from "../i18n";
import { useAuth } from "../context/AuthContext";
import NavIcon from "./NavIcon";
import BrandLogo from "./BrandLogo";

const NAV = [
  { to: "/", icon: "home" as const, key: "nav.home", end: true },
  { to: "/transactions", icon: "transactions" as const, key: "nav.transactions" },
  { to: "/analytics", icon: "analytics" as const, key: "nav.analytics" },
  { to: "/markets", icon: "markets" as const, key: "nav.markets" },
  { to: "/portfolio", icon: "portfolio" as const, key: "nav.portfolio" },
  { to: "/goals", icon: "goals" as const, key: "nav.goals" },
  { to: "/calculator", icon: "calculator" as const, key: "nav.calculator" },
  { to: "/forecast", icon: "forecast" as const, key: "nav.forecast" },
  { to: "/extras", icon: "extras" as const, key: "nav.extras" },
  { to: "/settings", icon: "settings" as const, key: "nav.settings" },
];

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
        <div className="sidebar-brand">
          <BrandLogo />
          <span>{t("app_name")}</span>
        </div>
        {NAV.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={"end" in item ? item.end : false}
            className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}
          >
            <NavIcon name={item.icon} />
            <span className="nav-label">{t(item.key)}</span>
          </NavLink>
        ))}
        <div style={{ marginTop: "auto", padding: "1rem 0.85rem" }} className="nav-label">
          <small style={{ color: "var(--muted)" }}>{user?.display_name || user?.email}</small>
          <div style={{ marginTop: "0.65rem", display: "flex", gap: "0.35rem" }}>
            <button type="button" className="secondary" onClick={() => setLocale("ru")}>
              RU
            </button>
            <button type="button" className="secondary" onClick={() => setLocale("en")}>
              EN
            </button>
          </div>
          <button type="button" className="secondary" style={{ marginTop: "0.5rem", width: "100%" }} onClick={logout}>
            {t("auth.logout")}
          </button>
        </div>
      </aside>
      <main className="main">{children}</main>
    </section>
  );
}
