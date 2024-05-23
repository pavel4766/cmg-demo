import React, { useState, useEffect } from 'react';
import axios from 'axios';

function App() {
  const [tasks, setTasks] = useState([]);
  const [input, setInput] = useState('');

  // Function to fetch tasks from backend
  const fetchTasks = async () => {
    const response = await axios.get('http://localhost:3001/api');
    setTasks(response.data.todos);
  };

  // Fetch tasks when the component mounts
  useEffect(() => {
    fetchTasks();
  }, []);

  // Function to add a task
  const addTask = async () => {
    if (input) {
      await axios.post('http://localhost:3001/api', { item: input });
      setInput('');
      fetchTasks();
    }
  };

  // Function to update a task
  const updateTask = async (id, item) => {
    await axios.put('http://localhost:3001/api/update', { id, item });
    fetchTasks();
  };

  // Function to delete a task
  const deleteTask = async (id) => {
    await axios.delete('http://localhost:3001/api/delete', { data: { id } });
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
