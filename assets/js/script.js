document.title = 'Task Day 3 - Form Submission';


function getValues() {
    let name = document.querySelector('#name').value;
    let email = document.querySelector('#email').value;
    let phone = document.querySelector('#phone').value;
    let subject = document.querySelector('#subject').value;
    let message = document.querySelector('#message').value;

    if(name == "") {
        alert("Nama harus di isi")
    } else if(email == "") {
        alert("email harus di isi")
    } else if(phone == "") {
        alert("phone harus di isi")
    } else if(subject == "") {
        alert("subject harus di isi")
    } else if(message == "") {
        alert("Pesan harus di isi")
    } else {
        const defaultEmail = "vindokountur@gmail.com"
    
        let mailTo = document.createElement('a')
        mailTo.href = `mailto:${defaultEmail}?subject=${subject}&body=Halo saya ${name}, ${message} jika tertarik bisa hubungi di nomor ini ${phone}`
        mailTo.click()
    }
}