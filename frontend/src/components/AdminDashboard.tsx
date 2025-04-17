import { useState, useEffect } from "react";
import { getAdminUsers, getAdminBoards, getAdminTasks } from '../api/api';
import styles from './AdminDashboard.module.css'

interface User {
    id: string;
    name: string;
    role: string;
}

const AdminDashboard = ({ token }: { token: string }) => {
    const [users, setUsers] = useState<User[]>([]);
    const [boards, setBoards] = useState([]);
    const [tasks, setTasks] = useState([]);

    useEffect(() => {
        getAdminUsers(token).then(setUsers);
        getAdminBoards(token).then(setBoards);
        getAdminTasks(token).then(setTasks);
    }, [token]);

    return (
        <div className={styles.adminPanel}> {/* Aplica el contenedor principal si lo deseas */}
          <h1>Panel de Administración</h1>
    
          <section className={styles.adminPanelSection}>
            <h2>Usuarios</h2>
            <ul>
              {users.map((user: any) => (
                <li key={user.id}>{user.name} - {user.role}</li>
              ))}
            </ul>
          </section>
    
          <section className={styles.adminPanelSection}>
            <h2>Boards</h2>
            <ul>
              {boards.map((board: any) => (
                <li key={board.id}>{board.name} (Owner: {board.owner_name})</li>
              ))}
            </ul>
          </section>
    
          <section className={styles.adminPanelSection}>
            <h2>Tareas</h2>
            <ul>
              {tasks.map((task: any) => (
                <li key={task.id}>{task.title} : {task.description}</li>
              ))}
            </ul>
          </section>
        </div>
      );
}

export default AdminDashboard;