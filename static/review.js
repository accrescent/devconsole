const buttons = document.querySelectorAll("button[name='approve']");

for (let button of buttons) {
    button.onclick = event => {
        event.preventDefault();

        let appId = button.value;

        fetch("/api/apps/approve", {
            method: "POST",
            mode: "same-origin",
            body: JSON.stringify({ app_id: appId }),
        }).then(resp => {
            if (!resp.ok) {
                return Promise.reject();
            }
        }).catch(console.error);
    }
}
