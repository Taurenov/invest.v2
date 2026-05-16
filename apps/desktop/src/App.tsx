import { Navigate, Route, Routes } from "react-router-dom";
import { AuthProvider, useAuth } from "./context/AuthContext";
import { ThemeProvider } from "./context/ThemeContext";
import { useAlertsPoll } from "./hooks/useAlerts";
import Layout from "./components/Layout";
import Home from "./pages/Home";
import Transactions from "./pages/Transactions";
import Analytics from "./pages/Analytics";
import Markets from "./pages/Markets";
import Forecast from "./pages/Forecast";
import GoalsPage from "./pages/Goals";
import CalculatorPage from "./pages/Calculator";
import PortfolioPage from "./pages/Portfolio";
import SettingsPage from "./pages/Settings";
import Login from "./pages/Login";
import Register from "./pages/Register";

function PrivateRoutes() {
  const { user, loading } = useAuth();
  useAlertsPoll();
  if (loading) return <p className="main">…</p>;
  if (!user) return <Navigate to="/login" replace />;
  return (
    <ThemeProvider>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/transactions" element={<Transactions />} />
          <Route path="/analytics" element={<Analytics />} />
          <Route path="/markets" element={<Markets />} />
          <Route path="/portfolio" element={<PortfolioPage />} />
          <Route path="/forecast" element={<Forecast />} />
          <Route path="/goals" element={<GoalsPage />} />
          <Route path="/calculator" element={<CalculatorPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Routes>
      </Layout>
    </ThemeProvider>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/*" element={<PrivateRoutes />} />
      </Routes>
    </AuthProvider>
  );
}
