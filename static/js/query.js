const limit = document.querySelector("#query_limit")
const sortBy = document.querySelector("#query_sort_by")
const sortDirec = document.querySelector("#query_sort_direc")
const sortByDirectionSet = document.querySelector("#query_sort_direc_set")
const queryButton = document.querySelector("#query_button")

const params = new URLSearchParams(window.location.search)
const tableName = params.get("name")

sortBy.addEventListener("change", (e) => {
    if (e.target.value === "NONE 😵") {
        if (!sortByDirectionSet.classList.contains("hidden")) {
            sortByDirectionSet.classList.add("hidden")
        }
    } else {
        sortByDirectionSet.classList.remove("hidden")
    }
})

queryButton.addEventListener("click", async (e) => {
    e.preventDefault()
    console.log(tableName)
    const rawResp = await fetch(`/query`,
        {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                tableName: tableName,
                limit: limit.value < 0 ? 0 : limit.value,
                sortBy: sortBy.value,
                sortDirec: sortDirec.value
            })
        }
    )
    const resp = await rawResp.json()
    console.log(resp)
})