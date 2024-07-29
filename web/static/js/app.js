// data-tooltip handling

const tooltip = document.createElement("span");
tooltip.id = "tooltip";


const body = document.body
if (body !== null) {
  body.addEventListener('pointerover', (event) => {
    if (event.target.hasAttribute("data-tooltip")) {
      tooltip.textContent = event.target.getAttribute("data-tooltip")
      body.appendChild(tooltip);
      positionTooltip(tooltip, event.target);
    }
  });

  body.addEventListener('pointerout', (event) => {
    if (event.target.hasAttribute("data-tooltip")) {
      body.removeChild(tooltip);
    }
  });
}

function positionTooltip(tooltip, element) {

  // Initial positioning below the element
  tooltip.style.position = 'absolute';
  tooltip.style.top = `${element.offsetTop + element.offsetHeight}px`;
  tooltip.style.left = `${element.offsetLeft}px`;

  // Check if the tooltip is fully visible
  const tooltipRect = tooltip.getBoundingClientRect();
  const windowHeight = window.innerHeight;
  if (tooltipRect.bottom > windowHeight) {
    // Reposition the tooltip above the element
    tooltip.style.top = `${element.offsetTop - tooltip.offsetHeight}px`;
  }
}


// popup handling

function showPopup(message, duration = 3600) {
  const popup = document.createElement("span");
  popup.id = "popup"
  body.appendChild(popup)
  popup.textContent = message;

  tenPart = duration / 10
  setTimeout(() => {
    dur1 = 3 * tenPart
    op = 1.0 / dur1
    for (let itr = 1; itr <= dur1; itr += 10) {
      setTimeout(() => {
        popup.style.opacity = (dur1 - itr) * op;
      }, itr)
    }
  }, tenPart * 7);

  setTimeout(() => {
    popup.remove()
  }, duration);
}

function getCookie(name) {
  const cookieValue = `; `;
  const parts = document.cookie.split(cookieValue);
  for (let i = 0; i < parts.length; i++) {
    const part = parts[i];
    if (part.startsWith(name + '=')) {
      return part.substring(name.length + 1);
    }
  }
  return null;
}

function deleteCookie(name) {
  document.cookie = name + "=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
}

function setCookie(cname, cvalue) {
  const date = new Date();
  document.cookie = cname + "=" + cvalue + "; path=/";
}

const popupCookie = getCookie('popup');
if (popupCookie) {
  showPopup(popupCookie)
  deleteCookie('popup')
}


// common function for fetch

async function fetchData(url, method = "GET") {
  const response = await fetch(url, { method: method });
  if (response.ok) {
    return response;
  } else {
    throw new Error(`API info failed with status ${response.status}`);
  }
}


// Add Project Floating Button

