const limit = document.querySelector("#query_limit")
const orderBy = document.querySelector("#query_order_by")
const orderDirec = document.querySelector("#query_order_direc")
const sortByDirectionSet = document.querySelector("#query_order_direc_set")
const queryButton = document.querySelector("#query_button")

const queryHead = document.querySelector("#query_results_head")
const queryBody = document.querySelector("#query_results_body")

const params = new URLSearchParams(window.location.search)
const tableName = params.get("name")

async function query() {
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

        let columns = Object.keys(resp[0])

        columns = ["id", ...columns.filter(c => c !== "id")]

        columns.forEach(col => {
            const newth = document.createElement("th")
            newth.innerText = col
            queryHead.appendChild(newth)
        });

        resp.forEach(entry => {
            const newTr = document.createElement("tr")
            columns.forEach(col => {
                const newTd = document.createElement("td")
                if (col === "id") {
                    const wrapper = document.createElement("div")
                    wrapper.classList.add("query_id_td")
                    const id = entry[col]
                    const deleteButton = document.createElement("button")
                    deleteButton.classList.add("query_delete_button")
                    deleteButton.innerText = "X"
                    deleteButton.addEventListener("click", async (e) => {
                        e.preventDefault()
                        let url = "/delete_one?"
                        url += "table_name=" + tableName.replaceAll(" ", "%20")
                        url += `&id=${id}`
                        const rawResp = await fetch(url, { method: "DELETE" })
                        if (rawResp.status === 200) {
                            query()
                        }
                    })
                    const idSpan = document.createElement("span")
                    idSpan.innerText = `${id}`
                    wrapper.appendChild(deleteButton)
                    wrapper.appendChild(idSpan)
                    newTd.append(wrapper)
                } else {
                    newTd.innerText = entry[col]
                }
                newTr.appendChild(newTd)
            })
            queryBody.appendChild(newTr)
        })
    }
}

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
    query()
})