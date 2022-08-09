const buttons = document.querySelectorAll("button[name='app_id']");

for (let i = 0; i < buttons.length; i++) {
    buttons[i].onclick = event => {
        event.preventDefault();

        let appId = buttons[i].value;

        fetch(`/api/apps/${appId}`, { method: "POST", mode: "same-origin" }).then(resp => {
            if (!resp.ok) {
                return Promise.reject();
            }
        }).catch(console.error);
    }
}
