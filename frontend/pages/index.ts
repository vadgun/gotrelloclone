import { useEffect, useState } from "react";
import axios from "axios";
import socket from "../utils/socket";

type Task = {
  id: string;
  title: string;
  status: "todo" | "in-progress" | "done";
};

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);

  useEffect(() => {
    // Cargar tareas al inicio
    axios.get("http://localhost:8082/tasks").then((res) => setTasks(res.data));

    // Escuchar eventos en tiempo real
    socket.on("task-updated", (updatedTask: Task) => {
      setTasks((prev) =>
        prev.map((task) => (task.id === updatedTask.id ? updatedTask : task))
      );
    });

    return () => {
      socket.off("task-updated");
    };
  }, []); 

  return(
    tasks
  )

}