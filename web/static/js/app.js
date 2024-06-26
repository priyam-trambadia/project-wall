


// Add Project Floating Button

const addProjectFloatingButton = document.querySelector("#add-project-fbtn")
if (addProjectFloatingButton !== null) {
  addProjectFloatingButton.addEventListener("click", () => {
    const github_url = window.prompt("Enter project github URL")

    if (github_url !== null) {
      const info = extractOwnerAndRepo(github_url)
      if (info.owner === null || info.repo === null) {
        invalidURLAlert()
      }
      else {
        const repo_api_url = "https://api.github.com/repos/" + info.owner + "/" + info.repo

        // extaract repo info
        fetchData(repo_api_url)
          .then((data) => {
            let info = new Object()
            info["title"] = data["name"]
            info["description"] = data["description"]
            info["github_link"] = data["html_url"]

            return fetchData(data["languages_url"])
              .then((language_data) => {
                info["languages"] = Object.keys(language_data)
                return info
              })
              .catch(() => {
                throw new Error(`API github language fatching error`);
              })
          })
          .then((repo_info) => {
            makeRequestAddProject(repo_info)
          })
          .catch(() => {
            fetchErrorAlert()
          })

      }

    }

  })
}

function extractOwnerAndRepo(url) {
  const regex = /https:\/\/github.com\/([^\/]+)\/([^\/]+)/;
  const match = url.match(regex);

  if (match) {
    return {
      owner: match[1],
      repo: match[2]
    };
  } else {
    return {
      owner: null,
      repo: null
    };
  }
}

function invalidURLAlert() {
  window.alert("Invalid github URL!")
}

async function fetchData(url) {
  try {
    const response = await fetch(url);
    if (response.ok) {
      const data = await response.json();
      return data;
    } else {
      throw new Error(`API info failed with status ${response.status}`);
    }
  } catch (error) {
    fetchErrorAlert()
  }
}

function fetchErrorAlert() {
  window.alert("Unable to fetch project info. Try later.")
}

function makeRequestAddProject(json_payload) {
  const request_url = document.location.origin + "/project/add"

  const form = document.createElement('form');
  form.method = 'GET';
  form.action = request_url;

  const payloadInput = document.createElement('input');
  payloadInput.type = 'hidden';
  payloadInput.name = 'project_details_json';
  payloadInput.value = JSON.stringify(json_payload);

  form.appendChild(payloadInput);
  document.body.appendChild(form);
  form.submit();
};



// clickable span handling

const searchSelectionDiv = document.querySelector('.search-selection-container');
if (searchSelectionDiv !== null) {
  searchSelectionDiv.addEventListener('click', function (event) {
    if (event.target.classList.contains('clickable-span')) {

      const clickedSpan = event.target
      // const searchResultsDiv = searchSelectionDiv.querySelector(".search-results")
      const selectedDiv = searchSelectionDiv.querySelector(".selected")

      if (clickedSpan.parentNode === selectedDiv) {
        clickedSpan.remove()
      }
      else {
        let isSameValueSpanExists = false

        selectedDiv.querySelectorAll('*').forEach(childElement => {
          if (childElement.innerHTML === clickedSpan.innerHTML) {
            isSameValueSpanExists = true
          }
        });

        if (!isSameValueSpanExists) {
          selectedDiv.appendChild(clickedSpan.cloneNode(true));
        }
      }

      const form = document.querySelector("form");
      const payloadId = event.currentTarget.id;
      let payloadInput = form.querySelector(`input[name="${payloadId}"]`)

      if (payloadInput === null) {
        payloadInput = document.createElement('input');
        payloadInput.type = 'hidden';
        payloadInput.name = payloadId;
        form.appendChild(payloadInput);
      }

      const childTextContents = [];
      selectedDiv.querySelectorAll('*').forEach(childElement => {
        childTextContents.push(childElement.textContent.trim()); // Add trimmed text content
      });

      payloadInput.value = JSON.stringify(childTextContents);
    }
  });
}