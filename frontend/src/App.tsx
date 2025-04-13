import { useState } from "react";
import { BrowserRouter as Router, Route, Routes, Navigate, useNavigate } from "react-router-dom";
import Login from './components/Login.tsx'
import Boards from './components/Boards.tsx'
import BoardDetails from './components/BoardDetails.tsx';
import Navbar from "./components/Navbar";
import Footer from "./components/Footer";
import './index.css'

function AppWrapper() {
  const [token, setToken] = useState(localStorage.getItem("token") || "");
  const [username] = useState(localStorage.getItem("username") || "");
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    setToken("");
    navigate("/");
  };

  return (
    <div className="app-container">
      {token && <Navbar username={username} handleLogout={handleLogout} />}
      <div className="main-content">
        <Routes>
          {/* Si el usuario no está autenticado, lo enviamos a Login */}
          <Route path="/" element={token ? <Navigate to="/boards" /> : <Login token={token} setToken={setToken} />} />

          {/* Rutas privadas */}
          {token && (
            <>
              <Route path="/boards" element={<Boards />} />
              <Route path="/boards/:boardID" element={<BoardDetails token={token} />} />
            </>
          )}

          {/* Catch-all */}
          <Route path="*" element={<Navigate to={token ? "/boards" : "/"} />} />
        </Routes>
      </div>
      <Footer />
    </div>
  );
}


function App() {
  return (
    <Router>
      <AppWrapper />
    </Router>
  );
}

export default App;
