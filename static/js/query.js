import { openUpdateDialog } from "./update_row.js"

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
    const rawResp = await fetch(`/api/tables/${tableName}?limit=${limit.value < 0 ? 0 : limit.value}&order_by=${orderBy.value}&order_direction=${orderDirec.value}`,
        {
            method: "GET",
            headers: {
                'Accept': 'application/json'
            }
        }
    )
    const resp = await rawResp.json()

    console.log(resp)
    if (resp) {

        queryHead.innerHTML = ""
        queryBody.innerHTML = ""

        if (resp.length === 0) {
            queryBody.innerHTML = `
            <div id="table_is_empty_warning">
            🐻 This table is empty! 🫎
            </div>
            `
            return
        }

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
                    // ------ Creating the delete button ------
                    const deleteButton = document.createElement("button")
                    deleteButton.classList.add("query_delete_button")
                    deleteButton.innerText = "X"
                    deleteButton.addEventListener("click", async (e) => {
                        e.preventDefault() //TODO: add a confirm alert?
                        let url = `/api/tables/${tableName}/${id}`
                        const rawResp = await fetch(url, { method: "DELETE" })
                        if (rawResp.status === 200) {
                            query()
                        }
                    })
                    // ----- creating the actual text
                    const idSpan = document.createElement("span")
                    idSpan.innerText = `${id}`
                    // ------ Creating the edit button ------
                    const editButton = document.createElement("button")
                    editButton.classList.add("query_edit_button")
                    editButton.innerText = "🖋"
                    editButton.addEventListener("click", (e) => {
                        e.preventDefault()
                        openUpdateDialog(entry)
                    })
                    // ----- putting stuff together
                    wrapper.appendChild(deleteButton)
                    wrapper.appendChild(idSpan)
                    wrapper.appendChild(editButton)
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

export const TESTE = "aaaaaa"