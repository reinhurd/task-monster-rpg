package core

// TODO rework template
func (s *Service) GetTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Task Manager</title>
  <style>
    /* CSS styles */
    body {
      font-family: Arial, sans-serif;
      margin: 0;
      padding: 0;
      background-color: #f4f4f4;
    }
    header {
      background: #333;
      color: #fff;
      padding: 10px;
      text-align: center;
    }
    .container {
      width: 90%;
      max-width: 1200px;
      margin: 20px auto;
    }
    .button {
      background: #333;
      color: #fff;
      padding: 10px 20px;
      border: none;
      cursor: pointer;
    }
    .button:hover {
      background: #555;
    }
    .task {
      background: #fff;
      margin-bottom: 10px;
      padding: 10px;
      border-radius: 5px;
    }
    .task-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .task-details {
      display: none;
      padding: 10px;
      border-top: 1px solid #ccc;
    }
    /* Modal styles */
    .modal {
      display: none;
      position: fixed;
      z-index: 1000;
      left: 0;
      top: 0;
      width: 100%;
      height: 100%;
      overflow: auto;
      background-color: rgba(0,0,0,0.5);
    }
    .modal-content {
      background-color: #fefefe;
      margin: 15% auto;
      padding: 20px;
      border: 1px solid #888;
      width: 80%;
      max-width: 500px;
    }
    .close {
      color: #aaa;
      float: right;
      font-size: 28px;
      font-weight: bold;
      cursor: pointer;
    }
  </style>
