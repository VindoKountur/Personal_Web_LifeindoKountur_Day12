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

  if (projectName == "") return alert("Nama Project harus diisi");
  if (startDate == "") return alert("Start Date harus diisi");
  if (endDate == "") return alert("End Date harus diisi");
  if (description == "") return alert("Description belum diisi");
  if (!nodejs && !vuejs && !reactjs && !php)
    return alert("Technologi belum dipilih");
  if (uploadImage.length == 0) return alert("Image belum dipilih");

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
  showProjectList();
}

function showProjectList() {
  let showProject = document.querySelector("#project-list-container");
  showProject.innerHTML = "";

  if (projectList.length == 0) {
    showProject.innerHTML = `<p class='text-center md:col-span-2 lg:col-span-3'>Its empty here<p>`;
    return;
  }
  for (let i = 0; i < projectList.length; i++) {
    const project = projectList[i];
    showProject.innerHTML += `
      <div class=" bg-white drop-shadow-lg rounded p-3">
        <img class="rounded-xl object-cover" src="${project.uploadImage}" alt="project-image" />
        <a href="/detail" class="mt-2 font-bold text-lg">${project.name}</a>
        <p class="text-gray-600">Durasi : ${getTimeDifferences(project.startDate, project.endDate)}</p>
        <p class="mb-7 mt-5 font-semibold">${project.description}</p>
        <div class="flex items-center gap-5">
        ${project.tech.nodejs ? '<i class="fa-brands fa-2xl fa-node"></i>' : ''}
        ${project.tech.vuejs ? '<i class="fa-brands fa-2xl fa-vuejs"></i>' : ''}
        ${project.tech.reactjs ? '<i class="fa-brands fa-2xl fa-react"></i>' : ''}
        ${project.tech.php ? '<i class="fa-brands fa-2xl fa-php"></i>' : ''}
        </div>
        <div class="mt-8 flex justify-between gap-4">
          <button class="w-full bg-black text-white rounded-lg py-1 font-semibold">edit</button>
          <button class="w-full bg-black text-white rounded-lg py-1 font-semibold">delete</button>
        </div>
      </div>
        `;
  }
}

const getTimeDifferences = (start, end) => {
  let timeDifferent = new Date(end) - new Date(start);

  let monthDistance = Math.floor(timeDifferent / (30 * 24 * 60 * 60 * 1000));
  let weekDistance = Math.floor(timeDifferent / (7 * 24 * 60 * 60 * 1000));
  let dayDistance = Math.floor(timeDifferent / (24 * 60 * 60 * 1000));

  if (monthDistance > 0) return `${monthDistance} month`;
  if (weekDistance > 0) return `${weekDistance} week`;
  if (dayDistance > 0) return `${dayDistance} day`;
  return "Cannot count time";
};