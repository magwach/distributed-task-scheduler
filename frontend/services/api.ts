import { NewTask, Task } from "./types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getTasks(): Promise<Task[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/tasks`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!res.ok) {
      throw new Error("Failed to fetch tasks");
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
    });

    if (!res.ok) {
      throw new Error("Failed to fetch task");
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
    });

    if (!res.ok) {
      throw new Error("Failed to delete task");
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
      body: JSON.stringify(task),
    });

    if (!res.ok) {
      throw new Error("Failed to add task");
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
    });

    if (!res.ok) {
      throw new Error("Failed to retry task");
    }

    const data = await res.json();

    return "success";
  } catch (error) {
    console.error("Failed to retry tasks");
    return "failed";
  }
}
