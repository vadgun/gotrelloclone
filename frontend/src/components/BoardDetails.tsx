import { useState, useEffect } from "react";
import { useParams, useLocation } from "react-router-dom"
import { getTasks, createTask } from "../api/api";
import { useNavigate } from "react-router-dom";
import styles from "./BoardDetails.module.css";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTrashCan, faPencil } from '@fortawesome/free-solid-svg-icons';
import Swal from 'sweetalert2';

interface Task {
  id: string;
  title: string;
  description: string;
  board_id: string;
  status: 'TODO' | 'IN_PROGRESS' | 'DONE';
}

function BoardDetails({ token }: { token: any; }) {
  const { boardID } = useParams();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const navigate = useNavigate();
  const location = useLocation();
  const boardName = location.state?.boardName || "Tablero sin nombre";
  const [isModalOpen, setIsOpenModal] = useState(false);
  const [taskToEdit, setTaskToEdit] = useState<Task | null>();
  const [editTitle, setEditTitle] = useState("");
  const [editDescription, setEditDescription] = useState("");
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [page, setPage] = useState(1);
  const [limit] = useState(5);
  const [totalPages, setTotalPages] = useState(1);


  const handleNewTask = async (e: React.FormEvent) => {
    e.preventDefault();

    const response = await createTask({ title, description, board_id: boardID })

    if (response.success) {
      const newTask = response.task;
      setTasks([...tasks, newTask]);
      Swal.fire('Creada!', response.message, 'success');
      setTitle("");
      setDescription("");
    } else {
      Swal.fire('Error', 'Error al agregar tarea', 'error');
    }
  };

  useEffect(() => {
    getTasks({ boardID, token, page, limit }).then(({ tasks, totalPages }) => {
      setTasks(tasks);
      setTotalPages(totalPages);
    });
  }, [page]);

  const handleGoToBoards = () => {
    navigate("/boards");
  };

  const openEditModal = (task: Task) => {
    setTaskToEdit(task);
    setEditTitle(task.title);
    setEditDescription(task.description);
    setIsOpenModal(true);
  }

  const closeModal = () => {
    setIsOpenModal(false);
    setTaskToEdit(null);
    setEditTitle("");
    setEditDescription("");
  }

  const handleUpdateTask = async () => {
    if (!taskToEdit) return;
    const response = await fetch(`http://localhost:8082/tasks/${taskToEdit.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        title: editTitle,
        description: editDescription,
      }),
    });

    if (response.ok) {
      const updated = tasks.map((t) =>
        t.id === taskToEdit.id
          ? { ...t, title: editTitle, description: editDescription }
          : t
      );
      setTasks(updated);
      closeModal();
    } else {
      Swal.fire('Error', 'Error al actualizar la tarea', 'error');
    }
  }

  const handleDeleteTask = async (taskId: string) => {
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
          const response = await fetch(`http://localhost:8082/tasks/${taskId}`, {
            method: "DELETE",
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
  
          if (response.ok) {
            setTasks(tasks.filter((task) => task.id !== taskId));
            Swal.fire('¡Eliminado!', 'La tarea ha sido eliminada', 'success');
          } else {
            Swal.fire('Error', 'Error al eliminar la tarea', 'error');
          }
        } catch (error) {
          Swal.fire('Error', 'Error en la conexión', 'error');
        }
      }
    });
  };

  const getHumanReadableStatus = (status: 'TODO' | 'IN_PROGRESS' | 'DONE') => {
    switch (status) {
      case 'TODO':
        return 'Por Hacer';
      case 'IN_PROGRESS':
        return 'En Progreso';
      case 'DONE':
        return 'Hecho';
      default:
        return status;
    }
  };

  const toggleTaskStatus = async (task: Task) => {
    const statusTransition: Record<'TODO' | 'IN_PROGRESS' | 'DONE', 'TODO' | 'IN_PROGRESS' | 'DONE'> = {
      'TODO': 'IN_PROGRESS',
      'IN_PROGRESS': 'DONE',
      'DONE': 'TODO'
    };

    const newStatus = statusTransition[task.status];

    const response = await fetch(`http://localhost:8082/tasks/${task.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        title: task.title,
        description: task.description,
        status: newStatus,
      }),
    });

    if (response.ok) {
      const updatedTasks = tasks.map((t) =>
        t.id === task.id ? { ...t, status: newStatus } : t
      );
      setTasks(updatedTasks);
      Swal.fire('Actualizado!', 'Actualizado a ' + getHumanReadableStatus(newStatus), 'success');
    } else {
      Swal.fire('Error', 'Error al actualizar el estado de la tarea', 'error');
    }
  };

  const filteredTasks = tasks.filter((task) => {
    const matchesSearch =
      task.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      task.description.toLowerCase().includes(searchTerm.toLowerCase());

    const matchesStatus = statusFilter ? task.status === statusFilter : true;

    return matchesSearch && matchesStatus;
  });

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h3 className={styles.boardTitle}>{boardName}</h3>
        <button onClick={handleGoToBoards} className={styles.backButton}>
          Ir a Tableros
        </button>
      </div>
      <form onSubmit={handleNewTask} className={styles.addTaskForm}>
        <input
          type="text"
          placeholder="Título"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          className={styles.taskInput}
        />
        <input
          type="text"
          placeholder="Descripción"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          required
          className={styles.taskInput}
        />
        <button type="submit" className={styles.addTaskButton}>
          Agregar Tarea
        </button>
      </form>
      <input
        type="text"
        placeholder="Buscar tareas..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        className={styles.searchInput}
      />
      <select
        value={statusFilter}
        onChange={(e) => setStatusFilter(e.target.value)}
        className={styles.statusSelect}
      >
        <option value="">Todos los estados</option>
        <option value="TODO">Pendiente</option>
        <option value="IN_PROGRESS">En progreso</option>
        <option value="DONE">Hecho</option>
      </select>
      {tasks.length === 0 ? (
        <p className={styles.noTasks}>No hay tareas disponibles</p>
      ) : filteredTasks.length === 0 ? (
        <p className={styles.noTasks}> No hay tareas que cumplan con el filtro</p>
      ) : (
        <ul className={styles.tasksList}>
          {filteredTasks.map((task: Task) => (
            <li key={task.id} className={styles.taskItem}>
              <input
                type="checkbox"
                checked={task.status === "DONE"}
                onChange={() => toggleTaskStatus(task)}
                className={styles.checkbox}
              />
              <div className={`${styles.taskInfo} ${task.status === "DONE" ? styles.taskCompleted : ""}`}>
                <div className={styles.taskTitle}>
                  <strong>Título:</strong> {task.title}
                </div>
                <div className={styles.taskDescription}>
                  <strong>Descripción:</strong> {task.description}
                </div>
                <div className={styles.taskStatus}>
                  <strong>Estado: </strong>{getHumanReadableStatus(task.status)}
                </div>
              </div>
              <div className={styles.taskActions}>
                <button
                  className={styles.editButton}
                  onClick={(e) => {
                    e.stopPropagation();
                    e.preventDefault();
                    openEditModal(task);
                  }}
                >
                  <FontAwesomeIcon icon={faPencil} />
                </button>
                <button
                  className={styles.deleteButton}
                  onClick={(e) => {
                    e.stopPropagation();
                    e.preventDefault();
                    handleDeleteTask(task.id);
                  }}
                >
                  <FontAwesomeIcon icon={faTrashCan} />
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
      {isModalOpen && (
        <div className={styles.modalOverlay}>
          <div className={styles.modal}>
            <h3>Editar Tarea</h3>
            <input
              type="text"
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
              placeholder="Título"
              className={styles.editInput}
            />
            <textarea
              value={editDescription}
              onChange={(e) => setEditDescription(e.target.value)}
              placeholder="Descripción"
              className={styles.editTextarea}
            ></textarea>
            <div className={styles.modalActions}>
              <button onClick={handleUpdateTask} className={styles.saveButton}>
                Guardar
              </button>
              <button onClick={closeModal} className={styles.cancelButton}>
                Cancelar
              </button>
            </div>
          </div>
        </div>
      )}
      <div className={styles.pagination}>
        <button onClick={() => setPage((prev) => Math.max(prev - 1, 1))} disabled={page === 1}>
          ← Anterior
        </button>

        <span>Página {page} de {totalPages}</span>

        <button onClick={() => setPage((prev) => Math.min(prev + 1, totalPages))} disabled={page === totalPages}>
          Siguiente →
        </button>
      </div>
    </div>

  );
}

export default BoardDetails;