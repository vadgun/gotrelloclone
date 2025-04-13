import styles from './Navbar.module.css'; // opcional para estilizar
import Swal from 'sweetalert2';

interface NavbarProps {
  username: string;
  handleLogout: () => void;
}

const Navbar = ({ username, handleLogout }: NavbarProps) => {
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
      <div className={styles.username}>Hola, {username}</div>
      <button className={styles.logoutButton} onClick={onLogoutClick}>
        Cerrar sesión
      </button>
    </nav>
  );
};

export default Navbar;
