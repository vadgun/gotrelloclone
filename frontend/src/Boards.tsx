import { useState, useEffect  } from "react";
import { getBoards } from "./api";

interface Board {
  id: string;
  name: string;
  owner_id: string;
  created_at: string;
}

function Boards() {
    // const [boards, setBoards] = useState([]);
    const [boards, setBoards] = useState<Board[]>([]);
    const [boardName, setBoardName] = useState("");

    const token = localStorage.getItem("token");

    const handleCreateBoard = async (e:any) => {
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
      getBoards().then(setBoards);
    }, []);
  
    return (
      <div>
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
          {boards.map((board:Board) => (
            <li key={board.id}>{board.name}</li>
          ))}
        </ul>
      </div>
    );
  }
  
  export default Boards;