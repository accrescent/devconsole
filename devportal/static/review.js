const approveNewButtons = document.querySelectorAll("button[name='approve']");
const approveUpdateButtons = document.querySelectorAll("button[name='approve_update']");

for (const button of approveNewButtons) {
    button.onclick = event => {
        event.preventDefault();

        const appId = button.value;

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

for (const button of approveUpdateButtons) {
    button.onclick = event => {
        event.preventDefault();

        const appId = button.value;

        fetch(`/api/apps/${appId}/approve`, { method: "POST", mode: "same-origin" }).then(resp => {
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
