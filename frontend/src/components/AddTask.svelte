<script>
  import { AddTaskFromFrontend } from '../../wailsjs/go/main/App';

  let taskName = '';

  async function addTask() {
    if (taskName.trim() === '') {
      return;
    }
    try {
      await AddTaskFromFrontend(taskName);
      taskName = '';
      window.dispatchEvent(new CustomEvent('taskAdded'));
    } catch (error) {
      console.error('Error adding task:', error);
    }
  }
</script>

<div class="add-task">
  <h2>Добавить задачу</h2>
  <form on:submit|preventDefault={addTask}>
    <input type="text" bind:value={taskName} placeholder="Введите название задачи" />
    <button type="submit">Добавить</button>
  </form>
</div>

<style>
  .add-task {
    margin-top: 20px;
  }
  input {
    padding: 10px;
    width: 70%;
  }
  button {
    padding: 10px;
  }
</style>