const addProjectFloatingButton = document.querySelector("#add-project-fbtn")
if (addProjectFloatingButton !== null) {
  addProjectFloatingButton.addEventListener("click", () => {
    const github_url = window.prompt("Enter Project Github URL")

    if (github_url !== null) {
      const info = extractOwnerAndRepo(github_url)
      if (info.owner === null || info.repo === null) {
        invalidURLAlert()
      }
      else {
        const repo_api_url = "https://api.github.com/repos/" + info.owner + "/" + info.repo

        // extaract repo info
        fetchData(repo_api_url)
          .then(response => {
            return response.json()
          })
          .then((data) => {
            let info = new Object()
            info["title"] = data["name"]
            info["description"] = data["description"]
            info["github_url"] = data["html_url"]

            return fetchData(data["languages_url"])
              .then(response => {
                return response.json()
              })
              .then((language_data) => {
                info["languages"] = Object.keys(language_data).map((name, index) => {
                  return { id: 0, name };
                });
                return info
              })
              .catch((err) => {
                throw err
              })
          })
          .then((repo_info) => {
            makeRequestAddProject(repo_info)
          })
          .catch((err) => {
            console.log(err)
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
  window.alert("Invalid Github URL!")
}

function fetchErrorAlert() {
  window.alert("Unable to fetch project info. Try later.")
}

function makeRequestAddProject(json_payload) {
  const request_url = "/project/create"

  const form = document.createElement('form');
  form.method = 'GET';
  form.action = request_url;

  const payloadInput = document.createElement('input');
  payloadInput.type = 'hidden';
  payloadInput.name = 'project_detail_json';
  payloadInput.value = JSON.stringify(json_payload);

  form.appendChild(payloadInput);
  document.body.appendChild(form);
  form.submit();
};


// custom 

function customFormSubmission(form) {
  const formData = new FormData(form);
  const urlSearchParams = new URLSearchParams(formData);
  const url = `${form.action}?${urlSearchParams.toString()}`;

  return fetch(url)
    .then(response => {
      if (!response.ok) {
        throw new Error('On customFormSubmission function call | URL : ' + form.action);
      }
      return response.text()
    })
    .catch(error => {
      console.error('Error:', error);
    });
}


// project-search handling

function formSubmissionHandling(form) {
  customFormSubmission(form)
    .then((htmlContent) => {
      const projectList = document.querySelector("#project-list")
      projectList.innerHTML = htmlContent
    })
}

function projectSearchFormSubmit() {
  const projectForm = document.querySelector("#project-search-section")
  formSubmissionHandling(projectForm)
}


// clickable span handling

const searchSelectionDivArray = document.querySelectorAll('.search-selection-container');
for (const searchSelectionDiv of searchSelectionDivArray) {
  searchSelectionDiv.addEventListener('click', function (event) {
    if (event.target.classList.contains('clickable-span')) {

      const clickedSpan = event.target
      const clickedSpanLi = clickedSpan.parentNode
      // const searchResultsDiv = searchSelectionDiv.querySelector(".search-results")
      const selectedDiv = searchSelectionDiv.querySelector(".selected.span-list")

      if (clickedSpanLi.parentNode === selectedDiv) {
        clickedSpanLi.remove()
      }
      else {
        let isSameValueSpanExists = false
        selectedDiv.querySelectorAll('*').forEach(childElement => {
          if (childElement.textContent === clickedSpanLi.textContent) {
            isSameValueSpanExists = true
          }
        });

        if (!isSameValueSpanExists) {
          selectedDiv.appendChild(clickedSpanLi.cloneNode(true));
        }
      }

      const form = document.querySelector("form");
      const payloadId = event.currentTarget.id;
      let payloadInput = form.querySelector(`input[name="${payloadId}"]`)

      if (payloadInput === null) {
        // this should not be case
      }

      let childTextContents = [];
      selectedDiv.querySelectorAll('*').forEach((childElement) => {
        if (childElement.tagName === "SPAN") {
          childTextContents.push({
            id: parseInt(childElement.getAttribute("data-id")),
            name: childElement.textContent.trim()
          })
        }
      });

      payloadInput.value = JSON.stringify(childTextContents);

      const projectForm = document.querySelector("#project-search-section")
      if (projectForm !== null) {
        formSubmissionHandling(projectForm)
      }
    }
  });
}


// sort-by handling

const sortBy = document.querySelector("#sort-by")
if (sortBy !== null) {
  sortBy.addEventListener("change", (event) => {

    projectSearchFormSubmit()

  })
}


// project-search-bar handling
const projectSearchForm = document.querySelector("#project-search-section")
if (projectSearchForm !== null) {

  projectSearchForm.addEventListener('submit', (event) => {
    event.preventDefault();
    projectSearchFormSubmit()
  });
  projectSearchInput = projectSearchForm.querySelector("input[name='project-search']")
  projectSearchInput.addEventListener('input', (event) => {
    projectSearchFormSubmit()
  });

}

// sort-direction handling
const sortDirection = document.querySelector("#sort-direction")
if (sortDirection !== null) {
  sortDirection.addEventListener("click", (event) => {
    if (event.target.tagName === "SPAN") {
      span = event.target
      input = document.querySelector("input[name='sort-direction']")
      if (input.value === "desc") {
        span.classList.add("selected")
        input.value = "asc"
      } else {
        span.classList.remove("selected")
        input.value = "desc"
      }
      projectSearchFormSubmit()
    }
  })
}

// organization-only handling
const organizationOnly = document.querySelector("#organization-only")
if (organizationOnly !== null) {
  organizationOnly.addEventListener("click", (event) => {
    if (event.target.tagName === "SPAN") {
      span = event.target
      input = document.querySelector("input[name='organization-only']")
      if (input.value === "false") {
        span.classList.add("selected")
        input.value = "true"
      } else {
        span.classList.remove("selected")
        input.value = "false"
      }

      projectSearchFormSubmit()

    }
  })
}


// project-tabs handling

const tabs = document.querySelectorAll(".tab")
for (const tab of tabs) {
  tab.addEventListener("click", (event) => {
    inputField = document.querySelector("input[name='tab']")
    inputField.value = event.currentTarget.id;

    formSubmissionHandling(projectSearchForm)

    for (const tab of tabs) {
      tab.classList.remove("selected")
    }
    event.currentTarget.classList.add("selected")
  })
}



// project-card action handling


const projectList = document.querySelector('#project-list');
if (projectList !== null) {
  projectList.addEventListener('click', (event) => {

    // project-delete handling
    if (event.target.classList.contains("project-delete")) {
      projectDelete = event.target
      projectCard = projectDelete.closest(".project-card")
      projectID = projectCard.getAttribute("data-project-id")
      url = "/project/" + projectID

      fetchData(url, "DELETE")
        .then(() => {
          projectCard.remove()
        })
    }


    // project-bookmark handling
    if (event.target.classList.contains("project-toggle-bookmark")) {
      projectToggleBookmark = event.target
      projectCard = projectToggleBookmark.closest(".project-card")
      projectID = projectCard.getAttribute("data-project-id")
      bookmarkCount = projectCard.querySelector(".bookmark-count")
      url = "/project/" + projectID + "/toggle-bookmark"

      const isLogin = document.querySelector("#menu-login")
      if (isLogin === null) {
        fetch(url, { method: "PATCH" })
          .then((response) => {
            if (!response.ok) {
              return response
            }

            if (projectToggleBookmark.innerHTML === "bookmark") {
              bookmarkCount.innerHTML = parseInt(bookmarkCount.innerHTML) - 1
              projectToggleBookmark.innerHTML = "bookmark_border"
            } else {
              bookmarkCount.innerHTML = parseInt(bookmarkCount.innerHTML) + 1
              projectToggleBookmark.innerHTML = "bookmark"
            }
            return null
          }
          ).then((response) => {
            if (response !== null) {
              setCookie(
                "popup",
                "Unfortunately, the project you are trying to access is no longer available."
              )
              window.location.reload()
            }
          })
      } else {
        showPopup("Looking to bookmark this? Register for a free account and never lose track of your favorites.")
      }
    }
  })
}


// addProjectPage

const tagSearch = document.querySelector("input[type='search']")
if (tagSearch !== null) {
  tagSearch.addEventListener("keypress", function (event) {
    if (event.keyCode === 13) {
      event.preventDefault();
    }
  })
}


// userProfilePage

const copyEmail = document.querySelector("#copy-email-btn")
if (copyEmail !== null) {
  copyEmail.addEventListener("click", () => {
    const userEmail = document.querySelector("#user-email")
    navigator.clipboard.writeText(userEmail.textContent)
    showPopup("Email copied")
  })
}

const editProfileBtn = document.querySelector("#user-edit-btn")
if (editProfileBtn !== null) {
  editProfileBtn.addEventListener("click", () => {
    const aboutText = document.querySelector("#upp-about-text")
    const aboutTextarea = document.querySelector("#upp-about-textarea")
    const userUpdateBtn = document.querySelector("#user-update-btn")

    aboutText.style.display = "none"
    editProfileBtn.style.display = "none"
    aboutTextarea.style.display = "block"
    userUpdateBtn.style.display = "block"
  })
}

const userUpdateBtn = document.querySelector("#user-update-btn")
if (userUpdateBtn !== null) {
  userUpdateBtn.addEventListener("click", () => {
    const bio = document.querySelector("#upp-about-textarea")
    const bioValue = bio.value
    const url = window.location.href

    fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "bio": bioValue
      })
    })
      .then(response => {
        if (!response.ok) {
          setCookie(
            "popup",
            "We're unable to process your request at this time. Please try again later."
          )
        }
        window.location.reload()
      })
  })
}