const limit = document.querySelector("#query_limit")
const orderBy = document.querySelector("#query_order_by")
const orderDirec = document.querySelector("#query_order_direc")
const sortByDirectionSet = document.querySelector("#query_order_direc_set")
const queryButton = document.querySelector("#query_button")

const queryHead = document.querySelector("#query_results_head")
const queryBody = document.querySelector("#query_results_body")

const params = new URLSearchParams(window.location.search)
const tableName = params.get("name")

orderBy.addEventListener("change", (e) => {
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
                orderBy: orderBy.value,
                orderDirec: orderDirec.value
            })
        }
    )
    const resp = await rawResp.json()
    if (resp && resp.length > 0) {

        queryHead.innerHTML = ""
        queryBody.innerHTML = ""

        const columns = Object.keys(resp[0])

        columns.forEach(col => {
            const newth = document.createElement("th")
            newth.innerText = col
            queryHead.appendChild(newth)
        });

        resp.forEach(entry => {
            const newTr = document.createElement("tr")
            const values = Object.values(entry)
            values.forEach(v => {
                const newTd = document.createElement("td")
                newTd.innerText = v
                newTr.appendChild(newTd)
            })
            queryBody.appendChild(newTr)
        })
    }
})