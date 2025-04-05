import { useState } from "react";
import { loginUser, registerUser } from "./api";


function Login({ token, setToken }: { token: string; setToken: (token: string) => void }) {

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [phone, setPhone] = useState("");
  const [name, setName] = useState("");
  const [user, setUser] = useState("");
  const [isRegistering, setIsRegistering] = useState(false);

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

  const handleRegister = async (e: any) => {
    e.preventDefault();

    const result = await registerUser({ email, password, name, phone });

    if (result.success) {
      alert("Registro exitoso, ahora inicia sesión.");
      setIsRegistering(false);
    } else {
      alert("Error: " + result.error);
    }
  };

  return (
    <div>
      {!token ? (
        isRegistering ? (
          <form onSubmit={handleRegister}>
            <h2>Registro</h2>
            <input type="text" placeholder="Nombre" value={name} onChange={(e) => setName(e.target.value)} required />
            <input type="tel" placeholder="Teléfono" value={phone} onChange={(e) => setPhone(e.target.value)} required />
            <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} required />
            <input type="password" placeholder="Contraseña" value={password} onChange={(e) => setPassword(e.target.value)} required />
            <button type="submit">Registrarse</button>
            <p>
              ¿Ya tienes cuenta? <button type="button" onClick={() => setIsRegistering(false)}>Iniciar sesión</button>
            </p>
          </form>
        ) : (
          <form onSubmit={handleLogin}>
            <h2>Iniciar Sesión</h2>
            <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} required />
            <input type="password" placeholder="Contraseña" value={password} onChange={(e) => setPassword(e.target.value)} required />
            <button type="submit">Login</button>
            <p>
              ¿No tienes cuenta? <button type="button" onClick={() => setIsRegistering(true)}>Regístrate</button>
            </p>
          </form>
        )
      ) : (
        <div>
          <p>Redireccionando ... </p>
      </div>
      )}
    </div>
  );
}

export default Login;