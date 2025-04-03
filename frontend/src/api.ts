const API_URL = "http://localhost:8081"; // Endpoint de board-service

export const getBoards = async () => {
  const token = localStorage.getItem("token");
  const response = await fetch(`${API_URL}/boards`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  
  const data = await response.json();
  return data.boards || [];
};

