import { useTranslation } from "react-i18next";
import { useTheme } from "../context/ThemeContext";

export default function SettingsPage() {
  const { t } = useTranslation();
  const { settings, apply } = useTheme();

  if (!settings) return <p>{t("common.loading")}</p>;

  return (
    <>
      <h1>{t("settings.title")}</h1>
      <article className="card">
        <label>
          {t("settings.theme")}
          <select
            value={settings.theme}
            onChange={(e) => apply({ theme: e.target.value as "dark" | "light" | "system" })}
          >
            <option value="dark">{t("settings.theme_dark")}</option>
            <option value="light">{t("settings.theme_light")}</option>
            <option value="system">{t("settings.theme_system")}</option>
          </select>
        </label>
        <label style={{ display: "block", marginTop: "1rem" }}>
          {t("settings.locale")}
          <select value={settings.locale} onChange={(e) => apply({ locale: e.target.value })}>
            <option value="ru">Русский</option>
            <option value="en">English</option>
          </select>
        </label>
        <label style={{ display: "block", marginTop: "1rem" }}>
          {t("settings.currency")}
          <select
            value={settings.base_currency}
            onChange={(e) => apply({ base_currency: e.target.value })}
          >
            <option value="RUB">RUB</option>
            <option value="USD">USD</option>
            <option value="EUR">EUR</option>
          </select>
        </label>
        <label style={{ display: "block", marginTop: "1rem" }}>
          {t("settings.timezone")}
          <input
            value={settings.timezone}
            onChange={(e) => apply({ timezone: e.target.value })}
          />
        </label>
      </article>
    </>
  );
}
