import { useState, useEffect } from "react";
import { getBoards } from "./api";
import { Link } from "react-router-dom";

interface Board {
  id: string;
  name: string;
  owner_id: string;
  created_at: string;
}

function Boards({ handleLogout }: { handleLogout: () => void }) {
  const [boards, setBoards] = useState<Board[]>([]);
  const [boardName, setBoardName] = useState("");
  const token = localStorage.getItem("token");
  const username = localStorage.getItem("username");

  const handleCreateBoard = async (e: any) => {
    e.preventDefault();

    const response = await fetch("http://localhost:8081/boards", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ name: boardName }),
    });

    if (response.ok) {
      const data = await response.json();
      const newBoard = data.board;
      setBoards([...boards, newBoard]); // Agrega el nuevo tablero a la lista
      setBoardName(""); // Limpia el input
    } else {
      alert("Error al crear el tablero");
    }
  };

  useEffect(() => {
    getBoards(token).then(setBoards);
  }, []);

  return (
    <div>
      <h2>Bienvenido {username}</h2>
      <h3>Mis Tableros</h3>
      <form onSubmit={handleCreateBoard}>
        <input
          type="text"
          placeholder="Nombre del tablero"
          value={boardName}
          onChange={(e) => setBoardName(e.target.value)}
          required
        />
        <button type="submit">Crear Tablero</button>
      </form>
      <ul>
        {boards.map((board: Board) => (
          <li key={board.id}>
            <Link to={`/boards/${board.id}`} state={{ boardName: board.name }}>{board.name}</Link> {/* 👈 Enlace a detalles del tablero */}
          </li>
        ))}
      </ul>
      <button onClick={handleLogout}>Cerrar Sesión</button>
    </div>
  );
}

export default Boards;