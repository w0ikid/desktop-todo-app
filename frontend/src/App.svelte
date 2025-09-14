<script lang="ts">
  import { onMount } from "svelte";
  import * as TaskHandler from "../wailsjs/go/wails/TaskHandler";

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

  // State
  let tasks: Task[] = [];
  let dashboard: Dashboard | null = null;
  let currentView = "dashboard";
  let loading = false;
  let error = "";

  // Form state
  let newTitle = "";
  let newPriority = "medium";
  let newDueDate = "";
  let editingTask: Task | null = null;

  // Load dashboard data
  async function loadDashboard() {
    try {
      loading = true;
      error = "";
      const result = await TaskHandler.GetDashboard();
      dashboard = result;
    } catch (err) {
      error = `Error loading dashboard: ${err}`;
      console.error(err);
    } finally {
      loading = false;
    }
  }

  // Load tasks based on filter
  async function loadTasks(filter?: string, status?: string) {
    try {
      loading = true;
      error = "";
      const result = await TaskHandler.ListTasks(status || null, null, filter || null);
      tasks = result.tasks || [];
    } catch (err) {
      error = `Error loading tasks: ${err}`;
      console.error(err);
    } finally {
      loading = false;
    }
  }

  // Add new task
  async function addTask() {
    if (!newTitle.trim()) {
      error = "Task title is required";
      return;
    }

    try {
      loading = true;
      error = "";
      const dueDate = newDueDate ? new Date(newDueDate) : null;
      await TaskHandler.CreateTask(newTitle.trim(), newPriority, dueDate);
      
      // Reset form
      newTitle = "";
      newPriority = "medium";
      newDueDate = "";
      
      // Refresh current view
      await refreshCurrentView();
    } catch (err) {
      error = `Error creating task: ${err}`;
      console.error(err);
    } finally {
      loading = false;
    }
  }

  // Complete task
  async function completeTask(id: string) {
    try {
      loading = true;
      await TaskHandler.CompleteTask(id);
      await refreshCurrentView();
    } catch (err) {
      error = `Error completing task: ${err}`;
      console.error(err);
    } finally {
      loading = false;
    }
  }

  // Delete task
  async function deleteTask(id: string) {
    if (!confirm("Are you sure you want to delete this task?")) return;
    
    try {
      loading = true;
      await TaskHandler.DeleteTask(id);
      await refreshCurrentView();
    } catch (err) {
      error = `Error deleting task: ${err}`;
      console.error(err);
    } finally {
      loading = false;
    }
  }

  // Update task
  async function updateTask(task: Task) {
    try {
      loading = true;
      const dueDate = task.DueDate ? new Date(task.DueDate) : null;
      await TaskHandler.UpdateTask(
        task.ID,
        task.Title || null,
        task.Status || null,
        task.Priority || null,
        dueDate
      );
      editingTask = null;
      await refreshCurrentView();
    } catch (err) {
      error = `Error updating task: ${err}`;
      console.error(err);
    } finally {
      loading = false;
    }
  }

  // Switch view and load appropriate data
  async function switchView(view: string) {
    currentView = view;
    await refreshCurrentView();
  }

  async function refreshCurrentView() {
    switch (currentView) {
      case "dashboard":
        await loadDashboard();
        break;
      case "all":
        await loadTasks();
        break;
      case "today":
        await loadTasks("today");
        break;
      case "week":
        await loadTasks("week");
        break;
      case "overdue":
        await loadTasks("overdue");
        break;
      case "completed":
        await loadTasks(null, "completed");
        break;
    }
  }

  function getPriorityColor(priority: string): string {
    switch (priority) {
      case "high": return "priority-high";
      case "medium": return "priority-medium";
      case "low": return "priority-low";
      default: return "priority-default";
    }
  }

  function getStatusColor(status: string): string {
    switch (status) {
      case "completed": return "status-completed";
      case "active": return "status-active";
      default: return "status-default";
    }
  }

  function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric"
    });
  }

  function isOverdue(task: Task): boolean {
    if (!task.DueDate || task.Status === "completed") return false;
    return new Date(task.DueDate) < new Date();
  }

  onMount(() => {
    loadDashboard();
  });
</script>

