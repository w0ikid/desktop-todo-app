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

  // Theme state
  let isDarkMode = false;

  // Form state
  let newTitle = "";
  let newPriority = "medium";
  let newDueDate = "";
  let editingTask: Task | null = null;

  // Form validation
  let titleError = "";

  // Theme functions
  function toggleTheme() {
    isDarkMode = !isDarkMode;
    document.documentElement.setAttribute('data-theme', isDarkMode ? 'dark' : 'light');
    
    // Save theme preference (in a real app, you might want to store this in backend)
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('theme', isDarkMode ? 'dark' : 'light');
    }
  }

  function initializeTheme() {
    // Check for saved theme preference or default to light mode
    let savedTheme = 'light';
    if (typeof localStorage !== 'undefined') {
      savedTheme = localStorage.getItem('theme') || 'light';
    }
    
    // Check system preference if no saved preference
    if (savedTheme === 'light' && window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      savedTheme = 'dark';
    }
    
    isDarkMode = savedTheme === 'dark';
    document.documentElement.setAttribute('data-theme', savedTheme);
  }

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

  // Validate form
  function validateForm(): boolean {
    titleError = "";
    
    if (!newTitle.trim()) {
      titleError = "Task title is required";
      return false;
    }
    
    if (newTitle.trim().length > 200) {
      titleError = "Task title is too long (max 200 characters)";
      return false;
    }
    
    return true;
  }

  // Add new task
  async function addTask() {
    if (!validateForm()) {
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
      titleError = "";
      
      // Refresh current view
      await refreshCurrentView();
    } catch (err) {
      error = `Error creating task: ${err}`;
      console.error(err);
    } finally {
      loading = false;
    }
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      
      if (!loading && newTitle.trim()) {
        addTask();
      }
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
    
    const dueDate = new Date(task.DueDate);
    const currentDate = new Date();
    
    dueDate.setHours(0, 0, 0, 0);
    currentDate.setHours(0, 0, 0, 0);
    
    return dueDate < currentDate;
  }

  onMount(() => {
    initializeTheme();
    loadDashboard();
  });
</script>

