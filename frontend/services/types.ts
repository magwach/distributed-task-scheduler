export type Task = {
  id: string;
  title: string;
  description: string;
  schedule: string;
  status: "pending" | "running" | "success" | "failed";
  created_at: string;
  updated_at: string;
  next_run_at?: string;
  last_run_at?: string;
  executions?: TaskExecution[];
  task_id?: string;
  execution_id?: string;
  error_message?: string;
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
};

export type NewTask = {
  title: string;
  schedule: string;
  description?: string;
};

