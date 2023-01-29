document.title = "Task Day 5 - My Profile Page";

// Contact
function getValues() {
  let name = document.querySelector("#name").value;
  let email = document.querySelector("#email").value;
  let phone = document.querySelector("#phone").value;
  let subject = document.querySelector("#subject").value;
  let message = document.querySelector("#message").value;

  if (name == "") {
    alert("Nama harus di isi");
  } else if (email == "") {
    alert("email harus di isi");
  } else if (phone == "") {
    alert("phone harus di isi");
  } else if (subject == "") {
    alert("subject harus di isi");
  } else if (message == "") {
    alert("Pesan harus di isi");
  } else {
    const defaultEmail = "vindokountur@gmail.com";

    let mailTo = document.createElement("a");
    mailTo.href = `mailto:${defaultEmail}?subject=${subject}&body=Halo saya ${name}, ${message} jika tertarik bisa hubungi di nomor ini ${phone}`;
    mailTo.click();
  }
}

// My Project
let projectList = [];

function getProjectData(e) {
  e.preventDefault();

  let projectName = document.querySelector("#projectName").value;
  let startDate = document.querySelector("#startDate").value;
  let endDate = document.querySelector("#endDate").value;
  let description = document.querySelector("#description").value;
  let nodejs = document.querySelector("#nodejs").checked;
  let vuejs = document.querySelector("#vuejs").checked;
  let reactjs = document.querySelector("#reactjs").checked;
  let php = document.querySelector("#php").checked;
  let uploadImage = document.querySelector("#uploadImage").files;

  if (projectName == "") {
    return alert("Nama Project harus diisi");
  }
  if (startDate == "") {
    return alert("Start Date harus diisi");
  }
  if (endDate == "") {
    return alert("End Date harus diisi");
  }
  if (description == "") {
    return alert("Description belum diisi");
  }
  if (!nodejs && !vuejs && !reactjs && !php) {
    return alert("Technologi belum dipilih");
  }
  if (uploadImage.length == 0) {
    return alert("Image belum dipilih");
  }

  uploadImage = URL.createObjectURL(uploadImage[0]);

  let projectData = {
    name: projectName,
    startDate,
    endDate,
    description,
    tech: {
      nodejs,
      vuejs,
      reactjs,
      php,
    },
    uploadImage,
  };

  projectList.push(projectData);
  console.log(projectList);
  showProjectList();
}

function showProjectList() {
  let showProject = document.querySelector(".project-list-container");
  showProject.innerHTML = "";

  if (projectList.length == 0) {
    showProject.innerHTML = `<p class='empty'>Its empty here<p>`;
    return;
  }
  for (let i = 0; i < projectList.length; i++) {
    const project = projectList[i];
    showProject.innerHTML += `
        <div class="project-card">
        <img src="${project.uploadImage}" alt="project-image" />
        <a class="title" href='/detailproject.html'>${project.name}</a>
        <p class="duration">Durasi : ${project.startDate} - ${
      project.endDate
    }</p>
        <p class="description">${project.description}</p>
        <div class="tech-list">
            ${
              project.tech.nodejs
                ? '<i class="fa-brands fa-xl fa-node"></i>'
                : ""
            }
            ${
              project.tech.vuejs
                ? '<i class="fa-brands fa-xl fa-vuejs"></i>'
                : ""
            }
            ${
              project.tech.reactjs
                ? '<i class="fa-brands fa-xl fa-react"></i>'
                : ""
            }
            ${project.tech.php ? '<i class="fa-brands fa-xl fa-php"></i>' : ""}
        </div>
        <div class="option">
          <button>edit</button>
          <button>delete</button>
        </div>
      </div>
        `;
  }
}

showProjectList();
