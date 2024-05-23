import React, { useState, useEffect } from 'react';
import axios from 'axios';

function App() {
  const [tasks, setTasks] = useState([]);
  const [input, setInput] = useState('');
  const apiBaseUrl = process.env.REACT_APP_API_URL || 'http://localhost:3001';  // Fallback URL if environment variable is not set

  // Function to fetch tasks from the backend
  const fetchTasks = async () => {
    const response = await axios.get(`${apiBaseUrl}/api`);
    setTasks(response.data.todos);
  };

  // Fetch tasks when the component mounts
  useEffect(() => {
    fetchTasks();
  }, []);

  // Function to add a task
  const addTask = async () => {
    if (input) {
      await axios.post(`${apiBaseUrl}/api`, { item: input });
      setInput('');
      fetchTasks();
    }
  };

  // Function to update a task
  const updateTask = async (id, item) => {
    await axios.put(`${apiBaseUrl}/api/update`, { id, item });
    fetchTasks();
  };

  // Function to delete a task
  const deleteTask = async (id) => {
    await axios.delete(`${apiBaseUrl}/api/delete`, { data: { id } });
    fetchTasks();
  };

  return (
    <div className="App">
      <input
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Add a new task"
      />
      <button onClick={addTask}>Add Task</button>
      <ul>
        {tasks.map((task) => (
          <li key={task.id}>
            {task.item}
            <button onClick={() => updateTask(task.id, prompt('Update the task:', task.item))}>
              Update
            </button>
            <button onClick={() => deleteTask(task.id)}>
              Delete
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
