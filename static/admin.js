const buttons = document.querySelectorAll("button[name='app_id']");

for (let button of buttons) {
    button.onclick = event => {
        event.preventDefault();

        let appId = button.value;

        fetch(`/api/apps/${appId}`, { method: "POST", mode: "same-origin" }).then(resp => {
            if (!resp.ok) {
                return Promise.reject();
            }
        }).catch(console.error);
    }
}
