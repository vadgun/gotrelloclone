import styles from './Navbar.module.css'; // opcional para estilizar
import Swal from 'sweetalert2';
import { useEffect } from 'react';

interface NavbarProps {
  userName: string;
  setUserName: (username: string) => void;
  handleLogout: () => void;
}

const Navbar = ({ userName, setUserName, handleLogout }: NavbarProps) => {
  useEffect(() => {
    const storedUsername = localStorage.getItem("username");
    if (storedUsername && storedUsername !== userName) {
      setUserName(storedUsername);
    }
  }, [userName, setUserName]);

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
      <button className={styles.logoutButton} onClick={onLogoutClick}>
        Cerrar sesión
      </button>
    </nav>
  );
};

export default Navbar;