<main class="app-container">
  <!-- Animated Background -->
  <!-- <div class="background-animation"></div> -->
  
  <div class="main-content">
    <!-- Header with Gradient -->
    <header class="app-header">
      <div class="header-content">
        <h1 class="app-title">
          <span class="title-icon">✨</span>
          Todo Dashboard
        </h1>
        <p class="app-subtitle">Manage your tasks with style and efficiency</p>
      </div>
      <div class="header-decoration"></div>
    </header>

    <!-- Error Alert with Animation -->
    {#if error}
      <div class="error-alert animate-slide-down">
        <div class="error-content">
          <span class="error-icon">⚠️</span>
          <span class="error-text">{error}</span>
          <button on:click={() => error = ""} class="error-close">×</button>
        </div>
      </div>
    {/if}

    <!-- Navigation with Hover Effects -->
    <nav class="nav-container">
      <div class="nav-wrapper">
        {#each [
          { id: "dashboard", label: "Dashboard", icon: "📊" },
          { id: "all", label: "All Tasks", icon: "📝" },
          { id: "today", label: "Due Today", icon: "📅" },
          { id: "week", label: "This Week", icon: "🗓️" },
          { id: "overdue", label: "Overdue", icon: "⏰" },
          { id: "completed", label: "Completed", icon: "✅" }
        ] as navItem}
          <button
            on:click={() => switchView(navItem.id)}
            class="nav-button {currentView === navItem.id ? 'nav-active' : ''}"
          >
            <span class="nav-icon">{navItem.icon}</span>
            <span class="nav-label">{navItem.label}</span>
          </button>
        {/each}
      </div>
    </nav>

    <!-- Add Task Form with Glassmorphism -->
    <div class="add-task-card">
      <h2 class="form-title">
        <span class="form-icon">➕</span>
        Add New Task
      </h2>
      <div class="form-grid">
        <div class="input-group">
          <input
            type="text"
            placeholder="What needs to be done?"
            bind:value={newTitle}
            class="form-input"
          />
          <div class="input-border"></div>
        </div>
        <div class="input-group">
          <select bind:value={newPriority} class="form-select">
            <option value="low">🟢 Low </option>
            <option value="medium">🟡 Medium</option>
            <option value="high">🔴 High</option>
          </select>
          <div class="input-border"></div>
        </div>
        <div class="input-group">
          <input
            type="date"
            bind:value={newDueDate}
            class="form-input"
          />
          <div class="input-border"></div>
        </div>
        <button
          on:click={addTask}
          disabled={loading || !newTitle.trim()}
          class="add-button {loading ? 'loading' : ''}"
        >
          <span class="button-text">{loading ? "Adding..." : "Add Task"}</span>
          <span class="button-icon">🚀</span>
        </button>
      </div>
    </div>

    <!-- Dashboard View -->
    {#if currentView === "dashboard" && dashboard}
      <div class="dashboard-grid">
        <!-- Stats Cards with Animation -->
        <div class="stat-card stat-active animate-float">
          <div class="stat-icon">📋</div>
          <div class="stat-content">
            <h3 class="stat-title">Active Tasks</h3>
            <p class="stat-number">{dashboard.active_count}</p>
          </div>
        </div>
        <div class="stat-card stat-completed animate-float" style="animation-delay: 0.1s;">
          <div class="stat-icon">✅</div>
          <div class="stat-content">
            <h3 class="stat-title">Completed</h3>
            <p class="stat-number">{dashboard.completed_count}</p>
          </div>
        </div>
        <div class="stat-card stat-overdue animate-float" style="animation-delay: 0.2s;">
          <div class="stat-icon">⚡</div>
          <div class="stat-content">
            <h3 class="stat-title">Overdue</h3>
            <p class="stat-number">{dashboard.overdue_count}</p>
          </div>
        </div>
      </div>

      <!-- Dashboard Sections -->
      <div class="dashboard-sections">
        <!-- Due Today -->
        <div class="section-card">
          <h3 class="section-title">
            <span class="section-icon">🌅</span>
            Due Today
          </h3>
          {#if dashboard.due_today?.length}
            <div class="task-list">
              {#each dashboard.due_today as task, index}
                <div class="task-item animate-slide-up" style="animation-delay: {index * 0.1}s;">
                  <div class="task-content">
                    <p class="task-title">{task.Title}</p>
                    <div class="task-meta">
                      <span class="task-priority {getPriorityColor(task.Priority)}">
                        {task.Priority}
                      </span>
                    </div>
                  </div>
                  <button
                    on:click={() => completeTask(task.ID)}
                    class="complete-button"
                  >
                    Complete
                  </button>
                </div>
              {/each}
            </div>
          {:else}
            <div class="empty-state">
              <div class="empty-icon">🎉</div>
              <p class="empty-text">No tasks due today</p>
            </div>
          {/if}
        </div>

        <!-- Recent Tasks -->
        <div class="section-card">
          <h3 class="section-title">
            <span class="section-icon">⚡</span>
            Recent Tasks
          </h3>
          {#if dashboard.recent_tasks?.length}
            <div class="task-list">
              {#each dashboard.recent_tasks as task, index}
                <div class="task-item animate-slide-up" style="animation-delay: {index * 0.1}s;">
                  <div class="task-content">
                    <p class="task-title">{task.Title}</p>
                    <div class="task-meta">
                      <span class="task-priority {getPriorityColor(task.Priority)}">
                        {task.Priority}
                      </span>
                      {#if task.DueDate}
                        <span class="task-date">Due: {formatDate(task.DueDate)}</span>
                      {/if}
                    </div>
                  </div>
                  <button
                    on:click={() => completeTask(task.ID)}
                    class="complete-button"
                  >
                    Complete
                  </button>
                </div>
              {/each}
            </div>
          {:else}
            <div class="empty-state">
              <div class="empty-icon">📝</div>
              <p class="empty-text">No recent tasks</p>
            </div>
          {/if}
        </div>
      </div>
    {/if}

    <!-- Task List View -->
    {#if currentView !== "dashboard"}
      <div class="task-list-card">
        <div class="task-list-header">
          <h2 class="list-title">
            {#if currentView === "all"}📝 All Tasks
            {:else if currentView === "today"}📅 Due Today
            {:else if currentView === "week"}🗓️ Due This Week
            {:else if currentView === "overdue"}⏰ Overdue Tasks
            {:else if currentView === "completed"}✅ Completed Tasks
            {/if}
            {#if tasks.length > 0}
              <span class="task-count">({tasks.length})</span>
            {/if}
          </h2>
        </div>

        {#if loading}
          <div class="loading-state">
            <div class="loading-spinner"></div>
            <p class="loading-text">Loading tasks...</p>
          </div>
        {:else if tasks.length === 0}
          <div class="empty-state">
            <div class="empty-icon">🎯</div>
            <p class="empty-text">No tasks found</p>
          </div>
        {:else}
          <div class="tasks-container">
            {#each tasks as task, index (task.ID)}
              <div class="task-card animate-slide-up {isOverdue(task) ? 'overdue' : ''}" style="animation-delay: {index * 0.05}s;">
                {#if editingTask?.ID === task.ID}
                  <!-- Edit Mode -->
                  <div class="edit-form">
                    <div class="edit-input-group">
                      <input
                        bind:value={editingTask.Title}
                        class="edit-input"
                        placeholder="Task title..."
                      />
                    </div>
                    <div class="edit-controls">
                      <select bind:value={editingTask.Priority} class="edit-select">
                        <option value="low">🟢 Low</option>
                        <option value="medium">🟡 Medium</option>
                        <option value="high">🔴 High</option>
                      </select>
                      <input
                        type="date"
                        bind:value={editingTask.DueDate}
                        class="edit-date"
                      />
                      <button
                        on:click={() => updateTask(editingTask)}
                        class="save-button"
                      >
                        Save
                      </button>
                      <button
                        on:click={() => editingTask = null}
                        class="cancel-button"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                {:else}
                  <!-- View Mode -->
                  <div class="task-view">
                    <div class="task-main">
                      <h3 class="task-title">{task.Title}</h3>
                      <div class="task-badges">
                        <span class="task-status {getStatusColor(task.Status)}">
                          {task.Status}
                        </span>
                        <span class="task-priority {getPriorityColor(task.Priority)}">
                          {task.Priority}
                        </span>
                        <span class="task-date">Created: {formatDate(task.CreatedAt)}</span>
                        {#if task.DueDate}
                          <span class="task-date {isOverdue(task) ? 'overdue-badge' : ''}">
                            Due: {formatDate(task.DueDate)}
                            {#if isOverdue(task)}
                              (Overdue)
                            {/if}
                          </span>
                        {/if}
                      </div>
                    </div>
                    <div class="task-actions">
                      <button
                        on:click={() => editingTask = {...task}}
                        class="action-button edit-btn"
                      >
                        ✏️ Edit
                      </button>
                      {#if task.Status !== "completed"}
                        <button
                          on:click={() => completeTask(task.ID)}
                          class="action-button complete-btn"
                        >
                          ✅ Complete
                        </button>
                      {/if}
                      <button
                        on:click={() => deleteTask(task.ID)}
                        class="action-button delete-btn"
                      >
                        🗑️ Delete
                      </button>
                    </div>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        {/if}
      </div>
    {/if}
  </div>
</main>
<style>
  :global(html, body) {
    margin: 0;
    padding: 0;
    background-color: var(--background-color);
  }
  /* Modern CSS Variables */
  :global(:root) {
    --primary-color: #6a5acd; /* SlateBlue */
    --secondary-color: #f08080; /* LightCoral */
    --background-color: #f4f7fa;
    --surface-color: #ffffff;
    --text-color: #333333;
    --muted-text-color: #777777;
    --border-color: #e0e0e0;
    --shadow-color: rgba(0, 0, 0, 0.1);
    --success-color: #2ecc71;
    --error-color: #e74c3c;
    --warning-color: #f39c12;

    --font-family-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji",
      "Segoe UI Symbol";
    
    --border-radius: 12px;
    --transition-speed: 0.3s;
  }

  /* Global Styles */
  :global(body) {
    margin: 0;
    font-family: var(--font-family-sans);
    background-color: var(--background-color);
    color: var(--text-color);
    line-height: 1.6;
  }

  /* App Container */
  .app-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  /* Header */
  .app-header {
    text-align: center;
    margin-bottom: 3rem;
  }

  .app-title {
    font-size: 2.5rem;
    font-weight: 700;
    color: var(--primary-color);
    margin: 0 0 0.5rem 0;
  }

  .app-subtitle {
    font-size: 1.1rem;
    color: var(--muted-text-color);
    margin: 0;
  }

  /* Error Alert */
  .error-alert {
    background-color: var(--error-color);
    color: white;
    padding: 1rem;
    border-radius: var(--border-radius);
    margin-bottom: 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 4px 15px rgba(231, 76, 60, 0.3);
  }

  .error-close {
    background: none;
    border: none;
    color: white;
    font-size: 1.5rem;
    cursor: pointer;
  }

  /* Navigation */
  .nav-container {
    margin-bottom: 2rem;
  }
  
  .nav-wrapper {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    background-color: var(--surface-color);
    border-radius: var(--border-radius);
    padding: 0.5rem;
    box-shadow: 0 2px 10px var(--shadow-color);
  }

  .nav-button {
    flex: 1;
    padding: 0.75rem 1rem;
    border: none;
    border-radius: 8px;
    background-color: transparent;
    color: var(--muted-text-color);
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
    text-align: center;
  }
  
  .nav-button:hover {
    background-color: rgba(0,0,0,0.05);
    color: var(--primary-color);
  }
  
  .nav-active {
    background-color: var(--primary-color);
    color: white !important;
    box-shadow: 0 4px 10px rgba(106, 90, 205, 0.3);
  }

  /* Add Task Form */
  .add-task-card {
    background-color: var(--surface-color);
    border-radius: var(--border-radius);
    padding: 2rem;
    margin-bottom: 2rem;
    box-shadow: 0 4px 20px var(--shadow-color);
  }

  .form-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0 0 1.5rem 0;
  }

  .form-grid {
    display: grid;
    grid-template-columns: 1fr auto auto auto;
    gap: 1rem;
    align-items: center;
  }

  @media (max-width: 768px) {
    .form-grid {
      grid-template-columns: 1fr;
    }
  }

  .form-input, .form-select {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    font-size: 1rem;
    transition: all var(--transition-speed) ease;
  }

  .form-input:focus, .form-select:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(106, 90, 205, 0.2);
  }
  
  .add-button {
    padding: 0.75rem 1.5rem;
    background-color: var(--primary-color);
    border: none;
    border-radius: 8px;
    color: white;
    font-weight: 500;
    font-size: 1rem;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
  }

  .add-button:hover:not(:disabled) {
    background-color: #5a4acd;
    transform: translateY(-2px);
    box-shadow: 0 6px 15px rgba(106, 90, 205, 0.4);
  }
  
  .add-button:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }

  /* Dashboard Stats */
  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
  }

  .stat-card {
    background-color: var(--surface-color);
    border-radius: var(--border-radius);
    padding: 1.5rem;
    box-shadow: 0 4px 20px var(--shadow-color);
    text-align: center;
  }

  .stat-title {
    font-size: 1rem;
    color: var(--muted-text-color);
    margin: 0 0 0.5rem 0;
  }

  .stat-number {
    font-size: 2rem;
    font-weight: 700;
    color: var(--primary-color);
    margin: 0;
  }

  /* Dashboard Sections */
  .dashboard-sections {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }

  @media (max-width: 768px) {
    .dashboard-sections {
      grid-template-columns: 1fr;
    }
  }

  .section-card {
    background-color: var(--surface-color);
    border-radius: var(--border-radius);
    padding: 2rem;
    box-shadow: 0 4px 20px var(--shadow-color);
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0 0 1.5rem 0;
  }
  
  /* Task List & Items */
  .task-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .task-item, .task-card {
    background-color: #fdfdff;
    border-radius: 8px;
    padding: 1rem 1.5rem;
    border: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: all var(--transition-speed) ease;
  }

  .task-card {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .task-item:hover, .task-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 15px var(--shadow-color);
    border-color: var(--primary-color);
  }

  .task-title {
    font-weight: 500;
    margin: 0;
    font-size: 1.1rem;
  }

  .task-meta {
    display: flex;
    gap: 1rem;
    align-items: center;
    font-size: 0.9rem;
    color: var(--muted-text-color);
  }
  
  .task-priority {
    padding: 0.25rem 0.75rem;
    border-radius: 1rem;
    font-size: 0.8rem;
    font-weight: 500;
    text-transform: capitalize;
  }

  .priority-high { background-color: #ffcccb; color: #c0392b; }
  .priority-medium { background-color: #fff3cd; color: #f39c12; }
  .priority-low { background-color: #d4edda; color: #1e7e34; }

  /* Task List View */
  .task-list-card {
    background-color: var(--surface-color);
    border-radius: var(--border-radius);
    padding: 2rem;
    box-shadow: 0 4px 20px var(--shadow-color);
  }

  .list-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0 0 1.5rem 0;
  }
  
  .task-badges {
    display: flex;
    gap: 0.5rem;
    margin-top: 0.75rem;
  }
  
  .task-status {
    padding: 0.25rem 0.75rem;
    border-radius: 1rem;
    font-size: 0.8rem;
    font-weight: 500;
    text-transform: capitalize;
  }

  .status-active { background-color: #cfe2ff; color: #0c5460; }
  .status-completed { background-color: #d4edda; color: #1e7e34; }

  .task-actions {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
  }

  .action-button {
    padding: 0.5rem 1rem;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    background-color: transparent;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
  }
  
  .action-button:hover {
    background-color: #f0f0f0;
  }
  
  .complete-btn:hover { background-color: var(--success-color); color: white; }
  .delete-btn:hover { background-color: var(--error-color); color: white; }
  
  .task-card.overdue {
    border-left: 4px solid var(--error-color);
  }
  
  /* Edit Form */
  .edit-form {
    width: 100%;
    padding-top: 1rem;
  }
  
  .edit-input, .edit-select, .edit-date {
    width: 100%;
    margin-bottom: 0.5rem;
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 8px;
  }

  .edit-controls {
    display: flex;
    gap: 0.5rem;
  }
  
  /* Empty & Loading States */
  .empty-state, .loading-state {
    text-align: center;
    padding: 3rem;
    color: var(--muted-text-color);
  }

  .loading-spinner {
    width: 30px;
    height: 30px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid var(--primary-color);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 1rem;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

</style>