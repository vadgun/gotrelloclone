import { useState, useEffect } from "react";
import { getBoards } from "../api/api";
import { Link } from "react-router-dom";
import styles from "./Boards.module.css";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTrashCan, faPencil } from '@fortawesome/free-solid-svg-icons';

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
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editBoard, setEditBoard] = useState<Board | null>(null);
  const [editName, setEditName] = useState("");

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

  const handleDeleteBoard = async (boardId: string) => {
    const confirmed = window.confirm("¿Estás seguro de eliminar este tablero?");
    if (!confirmed) return;

    const response = await fetch(`http://localhost:8081/boards/${boardId}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (response.ok) {
      setBoards(boards.filter((board) => board.id !== boardId));
    } else {
      alert("Error al eliminar el tablero");
    }
  };

  const openEditModal = (board: Board) => {
    setEditBoard(board);
    setEditName(board.name);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setEditBoard(null)
    setEditName("");
  }

  const handleUpdateBoard = async () => {
    if (!editBoard) return;

    const response = await fetch(`http://localhost:8081/boards/${editBoard.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ name: editName }),
    });

    if (response.ok) {
      const updatedBoards = boards.map((b) =>
        b.id === editBoard.id ? { ...b, name: editName } : b
      );
      setBoards(updatedBoards);
      closeModal();
    } else {
      alert("Error al actualizar el tablero");
    }
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2>Bienvenido {username}</h2>
        <h3>Mis Tableros</h3>
      </div>
      <form onSubmit={handleCreateBoard} className={styles.createBoardForm}>
        <input
          type="text"
          placeholder="Nombre del tablero"
          value={boardName}
          onChange={(e) => setBoardName(e.target.value)}
          required
          className={styles.createBoardInput}
        />
        <button type="submit" className={styles.createBoardButton}>Crear Tablero</button>
      </form>
      <ul className={styles.boardsList}>
        {boards.map((board: Board) => (
          <Link
            to={`/boards/${board.id}`}
            state={{ boardName: board.name }}
            className={styles.boardLink}
            key={board.id}
          >
            <li className={styles.boardItem}>
              <span className={styles.boardName}>{board.name}</span>
              <button className={styles.editButton} onClick={(e) => {
                e.stopPropagation();
                e.preventDefault();
                openEditModal(board);
              }}>
                <FontAwesomeIcon icon={faPencil} />
              </button>
              <button className={styles.deleteButton} onClick={(e) => {
                e.stopPropagation();
                e.preventDefault();
                handleDeleteBoard(board.id);
              }}>
                <FontAwesomeIcon icon={faTrashCan} />
              </button>
            </li>
          </Link>
        ))}
      </ul>
      <button onClick={handleLogout} className={styles.logoutButton}>Cerrar Sesión</button>
      {isModalOpen && (
        <div className={styles.modalOverlay}>
          <div className={styles.modal}>
            <h3>Editar Tablero</h3>
            <input
              type="text"
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
              className={styles.editInput}
            />
            <div className={styles.modalActions}>
              <button onClick={handleUpdateBoard} className={styles.saveButton}>Guardar</button>
              <button onClick={closeModal} className={styles.cancelButton}>Cancelar</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default Boards;