import { useState } from "react";
import { BrowserRouter as Router, Route, Routes, Navigate, useNavigate } from "react-router-dom";
import Login from './Login.tsx'
import Boards from './Boards.tsx'
import BoardDetails from './BoardDetails.tsx';

function AppWrapper() {
  const [token, setToken] = useState(localStorage.getItem("token") || "");
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    setToken("");
    navigate("/");
  };

  return (
    <Routes>
      {/* Si el usuario no está autenticado, lo enviamos a Login */}
      <Route path="/" element={token ? <Navigate to="/boards" /> : <Login token={token} setToken={setToken} />} />

      {/* Rutas privadas */}
      {token && (
        <>
          <Route path="/boards" element={<Boards handleLogout={handleLogout} />} />
          <Route path="/boards/:boardID" element={<BoardDetails handleLogout={handleLogout} />} />
        </>
      )}

      {/* Catch-all */}
      <Route path="*" element={<Navigate to={token ? "/boards" : "/"} />} />
    </Routes>
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
