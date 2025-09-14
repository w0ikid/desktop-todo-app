interface Task {
    ID: string;
    Title: string;
    Status: string;
    Priority: string;
    CreatedAt: string;
    DueDate?: string;
    }

interface Dashboard {
    active_count: number;
    completed_count: number;
    overdue_count: number;
    due_today: Task[];
    due_this_week: Task[];
    recent_tasks: Task[];
}