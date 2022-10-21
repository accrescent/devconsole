const updateInfo = document.getElementById("update_info");

const reviewSection = document.getElementById("review_section");
const reviewErrors = document.getElementById("review_errors");
const versionChange = document.getElementById("version_change");

const appId = document.getElementById("app_id").textContent;

document.getElementById("update_app_form").onsubmit = event => {
    event.preventDefault();

    const input = document.querySelector("input[type='file']");
    const data = new FormData();
    data.append("file", input.files[0]);

    fetch(`/api/apps/${appId}/updates`, {
        method: "POST",
        mode: "same-origin",
        body: data,
    }).then(resp => {
        if (!resp.ok) {
            return Promise.reject();
        }
        return resp.json();
    }).then(diff => {
        versionChange.innerText = `Update app ${appId} from\n
            ${diff.current_vcode} (${diff.current_vname})\n
            to\n
            ${diff.new_vcode} (${diff.new_vname})`;

        while (reviewErrors.firstChild) {
            reviewErrors.removeChild(reviewErrors.lastChild);
        }
        for (const error of diff.review_errors) {
            const err = document.createElement("li");
            err.innerText = error;
            reviewErrors.appendChild(err);
        }
        if (diff.review_errors.length > 0) {
            reviewSection.hidden = false;
        } else {
            reviewSection.hidden = true;
        }

        updateInfo.hidden = false;
    }).catch(err => {
        updateInfo.hidden = true;
        reviewSection.hidden = true;

        console.error(err);
    });
};

document.getElementById("submit").onclick = () => {
    fetch(`/api/apps/${appId}/updates`, {
        method: "PATCH",
        mode: "same-origin",
    }).then(resp => {
        if (!resp.ok) {
            return Promise.reject();
        }
        location.replace(`/apps/${appId}`);
    }).catch(console.error);
};
