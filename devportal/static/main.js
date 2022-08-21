document.getElementById("logout_button").onclick = () => {
    fetch("/api/logout", { method: "POST", mode: "same-origin" }).then(resp => {
        if (!resp.ok) {
            return Promise.reject();
        }
        location.replace("/");
    }).catch(console.error);
};
