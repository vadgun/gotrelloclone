import { useState, useEffect } from "react";
import { getBoards } from "./api";
import { Link } from "react-router-dom";
import styles from "./Boards.module.css";

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
    <div className={styles.container}> {/* Aplica el contenedor principal */}
      <div className={styles.header}>
        <h2>Bienvenido {username}</h2>
        <h3>Mis Tableros</h3>
      </div>
      <form onSubmit={handleCreateBoard} className={styles.createBoardForm}> {/* Aplica la clase al formulario */}
        <input
          type="text"
          placeholder="Nombre del tablero"
          value={boardName}
          onChange={(e) => setBoardName(e.target.value)}
          required
          className={styles.createBoardInput} // Aplica la clase al input
        />
        <button type="submit" className={styles.createBoardButton}>Crear Tablero</button> {/* Aplica la clase al botón */}
      </form>
      <ul className={styles.boardsList}> {/* Aplica la clase a la lista */}
        {boards.map((board: Board) => (
          <li key={board.id} className={styles.boardItem}> {/* Aplica la clase al elemento de la lista */}
            <Link to={`/boards/${board.id}`} state={{ boardName: board.name }} className={styles.boardLink}>
              {board.name}
            </Link> {/* Aplica la clase al enlace */}
          </li>
        ))}
      </ul>
      <button onClick={handleLogout} className={styles.logoutButton}>Cerrar Sesión</button> {/* Aplica la clase al botón de cerrar sesión */}
    </div>
  );
}

export default Boards;