<main class="app-container">
  <div class="main-content">
    <!-- Header with Gradient and Theme Toggle -->
    <header class="app-header">
      <div class="header-content">
        <div class="header-top">
          <h1 class="app-title">
            <span class="title-icon">✨</span>
            Todo Dashboard
          </h1>
          <button 
            class="theme-toggle" 
            on:click={toggleTheme}
            aria-label="Toggle theme"
            title={isDarkMode ? 'Switch to light mode' : 'Switch to dark mode'}
          >
            <span class="theme-icon">
              {#if isDarkMode}
                ☀️
              {:else}
                🌙
              {/if}
            </span>
          </button>
        </div>
        <p class="app-subtitle">Manage your tasks with style and efficiency</p>
      </div>
    </header>

    <!-- Error Alert with Animation -->
    {#if error}
      <div class="error-alert animate-slide-down">
        <div class="error-content">
          <span class="error-icon">⚠️</span>
          <span class="error-text">{error}</span>
          <button on:click={() => error = ""} class="error-close" aria-label="Close error">×</button>
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
            aria-pressed={currentView === navItem.id}
          >
            <span class="nav-icon">{navItem.icon}</span>
            <span class="nav-label">{navItem.label}</span>
          </button>
        {/each}
      </div>
    </nav>

    <!-- Add Task Form with better layout -->
    <div class="add-task-card">
      <h2 class="form-title">
        <span class="form-icon">➕</span>
        Add New Task
      </h2>
      <form class="task-form" on:submit|preventDefault={addTask}>
        <div class="form-row">
          <div class="input-group title-input">
            <label for="task-title" class="sr-only">Task Title</label>
            <input
              id="task-title"
              type="text"
              placeholder="What needs to be done?"
              on:keydown={handleKeyPress}
              bind:value={newTitle}
              class="form-input {titleError ? 'input-error' : ''}"
              maxlength="200"
              required
            />
            {#if titleError}
              <span class="input-error-text">{titleError}</span>
            {/if}
          </div>
          <div class="input-group">
            <label for="task-priority" class="sr-only">Priority</label>
            <select id="task-priority" bind:value={newPriority} class="form-select">
              <option value="low">🟢 Low Priority</option>
              <option value="medium">🟡 Medium Priority</option>
              <option value="high">🔴 High Priority</option>
            </select>
          </div>
          <div class="input-group">
            <label for="task-due-date" class="sr-only">Due Date</label>
            <input
              id="task-due-date"
              type="date"
              bind:value={newDueDate}
              class="form-input"
              min={new Date().toISOString().split('T')[0]}
            />
          </div>
          <button
            type="submit"
            disabled={loading || !newTitle.trim()}
            class="add-button {loading ? 'loading' : ''}"
          >
            <span class="button-text">{loading ? "Adding..." : "Add Task"}</span>
            <span class="button-icon">🚀</span>
          </button>
        </div>
      </form>
    </div>

    <!-- Dashboard View -->
    {#if currentView === "dashboard" && dashboard}
      <div class="dashboard-grid">
        <!-- Stats Cards with Animation -->
        <div class="stat-card stat-active">
          <div class="stat-icon">📋</div>
          <div class="stat-content">
            <h3 class="stat-title">Active Tasks</h3>
            <p class="stat-number">{dashboard.active_count}</p>
          </div>
        </div>
        <div class="stat-card stat-completed">
          <div class="stat-icon">✅</div>
          <div class="stat-content">
            <h3 class="stat-title">Completed</h3>
            <p class="stat-number">{dashboard.completed_count}</p>
          </div>
        </div>
        <div class="stat-card stat-overdue">
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
              {#each dashboard.due_today as task}
                <div class="task-item">
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
                    disabled={loading}
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
              {#each dashboard.recent_tasks as task}
                <div class="task-item">
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
                    disabled={loading}
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
            {#each tasks as task (task.ID)}
              <div class="task-card {isOverdue(task) ? 'overdue' : ''}">
                {#if editingTask?.ID === task.ID}
                  <!-- Edit Mode -->
                  <div class="edit-form">
                    <div class="edit-input-group">
                      <input
                        bind:value={editingTask.Title}
                        class="edit-input"
                        placeholder="Task title..."
                        maxlength="200"
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
                        disabled={loading}
                      >
                        {loading ? "Saving..." : "Save"}
                      </button>
                      <button
                        on:click={() => editingTask = null}
                        class="cancel-button"
                        disabled={loading}
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
                        disabled={loading}
                      >
                        ✏️ Edit
                      </button>
                      {#if task.Status !== "completed"}
                        <button
                          on:click={() => completeTask(task.ID)}
                          class="action-button complete-btn"
                          disabled={loading}
                        >
                          ✅ Complete
                        </button>
                      {/if}
                      <button
                        on:click={() => deleteTask(task.ID)}
                        class="action-button delete-btn"
                        disabled={loading}
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
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  /* Light Theme Variables */
  :global(:root[data-theme="light"]) {
    --primary-color: #6a5acd;
    --secondary-color: #f08080;
    --background-color: #f4f7fa;
    --surface-color: #ffffff;
    --surface-hover: #fafaff;
    --text-color: #333333;
    --muted-text-color: #777777;
    --border-color: #e0e0e0;
    --shadow-color: rgba(0, 0, 0, 0.1);
    --success-color: #2ecc71;
    --error-color: #e74c3c;
    --warning-color: #f39c12;
    --card-gradient-start: #ffffff;
    --card-gradient-end: #fafaff;

    --font-family-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji",
      "Segoe UI Symbol";
    
    --border-radius: 12px;
    --transition-speed: 0.3s;
  }

  /* Dark Theme Variables */
  :global(:root[data-theme="dark"]) {
    --primary-color: #8a7cd8;
    --secondary-color: #ff9999;
    --background-color: #1a1a1a;
    --surface-color: #2d2d30;
    --surface-hover: #353538;
    --text-color: #e4e4e7;
    --muted-text-color: #a1a1aa;
    --border-color: #404040;
    --shadow-color: rgba(0, 0, 0, 0.3);
    --success-color: #22c55e;
    --error-color: #ef4444;
    --warning-color: #f59e0b;
    --card-gradient-start: #2d2d30;
    --card-gradient-end: #353538;

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
    transition: background-color var(--transition-speed) ease, color var(--transition-speed) ease;
  }

  /* Screen Reader Only */
  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
  }

  /* App Container - Improved responsiveness */
  .app-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 1rem;
  }

  @media (min-width: 768px) {
    .app-container {
      padding: 2rem;
    }
  }

  /* Header with Theme Toggle */
  .app-header {
    text-align: center;
    margin-bottom: 3rem;
  }

  .header-top {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 2rem;
    margin-bottom: 0.5rem;
    position: relative;
  }

  .app-title {
    font-size: 2.5rem;
    font-weight: 700;
    color: var(--primary-color);
    margin: 0;
  }

  .title-icon {
    margin-right: 0.5rem;
  }

  .theme-toggle {
    position: absolute;
    right: 0;
    background: linear-gradient(135deg, var(--surface-color), var(--surface-hover));
    border: 2px solid var(--border-color);
    border-radius: 50%;
    width: 50px;
    height: 50px;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 10px var(--shadow-color);
  }

  .theme-toggle:hover {
    transform: scale(1.1) rotate(15deg);
    border-color: var(--primary-color);
    box-shadow: 0 4px 20px rgba(138, 124, 216, 0.3);
  }

  .theme-icon {
    font-size: 1.5rem;
    transition: transform var(--transition-speed) ease;
  }

  .theme-toggle:hover .theme-icon {
    transform: scale(1.2);
  }

  @media (max-width: 768px) {
    .header-top {
      flex-direction: column;
      gap: 1rem;
    }

    .theme-toggle {
      position: static;
    }
  }

  .app-subtitle {
    font-size: 1.1rem;
    color: var(--muted-text-color);
    margin: 0;
  }

  /* Error Alert - Improved animation */
  .error-alert {
    background: linear-gradient(135deg, var(--error-color), #c0392b);
    color: white;
    padding: 1rem 1.5rem;
    border-radius: var(--border-radius);
    margin-bottom: 2rem;
    box-shadow: 0 4px 20px rgba(231, 76, 60, 0.3);
    animation: slideDown 0.3s ease-out;
  }

  .error-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
  }

  .error-text {
    flex: 1;
  }

  .error-close {
    background: none;
    border: none;
    color: white;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: background-color var(--transition-speed) ease;
  }

  .error-close:hover {
    background-color: rgba(255, 255, 255, 0.2);
  }

  /* Animations */
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* Navigation - Improved responsive design */
  .nav-container {
    margin-bottom: 2rem;
  }
  
  .nav-wrapper {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    background: linear-gradient(135deg, var(--card-gradient-start), var(--card-gradient-end));
    border-radius: var(--border-radius);
    padding: 0.75rem;
    box-shadow: 0 2px 15px var(--shadow-color);
    border: 1px solid var(--border-color);
  }

  .nav-button {
    flex: 1;
    min-width: 120px;
    padding: 0.875rem 1rem;
    border: none;
    border-radius: 8px;
    background-color: transparent;
    color: var(--muted-text-color);
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
  }
  
  .nav-button:hover {
    background-color: rgba(138, 124, 216, 0.1);
    color: var(--primary-color);
    transform: translateY(-2px);
  }
  
  .nav-active {
    background: linear-gradient(135deg, var(--primary-color), #5a4acd);
    color: white !important;
    box-shadow: 0 4px 15px rgba(138, 124, 216, 0.4);
  }

  .nav-icon {
    font-size: 1.1rem;
  }

  @media (max-width: 768px) {
    .nav-wrapper {
      flex-direction: column;
    }
    
    .nav-button {
      min-width: auto;
    }
  }

  /* Add Task Form - FIXED LAYOUT */
  .add-task-card {
    background: linear-gradient(135deg, var(--card-gradient-start), var(--card-gradient-end));
    border-radius: var(--border-radius);
    padding: 2rem;
    margin-bottom: 2rem;
    box-shadow: 0 4px 25px var(--shadow-color);
    border: 1px solid var(--border-color);
  }

  .form-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0 0 1.5rem 0;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .form-icon {
    color: var(--primary-color);
  }

  /* FIXED: Better form layout */
  .task-form {
    width: 100%;
  }

  .form-row {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
    flex-wrap: wrap;
  }

  .input-group {
    position: relative;
  }

  .title-input {
    flex: 2;
    min-width: 250px;
  }

  .input-group:not(.title-input) {
    flex: 0 0 auto;
    min-width: 150px;
  }

  /* FIXED: Input styles with better borders */
  .form-input, .form-select {
    width: 100%;
    padding: 0.875rem 1rem;
    border: 2px solid var(--border-color);
    border-radius: 10px;
    font-size: 1rem;
    transition: all var(--transition-speed) ease;
    background-color: var(--surface-color);
    color: var(--text-color);
    box-sizing: border-box;
  }

  .form-input:focus, .form-select:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(138, 124, 216, 0.1);
    transform: translateY(-1px);
  }

  .input-error {
    border-color: var(--error-color) !important;
    box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
  }

  .input-error-text {
    position: absolute;
    top: 100%;
    left: 0;
    font-size: 0.875rem;
    color: var(--error-color);
    margin-top: 0.25rem;
  }
  
  /* FIXED: Better button styling */
  .add-button {
    padding: 0.875rem 1.5rem;
    background: linear-gradient(135deg, var(--primary-color), #5a4acd);
    border: none;
    border-radius: 10px;
    color: white;
    font-weight: 600;
    font-size: 1rem;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    white-space: nowrap;
    min-width: 140px;
    justify-content: center;
  }

  .add-button:hover:not(:disabled) {
    background: linear-gradient(135deg, #5a4acd, #4a3abd);
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(138, 124, 216, 0.4);
  }
  
  .add-button:disabled {
    background: linear-gradient(135deg, #ccc, #bbb);
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }

  .loading .button-text {
    opacity: 0.8;
  }

  /* Responsive form layout */
  @media (max-width: 768px) {
    .form-row {
      flex-direction: column;
    }
    
    .input-group {
      width: 100%;
    }
    
    .title-input {
      min-width: auto;
    }
  }

  /* Dashboard Stats - Improved */
  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
  }

  .stat-card {
    background: linear-gradient(135deg, var(--card-gradient-start), var(--card-gradient-end));
    border-radius: var(--border-radius);
    padding: 2rem 1.5rem;
    box-shadow: 0 4px 25px var(--shadow-color);
    text-align: center;
    border: 1px solid var(--border-color);
    transition: all var(--transition-speed) ease;
  }

  .stat-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 35px var(--shadow-color);
  }

  .stat-icon {
    font-size: 2.5rem;
    margin-bottom: 1rem;
  }

  .stat-title {
    font-size: 1rem;
    color: var(--muted-text-color);
    margin: 0 0 0.5rem 0;
    font-weight: 500;
  }

  .stat-number {
    font-size: 2.5rem;
    font-weight: 700;
    color: var(--primary-color);
    margin: 0;
  }

  /* Dashboard Sections - Improved */
  .dashboard-sections {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }

  @media (max-width: 968px) {
    .dashboard-sections {
      grid-template-columns: 1fr;
    }
  }

  .section-card {
    background: linear-gradient(135deg, var(--card-gradient-start), var(--card-gradient-end));
    border-radius: var(--border-radius);
    padding: 2rem;
    box-shadow: 0 4px 25px var(--shadow-color);
    border: 1px solid var(--border-color);
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0 0 1.5rem 0;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .section-icon {
    font-size: 1.3rem;
  }
  
  /* Task List & Items - FIXED PRIORITY COLORS */
  .task-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .task-item, .task-card {
    background-color: var(--surface-color);
    border-radius: 10px;
    padding: 1.25rem 1.5rem;
    border: 2px solid var(--border-color);
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
    box-shadow: 0 6px 20px var(--shadow-color);
    border-color: var(--primary-color);
  }

  .task-content {
    flex: 1;
  }

  .task-title {
    font-weight: 600;
    margin: 0 0 0.5rem 0;
    font-size: 1.1rem;
    color: var(--text-color);
  }

  .task-meta {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    font-size: 0.9rem;
    flex-wrap: wrap;
  }
  
  /* FIXED: Better priority badge styling with distinct borders - Dark theme compatible */
  .task-priority {
    padding: 0.375rem 0.875rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: capitalize;
    border: 2px solid;
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
  }

  /* Light theme priority colors */
  :global([data-theme="light"]) .priority-high { 
    background-color: #ffebee; 
    color: #c62828; 
    border-color: #e57373;
  }
  
  :global([data-theme="light"]) .priority-medium { 
    background-color: #fff8e1; 
    color: #f57c00; 
    border-color: #ffb74d;
  }
  
  :global([data-theme="light"]) .priority-low { 
    background-color: #e8f5e8; 
    color: #2e7d32; 
    border-color: #81c784;
  }

  :global([data-theme="light"]) .priority-default { 
    background-color: #f5f5f5; 
    color: #666; 
    border-color: #ddd;
  }

  /* Dark theme priority colors */
  :global([data-theme="dark"]) .priority-high { 
    background-color: #4c1d1d; 
    color: #ff6b6b; 
    border-color: #d32f2f;
  }
  
  :global([data-theme="dark"]) .priority-medium { 
    background-color: #4c3d1a; 
    color: #ffb347; 
    border-color: #ff9800;
  }
  
  :global([data-theme="dark"]) .priority-low { 
    background-color: #1e4c1e; 
    color: #66bb6a; 
    border-color: #4caf50;
  }

  :global([data-theme="dark"]) .priority-default { 
    background-color: #404040; 
    color: #a1a1aa; 
    border-color: #525252;
  }

  /* Task status badges */
  .task-status {
    padding: 0.375rem 0.875rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: capitalize;
    border: 2px solid;
  }

  /* Light theme status colors */
  :global([data-theme="light"]) .status-active { 
    background-color: #e3f2fd; 
    color: #1565c0; 
    border-color: #64b5f6;
  }
  
  :global([data-theme="light"]) .status-completed { 
    background-color: #e8f5e8; 
    color: #2e7d32; 
    border-color: #81c784;
  }

  :global([data-theme="light"]) .status-default { 
    background-color: #f5f5f5; 
    color: #666; 
    border-color: #ddd;
  }

  /* Dark theme status colors */
  :global([data-theme="dark"]) .status-active { 
    background-color: #1a2332; 
    color: #60a5fa; 
    border-color: #3b82f6;
  }
  
  :global([data-theme="dark"]) .status-completed { 
    background-color: #1e4c1e; 
    color: #66bb6a; 
    border-color: #4caf50;
  }

  :global([data-theme="dark"]) .status-default { 
    background-color: #404040; 
    color: #a1a1aa; 
    border-color: #525252;
  }

  .task-date {
    color: var(--muted-text-color);
    font-size: 0.85rem;
  }

  .overdue-badge {
    color: var(--error-color);
    font-weight: 600;
  }

  /* Complete button */
  .complete-button {
    padding: 0.5rem 1rem;
    background: linear-gradient(135deg, var(--success-color), #27ae60);
    border: none;
    border-radius: 8px;
    color: white;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
    white-space: nowrap;
  }

  .complete-button:hover:not(:disabled) {
    background: linear-gradient(135deg, #27ae60, #229954);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(46, 204, 113, 0.4);
  }

  .complete-button:disabled {
    background: #ccc;
    cursor: not-allowed;
  }

  /* Task List View - Improved */
  .task-list-card {
    background: linear-gradient(135deg, var(--card-gradient-start), var(--card-gradient-end));
    border-radius: var(--border-radius);
    padding: 2rem;
    box-shadow: 0 4px 25px var(--shadow-color);
    border: 1px solid var(--border-color);
  }

  .task-list-header {
    margin-bottom: 1.5rem;
  }

  .list-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .task-count {
    color: var(--muted-text-color);
    font-weight: 400;
  }
  
  .task-badges {
    display: flex;
    gap: 0.75rem;
    margin-top: 1rem;
    flex-wrap: wrap;
  }

  .task-view {
    width: 100%;
  }

  .task-main {
    margin-bottom: 1rem;
  }

  .task-actions {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .action-button {
    padding: 0.5rem 1rem;
    border: 2px solid var(--border-color);
    border-radius: 8px;
    background-color: var(--surface-color);
    color: var(--text-color);
    cursor: pointer;
    transition: all var(--transition-speed) ease;
    font-weight: 500;
    white-space: nowrap;
  }
  
  .action-button:hover:not(:disabled) {
    transform: translateY(-1px);
  }

  .action-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .edit-btn:hover:not(:disabled) { 
    background-color: rgba(59, 130, 246, 0.1);
    border-color: #3b82f6; 
    color: #3b82f6; 
  }
  
  .complete-btn:hover:not(:disabled) { 
    background-color: var(--success-color); 
    border-color: var(--success-color); 
    color: white; 
  }
  
  .delete-btn:hover:not(:disabled) { 
    background-color: var(--error-color); 
    border-color: var(--error-color); 
    color: white; 
  }
  
  .task-card.overdue {
    border-left: 4px solid var(--error-color);
    background: linear-gradient(135deg, rgba(239, 68, 68, 0.05), var(--surface-color));
  }
  
  /* Edit Form - Improved */
  .edit-form {
    width: 100%;
    padding: 1rem;
    background: var(--surface-hover);
    border-radius: 8px;
    border: 2px solid var(--border-color);
  }
  
  .edit-input-group {
    margin-bottom: 1rem;
  }
  
  .edit-input, .edit-select, .edit-date {
    width: 100%;
    padding: 0.75rem;
    border: 2px solid var(--border-color);
    border-radius: 8px;
    font-size: 1rem;
    background-color: var(--surface-color);
    color: var(--text-color);
    transition: all var(--transition-speed) ease;
    box-sizing: border-box;
  }

  .edit-input:focus, .edit-select:focus, .edit-date:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(138, 124, 216, 0.1);
  }

  .edit-controls {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    align-items: center;
  }

  .save-button, .cancel-button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
    white-space: nowrap;
  }

  .save-button {
    background: linear-gradient(135deg, var(--success-color), #27ae60);
    color: white;
  }

  .save-button:hover:not(:disabled) {
    background: linear-gradient(135deg, #27ae60, #229954);
    transform: translateY(-1px);
  }

  .cancel-button {
    background: linear-gradient(135deg, #95a5a6, #7f8c8d);
    color: white;
  }

  .cancel-button:hover:not(:disabled) {
    background: linear-gradient(135deg, #7f8c8d, #6c7b7d);
    transform: translateY(-1px);
  }

  .save-button:disabled, .cancel-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }
  
  /* Empty & Loading States - Enhanced */
  .empty-state, .loading-state {
    text-align: center;
    padding: 4rem 2rem;
    color: var(--muted-text-color);
  }

  .empty-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
  }

  .empty-text {
    font-size: 1.1rem;
    margin: 0;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid var(--border-color);
    border-top: 4px solid var(--primary-color);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 1rem;
  }

  .loading-text {
    font-size: 1.1rem;
    margin: 0;
    color: var(--muted-text-color);
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  /* Tasks container */
  .tasks-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  /* Improved responsive design */
  @media (max-width: 480px) {
    .app-container {
      padding: 1rem;
    }

    .task-actions {
      justify-content: center;
    }

    .action-button {
      flex: 1;
      text-align: center;
    }

    .edit-controls {
      flex-direction: column;
    }

    .edit-select, .edit-date {
      margin-bottom: 0.5rem;
    }
  }

  /* Better focus states for accessibility */
  button:focus-visible {
    outline: 2px solid var(--primary-color);
    outline-offset: 2px;
  }

  input:focus-visible, select:focus-visible {
    outline: 2px solid var(--primary-color);
    outline-offset: 2px;
  }

  /* Additional dark theme improvements */
  :global([data-theme="dark"]) .loading-spinner {
    border-color: var(--border-color);
    border-top-color: var(--primary-color);
  }

  /* Smooth transitions for theme switching */
  * {
    transition: background-color var(--transition-speed) ease, 
                border-color var(--transition-speed) ease, 
                color var(--transition-speed) ease;
  }
</style>