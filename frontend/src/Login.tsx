import { useState } from "react";
import { loginUser, registerUser } from "./api";
import styles from './Login.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';


function Login({ token, setToken }: { token: string; setToken: (token: string) => void }) {

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [phone, setPhone] = useState("");
  const [name, setName] = useState("");
  const [user, setUser] = useState("");
  const [isRegistering, setIsRegistering] = useState(false);
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const handleLogin = async (e: any) => {
    e.preventDefault();

    const result = await loginUser(email, password);

    if (result.success) {
      localStorage.setItem("token", result.token);
      localStorage.setItem("username", result.user)
      setToken(result.token);
      setUser(result.user);
    } else {
      alert("Error: " + result.error);
    }
  };

  const toggleShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleRegister = async (e: any) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      alert('Las contraseñas no coinciden.');
      return;
    }

    const result = await registerUser({ email, password, name, phone });

    if (result.success) {
      alert("Registro exitoso, ahora inicia sesión.");
      setIsRegistering(false);
    } else {
      alert("Error: " + result.error);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.loginBox}>
        {!token ? (
          isRegistering ? (
            <form onSubmit={handleRegister}>
              <h2>Registro</h2>
              <div className={styles.inputGroup}>
                <label htmlFor="register-name">Nombre</label>
                <input
                  type="text"
                  id="register-name"
                  placeholder="Nombre"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className={styles.input}
                  required
                />
              </div>
              <div className={styles.inputGroup}>
                <label htmlFor="register-phone">Teléfono</label>
                <input
                  type="tel"
                  id="register-phone"
                  placeholder="Teléfono"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  className={styles.input}
                  required
                />
              </div>
              <div className={styles.inputGroup}>
                <label htmlFor="register-email">Email</label>
                <input
                  type="email"
                  id="register-email"
                  placeholder="Email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className={styles.input}
                  required
                />
              </div>
              <div className={styles.inputGroup}>
                <label htmlFor="register-password">Contraseña</label>
                <div className={styles.passwordInputContainer}>
                  <input
                    type={showPassword ? 'text' : 'password'}
                    id="register-password"
                    placeholder="Contraseña"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className={styles.passwordInput}
                    required
                  />
                  <button
                    type="button"
                    className={styles.passwordToggleButton}
                    onClick={toggleShowPassword}
                  >
                    <FontAwesomeIcon icon={showPassword ? faEye : faEyeSlash} />
                  </button>
                </div>
              </div>
              <div className={styles.inputGroup}>
                <label htmlFor="register-confirm-password">Confirmar contraseña</label>
                <div className={styles.passwordInputContainer}>
                  <input
                    type="password"
                    id="register-confirm-password"
                    placeholder="Confirmar contraseña"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    className={styles.passwordInput}
                    required
                  />
                </div>
              </div>
              <button type="submit" className={styles.loginButton}>
                Registrarse
              </button>
              <p className={styles.signUpText}>
                ¿Ya tienes cuenta?{' '}
                <button
                  type="button"
                  onClick={() => setIsRegistering(false)}
                  className={styles.signUpLinkButton}
                >
                  Iniciar sesión
                </button>
              </p>
            </form>
          ) : (
            <form onSubmit={handleLogin}>
              <h2>Iniciar Sesión</h2>
              <div className={styles.inputGroup}>
                <label htmlFor="login-email">Email</label>
                <input
                  type="email"
                  id="login-email"
                  placeholder="Email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className={styles.input}
                  required
                />
              </div>
              <div className={styles.inputGroup}>
                <label htmlFor="login-password">Contraseña</label>
                <div className={styles.passwordInputContainer}>
                  <input
                    type={showPassword ? 'text' : 'password'}
                    id="login-password"
                    placeholder="Contraseña"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className={styles.passwordInput}
                    required
                  />
                  <button
                    type="button"
                    className={styles.passwordToggleButton}
                    onClick={toggleShowPassword}
                  >
                    <FontAwesomeIcon icon={showPassword ? faEye : faEyeSlash} />
                  </button>
                </div>
              </div>
              <button type="submit" className={styles.loginButton}>
                Login
              </button>
              <p className={styles.signUpText}>
                ¿No tienes cuenta?{' '}
                <button
                  type="button"
                  onClick={() => setIsRegistering(true)}
                  className={styles.signUpLinkButton}
                >
                  Regístrate
                </button>
              </p>
            </form>
          )
        ) : (
          <div>
            <p>Redireccionando ... </p>
          </div>
        )}
      </div>
    </div>
  );
}

export default Login;