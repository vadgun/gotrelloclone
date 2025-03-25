import { io } from "socket.io-client";

const socket = io("http://localhost:8083"); // Conectamos con notification-service

export default socket;
