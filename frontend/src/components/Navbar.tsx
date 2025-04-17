import styles from './Navbar.module.css'; // opcional para estilizar
import Swal from 'sweetalert2';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

interface NavbarProps {
  userName: string;
  role: string;
  setUserName: (username: string) => void;
  setRole: (role: string) => void;
  handleLogout: () => void;
}

const Navbar = ({ userName, role, setUserName, setRole, handleLogout }: NavbarProps) => {
  const [isAdmin, setIsAdmin] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const storedUsername = localStorage.getItem("username");
    const storedRole = localStorage.getItem("role");

    if (storedRole && storedRole !== role) {
      setRole(storedRole);
    }

    if (storedRole === "admin") setIsAdmin(true);
    if (storedUsername && storedUsername !== userName) {
      setUserName(storedUsername);
    }
  }, [userName, setUserName, role, setRole]);

  const onLogoutClick = () => {
    Swal.fire({
      title: '¿Cerrar sesión?',
      text: 'Tu sesión actual se cerrará',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Sí, salir',
      cancelButtonText: 'Cancelar',
    }).then((result) => {
      if (result.isConfirmed) {
        handleLogout();
      }
    });
  };

  return (
    <nav className={styles.navbar}>
      <div className={styles.username}>Hola, {userName}</div>
      {isAdmin && (<button onClick={() => navigate("/admin")} className={styles.adminButton}>Administración</button>)}
      {isAdmin && (<button onClick={() => navigate("/boards")} className={styles.logoutButton}>Boards</button>)}
      <button className={styles.logoutButton} onClick={onLogoutClick}>
        Cerrar sesión
      </button>
    </nav>
  );
};

export default Navbar;
