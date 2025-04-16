import { useState, useEffect } from "react";
import styles from './AdminDashboard.module.css'

interface User {
    id: string;
    name: string;
    role: string;
}

const AdminDashboard = ({ token }: { token: string }) => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const response = await fetch("http://localhost:8080/admin/users", {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    }
                });
                const data = await response.json();
                setUsers(data);
            } catch (err) {
                console.log("Error al obtener usuarios", err);
            } finally {
                setLoading(false);
            }
        };
        fetchUsers();
    }, [token])

    return (
        <div className={styles.dashboard}>
            <h2>Panel de Administración</h2>
            {loading ? (
                <p>Cargando usuarios...</p>
            ) : (
                <table className={styles.table}>
                    <thead>
                        <tr>
                            <th>Nombre</th>
                            <th>Rol</th>
                        </tr>
                    </thead>
                    <tbody>
                        {users.map(user => (
                            <tr key={user.id}>
                                <td>{user.name}</td>
                                <td>{user.role}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    );
}

export default AdminDashboard;