</head>
<body>
  <header>
    <h1>Task Manager</h1>
  </header>
  <div class="container">
    <!-- Login Section -->
    <div id="loginSection">
      <h2>Login</h2>
      <form id="loginForm">
        <label for="login">Login:</label><br>
        <input type="text" id="login" name="login" required><br>
        <label for="password">Password:</label><br>
        <input type="password" id="password" name="password" required><br><br>
        <button type="submit" class="button">Login</button>
      </form>
      <p>Or <a href="/api/users">Create New User</a></p>
    </div>

    <!-- Tasks Section (visible only if authorized) -->
    <div id="tasksSection" style="display:none;">
      <h2>Your Tasks</h2>
      <div id="tasksContainer"></div>
      <h3>Create New Task</h3>
      <div>
        <button id="manualCreateButton" class="button">Manual Create</button>
        <button id="gptCreateButton" class="button">GPT Help Create</button>
      </div>
      <!-- Manual Create Form (hidden by default) -->
      <div id="manualCreateForm" style="display:none; margin-top:20px;">
        <h3>Manual Task Creation</h3>
        <form id="createTaskForm">
          <label for="title">Title:</label><br>
          <input type="text" id="title" name="title" required><br>
          <label for="description">Description:</label><br>
          <textarea id="description" name="description" required></textarea><br>
          <label for="executor">Executor:</label><br>
          <input type="text" id="executor" name="executor" required><br>
          <label for="deadline">Deadline:</label><br>
          <input type="datetime-local" id="deadline" name="deadline" required><br><br>
          <button type="submit" class="button">Create Task</button>
        </form>
      </div>
    </div>
  </div>

  <!-- GPT Help Modal -->
  <div id="gptModal" class="modal">
    <div class="modal-content">
      <span class="close" id="gptModalClose">&times;</span>
      <h3>GPT Help for Task Creation</h3>
      <label for="gptRequest">Enter your request:</label>
      <input type="text" id="gptRequest" name="gptRequest" style="width:100%;"><br><br>
      <button id="sendGptRequest" class="button">Send Request</button>
    </div>
  </div>

  <script>
    // For demonstration purposes we use localStorage to keep a token.
    // In production, you might use cookies or another secure storage.
    var isAuthorized = false;

    // Handle login submission
    document.getElementById('loginForm').addEventListener('submit', function(e) {
      e.preventDefault();
      var login = document.getElementById('login').value;
      var password = document.getElementById('password').value;
      // POST to /api/users/login with the userCreateRequest payload
      fetch('/api/users/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: login, password: password })
      })
      .then(response => response.json())
      .then(data => {
        if(data.success) {
          isAuthorized = true;
          // Save token if returned
          localStorage.setItem('token', data.token);
          document.getElementById('loginSection').style.display = 'none';
          document.getElementById('tasksSection').style.display = 'block';
          loadTasks();
        } else {
          alert('Login failed');
        }
      })
      .catch(err => console.error('Error during login:', err));
    });

    // Load tasks from /api/tasks and render them in the tasksContainer
    function loadTasks() {
      fetch('/api/tasks', {
        headers: { 'Authorization': '' + localStorage.getItem('token') }
      })
      .then(response => response.json())
      .then(tasks => {
        var tasksContainer = document.getElementById('tasksContainer');
        tasksContainer.innerHTML = '';
        tasks.forEach(task => {
          var taskDiv = document.createElement('div');
          taskDiv.className = 'task';

          // Task header with title and a toggle button for details
          var headerDiv = document.createElement('div');
          headerDiv.className = 'task-header';
          headerDiv.innerHTML = '<strong>' + task.Title + '</strong>';
          var toggleButton = document.createElement('button');
          toggleButton.className = 'button';
          toggleButton.textContent = 'Details';
          toggleButton.addEventListener('click', function() {
            if(detailsDiv.style.display === 'none' || detailsDiv.style.display === '') {
              detailsDiv.style.display = 'block';
            } else {
              detailsDiv.style.display = 'none';
            }
          });
          headerDiv.appendChild(toggleButton);
          taskDiv.appendChild(headerDiv);

          // Details section showing all task info
          var detailsDiv = document.createElement('div');
          detailsDiv.className = 'task-details';
          detailsDiv.innerHTML =
            '<p><strong>Description:</strong> ' + task.Description + '</p>' +
            '<p><strong>Executor:</strong> ' + task.Executor + '</p>' +
            '<p><strong>Reviewer:</strong> ' + (task.Reviewer || '') + '</p>' +
            '<p><strong>Completed:</strong> ' + task.Completed + '</p>' +
            '<p><strong>Created At:</strong> ' + new Date(task.Created_at).toLocaleString() + '</p>' +
            '<p><strong>Updated At:</strong> ' + new Date(task.Updated_at).toLocaleString() + '</p>';

          // Update form for the task – sends PUT request to /api/tasks/:id/status
          var updateForm = document.createElement('form');
          updateForm.innerHTML = '<h4>Update Task</h4>' +
            '<label>Title: <input type="text" name="title" value="' + task.Title + '"></label><br>' +
            '<label>Description: <textarea name="description">' + task.Description + '</textarea></label><br>' +
            '<label>Executor: <input type="text" name="executor" value="' + task.Executor + '"></label><br>' +
            '<label>Reviewer: <input type="text" name="reviewer" value="' + (task.Reviewer || '') + '"></label><br>' +
            '<label>Completed: <input type="checkbox" name="completed" ' + (task.Completed ? 'checked' : '') + '></label><br>' +
            '<button type="submit" class="button">Update Task</button>';

          updateForm.addEventListener('submit', function(e) {
            e.preventDefault();
            var formData = new FormData(updateForm);
            var updatedTask = {
              title: formData.get('Title'),
              description: formData.get('Description'),
              executor: formData.get('Executor'),
              reviewer: formData.get('Reviewer'),
              completed: formData.get('Completed') === 'on'
            };
            fetch('/api/tasks/' + task._id + '/status', {
              method: 'PUT',
              headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('token')
              },
              body: JSON.stringify(updatedTask)
            })
            .then(response => response.json())
            .then(result => {
              if(result.success) {
                alert('Task updated successfully');
                loadTasks();
              } else {
                alert('Failed to update task');
              }
            })
            .catch(err => console.error('Error updating task:', err));
          });

          detailsDiv.appendChild(updateForm);
          taskDiv.appendChild(detailsDiv);
          tasksContainer.appendChild(taskDiv);
        });
      })
      .catch(err => console.error('Error fetching tasks:', err));
    }

    // Toggle display of the manual create task form
    document.getElementById('manualCreateButton').addEventListener('click', function() {
      var form = document.getElementById('manualCreateForm');
      form.style.display = (form.style.display === 'none' || form.style.display === '') ? 'block' : 'none';
    });

    // Handle manual task creation (POST to /api/tasks/create)
    document.getElementById('createTaskForm').addEventListener('submit', function(e) {
      e.preventDefault();
      var formData = new FormData(document.getElementById('createTaskForm'));
      var newTask = {
        title: formData.get('title'),
        description: formData.get('description'),
        executor: formData.get('executor'),
        deadline: formData.get('deadline')
      };
      fetch('/api/tasks/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify(newTask)
      })
      .then(response => response.json())
      .then(result => {
        if(result.success) {
          alert('Task created successfully');
          loadTasks();
          document.getElementById('createTaskForm').reset();
          document.getElementById('manualCreateForm').style.display = 'none';
        } else {
          alert('Failed to create task');
        }
      })
      .catch(err => console.error('Error creating task:', err));
    });

    // GPT Help Modal logic
    var gptModal = document.getElementById('gptModal');
    var gptModalClose = document.getElementById('gptModalClose');
    document.getElementById('gptCreateButton').addEventListener('click', function() {
      gptModal.style.display = 'block';
    });
    gptModalClose.addEventListener('click', function() {
      gptModal.style.display = 'none';
    });
    window.addEventListener('click', function(event) {
      if (event.target == gptModal) {
        gptModal.style.display = 'none';
      }
    });
    // Send GPT request (GET to /api/tasks/create/gpt?req=...)
    document.getElementById('sendGptRequest').addEventListener('click', function() {
      var reqText = document.getElementById('gptRequest').value;
      fetch('/api/tasks/create/gpt?req=' + encodeURIComponent(reqText), {
        method: 'GET',
        headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
      })
      .then(response => response.json())
      .then(result => {
        if(result.success) {
          alert('GPT task created successfully');
          loadTasks();
          gptModal.style.display = 'none';
          document.getElementById('gptRequest').value = '';
        } else {
          alert('GPT task creation failed');
        }
      })
      .catch(err => console.error('Error with GPT task creation:', err));
    });
  </script>
</body>
</html>`
}
