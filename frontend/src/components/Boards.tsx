import { useState, useEffect } from "react";
import { getBoards } from "../api/api";
import { Link } from "react-router-dom";
import styles from "./Boards.module.css";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTrashCan, faPencil } from '@fortawesome/free-solid-svg-icons';
import Swal from 'sweetalert2';

interface Board {
  id: string;
  name: string;
  owner_id: string;
  created_at: string;
}

function Boards() {
  const [boards, setBoards] = useState<Board[]>([]);
  const [boardName, setBoardName] = useState("");
  const token = localStorage.getItem("token");
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
      Swal.fire('Error', 'Error al crear el tablero', 'error');
    }
  };

  useEffect(() => {
    getBoards(token).then(setBoards);
  }, []);

  const handleDeleteBoard = async (boardId: string) => {
    Swal.fire({
      title: '¿Estás seguro?',
      text: 'Esta acción no se puede deshacer',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Sí, eliminar',
      cancelButtonText: 'Cancelar',
    }).then(async (result) => {
      if (result.isConfirmed) {
        try {
          const response = await fetch(`http://localhost:8081/boards/${boardId}`, {
            method: "DELETE",
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });

          if (response.ok) {
            setBoards(boards.filter((board) => board.id !== boardId));
            Swal.fire('¡Eliminado!', 'El tablero ha sido eliminado', 'success');
          } else {
            Swal.fire('Error', 'Error al eliminar el tablero', 'error');
          }
        } catch (error) {
          Swal.fire('Error', 'Error en la conexión', 'error');
        }
      }
    });
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
      Swal.fire('Error', 'Error al actualizar el tablero', 'error');
    }
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
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
      {boards.length === 0 ? (
        <p className={styles.noBoards}>No hay tableros disponibles</p>
      ) : (
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
      )}
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