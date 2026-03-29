import { Login, NewTask, Register, Task, User } from "./types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getTasks(): Promise<Task[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/tasks`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });
    if (!res.ok) {
      throw new Error("Failed to fetch tasks");
    }

    if (res.status === 401) {
      window.location.href = "/login";
      throw new Error("Unauthorized");
    }

    const data = await res.json();

    return data.data;
  } catch (error) {
    console.error("Failed to fetch tasks");
    return [];
  }
}

export async function getTask(id: string): Promise<Task> {
  try {
    const res = await fetch(`${API_BASE_URL}/task/${id}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });

    if (!res.ok) {
      throw new Error("Failed to fetch task");
    }

    if (res.status === 401) {
      window.location.href = "/login";
      throw new Error("Unauthorized");
    }

    const data = await res.json();

    return data.data;
  } catch (error) {
    throw new Error("Failed to fetch tasks");
  }
}

export async function deleteTask(id: string): Promise<Task[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/task/${id}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });

    if (!res.ok) {
      throw new Error("Failed to delete task");
    }

    if (res.status === 401) {
      window.location.href = "/login";
      throw new Error("Unauthorized");
    }

    const data = await res.json();

    return data.data;
  } catch (error) {
    console.error("Failed to delete tasks");
    return [];
  }
}

export async function createTask(task: NewTask): Promise<Task[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/task`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(task),
    });

    if (!res.ok) {
      throw new Error("Failed to add task");
    }

    if (res.status === 401) {
      localStorage.removeItem("auth_token");
      throw new Error("Unauthorized");
    }

    const data = await res.json();

    return data.data;
  } catch (error) {
    console.error("Failed to add tasks");
    return [];
  }
}

export async function retryTask(taskId: String): Promise<string> {
  try {
    const res = await fetch(`${API_BASE_URL}/task/${taskId}/retry`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });

    if (!res.ok) {
      throw new Error("Failed to retry task");
    }
    if (res.status === 401) {
      window.location.href = "/login";
      throw new Error("Unauthorized");
    }

    await res.json();

    return "success";
  } catch (error) {
    console.error("Failed to retry tasks");
    return "failed";
  }
}

export async function fetchUser(): Promise<User> {
  try {
    const res = await fetch(`${API_BASE_URL}/me`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });

    if (!res.ok) {
      throw new Error("Failed to get user");
    }
    if (res.status === 401) {
      window.location.href = "/login";
      throw new Error("Unauthorized");
    }

    const data = await res.json();

    return data;
  } catch (error) {
    console.error("Failed to fetch user");
    throw new Error("Failed to get user");
  }
}

export async function login(form: Login): Promise<User> {
  try {
    const res = await fetch(`${API_BASE_URL}/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(form),
    });

    if (!res.ok) {
      throw new Error("Failed to login");
    }

    const data = await res.json();

    return data;
  } catch (error) {
    console.error("Failed to login");
    throw new Error("Failed to login");
  }
}

export async function register(form: Register): Promise<User> {
  try {
    const res = await fetch(`${API_BASE_URL}/auth/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(form),
    });

    if (!res.ok) {
      throw new Error("Failed to create account.");
    }

    const data = await res.json();

    return data;
  } catch (error) {
    console.error("Failed to create account.");
    throw new Error("Failed to create account.");
  }
}
