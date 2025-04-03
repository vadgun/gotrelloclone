import { useState } from "react";
import Boards from "./Boards";

function Login(){

      const [email, setEmail] = useState("");
      const [password, setPassword] = useState("");
      const [phone, setPhone] = useState("");
      const [name, setName] = useState("");
      const [user, setUser] = useState("");
      const [token, setToken] = useState(localStorage.getItem("token") || "");
      const [isRegistering, setIsRegistering] = useState(false);
    
      const handleLogin = async (e:any) => {
        e.preventDefault();
        const response = await fetch("http://localhost:8080/users/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email, password }),
        });
    
        const data = await response.json();
        if (response.ok) {
          localStorage.setItem("token", data.token);
          setToken(data.token);
          setUser(data.user);
        } else {
          alert("Error: " + data.error);
        }
      };

      const handleRegister = async (e:any) => {
        e.preventDefault();
        const response = await fetch("http://localhost:8080/users/register", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email, password, name, phone }),
        });
    
        if (response.ok) {
          alert("Registro exitoso, ahora inicia sesión.");
          setIsRegistering(false); // Cambiar a la pantalla de Login
        } else {
          const data = await response.json();
          alert("Error: " + data.message);
        }
      };
    
      const handleLogout = () => {
        localStorage.removeItem("token");
        setToken("");
        setEmail("");
        setUser("");
        setPhone("");
        setName("");
        setPassword("");
      };

      return (
        <div>
          {!token ? (
            isRegistering ? (
              <form onSubmit={handleRegister}>
                <h2>Registro</h2>
                <input
                  type="text"
                  placeholder="Nombre"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                />
                <input
                  type="tel"
                  placeholder="Telefono"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  pattern="[0-9]{3}-[0-9]{3}-[0-9]{4}"
                  required
                />
                <input
                  type="email"
                  placeholder="Email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
                <input
                  type="password"
                  placeholder="Contraseña"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
                <button type="submit">Registrarse</button>
                <p>
                  ¿Ya tienes cuenta?{" "}
                  <button type="button" onClick={() => setIsRegistering(false)}>
                    Iniciar sesión
                  </button>
                </p>
              </form>
            ) : (
              <form onSubmit={handleLogin}>
                <h2>Iniciar Sesión</h2>
                <input
                  type="email"
                  placeholder="Email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
                <input
                  type="password"
                  placeholder="Contraseña"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
                <button type="submit">Login</button>
                <p>
                  ¿No tienes cuenta?{" "}
                  <button type="button" onClick={() => setIsRegistering(true)}>
                    Regístrate
                  </button>
                </p>
              </form>
            )
          ) : (
            <div>
            <h2>¡Bienvenido! {user}</h2>
              <Boards></Boards>
            <button onClick={handleLogout}>Cerrar Sesión</button>
          </div>
          )}
        </div>
      );
}

export default Login;