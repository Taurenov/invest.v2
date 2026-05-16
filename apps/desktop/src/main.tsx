import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import "./styles/global.css";
import "./i18n";
import { getToken, getUser } from "./api/client";

// Dev: без логина — demo JWT path через dev-token (backend JWTOrDevToken)
if (!getToken() && !getUser() && import.meta.env.VITE_API_TOKEN) {
  localStorage.setItem("fin_token", import.meta.env.VITE_API_TOKEN);
  localStorage.setItem(
    "fin_user",
    JSON.stringify({
      id: "00000000-0000-4000-8000-000000000001",
      email: "demo@fin-helper.local",
      display_name: "Demo",
    })
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>
);
