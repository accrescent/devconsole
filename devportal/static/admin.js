const buttons = document.querySelectorAll("button[name='publish']");

for (const button of buttons) {
    button.onclick = event => {
        event.preventDefault();

        const appId = button.value;

        fetch(`/api/apps/${appId}`, { method: "POST", mode: "same-origin" }).then(resp => {
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
