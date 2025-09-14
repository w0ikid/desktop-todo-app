<script lang="ts">
  export let newTitle: string;
  export let newPriority: string;
  export let newDueDate: string;
  export let titleError: string;
  export let loading: boolean;
  export let addTask: () => Promise<void>;
  export let handleKeyPress: (event: KeyboardEvent) => void;
</script>

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