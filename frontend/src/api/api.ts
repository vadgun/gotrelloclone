const USER_API_URL = "http://localhost:8080"; // Endpoint de user-service
const BOARD_API_URL = "http://localhost:8081"; // Endpoint de board-service
const TASK_API_URL = "http://localhost:8082"; // Endpoint de task-service

export const loginUser = async (email: string, password: string) => {
  const response = await fetch(`${USER_API_URL}/users/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json();

  if (response.ok) {
    return { success: true, token: data.token, user: data.user, role: data.role };
  } else {
    return { success: false, error: data.error };
  }
};

export const registerUser = async (userData: { email: string, password: string, name: string, phone: string }) => {
  const response = await fetch(`${USER_API_URL}/users/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(userData),
  });

  if (response.ok) {
    return { success: true };
  } else {
    const data = await response.json();
    console.log(response)
    return { success: false, error: data.message };
  }
};

export const createTask = async (newTask: { title: string, description: string, board_id: any, token: any }) => {
  const response = await fetch(`${TASK_API_URL}/tasks`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${newTask.token}`,
    },
    body: JSON.stringify(newTask),
  });

  const data = await response.json();

  if (response.ok) {
    return { success: true, task: data.task, message: data.message };
  } else {
    const data = await response.json();
    return { success: false, error: data.message };
  }
};

export const getBoards = async (token: any) => {
  const response = await fetch(`${BOARD_API_URL}/boards`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  const data = await response.json();
  return data.boards || [];
};

export const getTasks = async ({
  boardID,
  token,
  page,
  limit
}: {
  boardID: any;
  token: string;
  page: number;
  limit: number;
}) => {
  const response = await fetch(
    `${TASK_API_URL}/tasks/board/${boardID}?page=${page}&limit=${limit}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  const data = await response.json();

  return {
    tasks: data.tasks || [],
    totalPages: data.totalPages || 1,
  };
};

export const getAdminUsers = async (token: string) => {
  const res = await fetch(`${USER_API_URL}/admin/users`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.json();
};

export const getAdminBoards = async (token: string) => {
  const res = await fetch(`${BOARD_API_URL}/admin/boards`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.json();
};

export const getAdminTasks = async (token: string) => {
  const res = await fetch(`${TASK_API_URL}/admin/tasks`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.json();
};