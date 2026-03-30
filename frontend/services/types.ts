export type Task = {
  id: string;
  title: string;
  description: string;
  schedule: string;
  status: "pending" | "running" | "success" | "failed";
  created_at: string;
  updated_at: string;
  next_run_at: string | null;
  last_run_at: string | null;
  max_retries: number;
  retry_count: number;
  retry_delay_seconds: number;
  priority: string;
  executions?: TaskExecution[];
};

export type TaskUpdateEvent = {
  task_id: string;
  execution_id: string;
  status: "pending" | "running" | "success" | "failed";
  updated_at: string;
  next_run_at?: string | null;
  error_message?: string | null;
  retry_count?: number;
  max_retries?: number;
};

export type TaskExecution = {
  id: string;
  task_id: string;
  status: "running" | "success" | "failed";
  started_at: string;
  finished_at: string;
  error_message: string;
  created_at: string;
  updated_at: string;
  logs?: TaskLog;
};

export type TaskLog = {
  id: string;
  execution_id: string;
  level: string;
  message: string;
  created_at: string;
};

export type NewTask = {
  title: string;
  schedule: string;
  priority: string;
  description?: string;
};

export type User = {
  user_id: string;
  email: string;
  role: string;
  avatar_url?: string;
};

export type Login = {
  email: string;
  password: string;
};

export type Register = {
  name: string;
} & Login;
