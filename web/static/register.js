document.getElementById("register_form").onsubmit = event => {
    event.preventDefault();

    let selectedEmail = document.querySelector("input[name='email']:checked").value;

    fetch("/register/new", {
        method: "POST",
        mode: "same-origin",
        body: JSON.stringify({email: selectedEmail}),
    }).then(resp => {
        if (!resp.ok) {
            if (resp.status === 409) {
                console.log("GitHub user already registered")
            }
            return Promise.reject();
        }
        location.replace("/portal");
    }).catch(err => {
        console.log(err)
    });
};
