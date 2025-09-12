<script>
  import { onMount } from 'svelte';
  import { GetTasksFromFrontend, MarkTaskCompletedFromFrontend, DeleteTaskFromFrontend } from '../../wailsjs/go/main/App';

  let tasks = [];

  async function loadTasks() {
    try {
      tasks = await GetTasksFromFrontend();
    } catch (error) {
      console.error('Error loading tasks:', error);
    }
  }

  async function toggleTaskCompletion(task) {
    try {
      await MarkTaskCompletedFromFrontend(task.id);
      loadTasks();
    } catch (error) {
      console.error('Error marking task as completed:', error);
    }
  }

  async function deleteTask(task) {
    try {
      await DeleteTaskFromFrontend(task.id);
      loadTasks();
    } catch (error) {
      console.error('Error deleting task:', error);
    }
  }

  onMount(() => {
    loadTasks();
  });
</script>

<div class="task-list">
  <h2>Список задач</h2>
  {#if tasks.length === 0}
    <p>Нет задач для отображения.</p>
  {:else}
    <ul>
      {#each tasks as task (task.id)}
        <li class:completed={task.completed}>
          <span on:click={() => toggleTaskCompletion(task)}>
            {task.name}
          </span>
          <button on:click={() => deleteTask(task)}>Удалить</button>
        </li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  .task-list {
    margin-top: 20px;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px;
    border-bottom: 1px solid #ccc;
  }
  li.completed span {
    text-decoration: line-through;
    color: #888;
  }
  span {
    cursor: pointer;
  }
</style>