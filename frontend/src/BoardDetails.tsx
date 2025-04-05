import { useState, useEffect } from "react";
import { useParams, useLocation } from "react-router-dom"
import { getTasks, createTask } from "./api";
import { useNavigate } from "react-router-dom";
import styles from "./BoardDetails.module.css";

interface Task {
    id: string;
    title: string;
    description: string;
    board_id: string;
  }

  function BoardDetails({ handleLogout, token }: { handleLogout: () => void; token: any; }) {
    const { boardID } = useParams();
    const [tasks, setTasks] = useState<Task[]>([]);
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const navigate = useNavigate();
    const location = useLocation();
    const boardName = location.state?.boardName || "Tablero sin nombre";

    const handleNewTask = async (e: React.FormEvent) => {
        e.preventDefault();

        const response = await createTask({ title, description, board_id: boardID })

        if (response.success) {
            const newTask = response.task;
            setTasks([...tasks, newTask]);
            alert(response.message);
            setTitle("");
            setDescription("");
        } else {
            alert("Error al agregar tarea");
        }
    };

    useEffect(() => {
        getTasks({boardID, token:token}).then(setTasks);
    }, []);

    const handleGoToBoards = () => {
        navigate("/boards");
      };

      return (
        <div className={styles.container}>
          <div className={styles.header}>
            <h3 className={styles.boardTitle}>Tareas del Tablero: {boardName}</h3>
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
          {tasks.length === 0 ? (
            <p className={styles.noTasks}>No hay tareas disponibles</p>
          ) : (
            <ul className={styles.tasksList}>
              {tasks.map((task: any) => (
                <li key={task.id} className={styles.taskItem}>
                  <span>{task.title} - {task.description}</span>
                  {/* Aquí podrías agregar más funcionalidades a cada tarea, como botones de editar o eliminar */}
                </li>
              ))}
            </ul>
          )}
          <button onClick={handleLogout} className={styles.logoutButton}>
            Cerrar Sesión
          </button>
        </div>
      );
}

export default BoardDetails;