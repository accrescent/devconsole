document.getElementById("logout_button").onclick = () => {
    fetch("/logout", { method: "POST", mode: "same-origin" }).then(resp => {
        if (!resp.ok) {
            return Promise.reject();
        }
        location.replace("/");
    }).catch(err => {
        console.log(err)
    });
};
