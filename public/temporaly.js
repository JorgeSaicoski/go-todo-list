
const taskList = document.getElementById('task-list');
const taskForm = document.getElementById('task-form');  


async function getTasks() {
    const response = await fetch('/tasks');
    const tasks = await response.json();

    taskList.innerHTML = '';
    tasks.forEach(task => {
        console.log(task)
        const li = document.createElement('li');
        li.textContent = `${task.title}:  
        ${task.content}`;
        taskList.appendChild(li);
    });
}

async function createTask(event) {
    event.preventDefault();
    const title = document.getElementById('title').value;
    const content = document.getElementById('content').value;  


    const response = await fetch('/task', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ title, content })
    });

    if (response.ok) {
        getTasks(); // Refresh task list after successful creation
    } else {
        console.error('Error creating task');
    }
}

taskForm.addEventListener('submit', createTask);
getTasks();
