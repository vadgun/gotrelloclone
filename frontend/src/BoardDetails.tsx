import { useState, useEffect } from "react";
import { useParams, useLocation } from "react-router-dom"
import { getTasks, createTask } from "./api";
import { useNavigate } from "react-router-dom";

interface Task {
    id: string;
    title: string;
    description: string;
    board_id: string;
  }

function BoardDetails({ handleLogout }: { handleLogout: () => void }) {
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
        getTasks(boardID).then(setTasks);
    }, []);

    const handleGoToBoards = () => {
        navigate("/boards");
      };

    return (
        <div>
        <h3>Tareas del Tablero : {boardName}  </h3>
        <button onClick={handleGoToBoards}>Ir a Tableros</button>
        <form onSubmit={handleNewTask}>
                <input 
                    type="text" 
                    placeholder="Título" 
                    value={title} 
                    onChange={(e) => setTitle(e.target.value)} 
                    required 
                />
                <input 
                    type="text" 
                    placeholder="Descripción" 
                    value={description} 
                    onChange={(e) => setDescription(e.target.value)} 
                    required 
                />
                <button type="submit">Agregar Tarea</button>
            </form>
        {tasks.length === 0 ? (
            <p>No hay tareas disponibles</p>
        ) : (
            <ul>
                {tasks.map((task: any) => (
                    <li key={task.id}>{task.title} - {task.description}</li>
                ))}
            </ul>
        )}
        <button onClick={handleLogout}>Cerrar Sesión</button>
    </div>
    );
}

export default BoardDetails;