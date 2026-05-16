import { useTranslation } from "react-i18next";

export default function Disclaimer() {
  const { t } = useTranslation();
  return <p className="disclaimer" role="note">{t("disclaimer")}</p>;
}
