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

            const container = button.parentNode.parentNode;
            if (container.childElementCount > 1) {
                container.removeChild(button.parentNode);
            } else {
                location.reload();
            }
        }).catch(console.error);
    }
}
