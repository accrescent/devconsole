const appInfo = document.getElementById("app_info");

const id = document.getElementById("id");
const label = document.getElementById("label");
const versionName = document.getElementById("version_name");
const versionCode = document.getElementById("version_code");

document.getElementById("new_app_form").onsubmit = event => {
    event.preventDefault();

    let input = document.querySelector("input[type='file']");
    let data = new FormData();
    data.append("file", input.files[0]);

    fetch("/api/apps", {
        method: "POST",
        mode: "same-origin",
        body: data,
    }).then(resp => {
        if (!resp.ok) {
            return Promise.reject();
        }
        return resp.json();
    }).then(app => {
        id.innerText = `App ID: ${app.id}`;
        label.innerText = `Display name: ${app.label}`;
        versionName.innerText = `Display version: ${app.version_name}`;
        versionCode.innerText = `Version code: ${app.version_code}`;

        appInfo.hidden = false;
    }).catch(err => {
        appInfo.hidden = true;

        console.error(err);
    });
};

document.getElementById("submit").onclick = () => {
    fetch("/api/apps", { method: "PATCH", mode: "same-origin", }).then(resp => {
        if (!resp.ok) {
            return Promise.reject();
        }
    }).catch(console.error);
};
