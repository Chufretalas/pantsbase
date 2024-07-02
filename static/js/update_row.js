import { query } from "./query.js"

const schema = JSON.parse(document.getElementById("schema").value)
const tableName = document.getElementById("table_name").value

const updateDialog = document.getElementById("update_dialog")
const idInput = document.getElementById("update_row_id")
const newRowForm = document.getElementById("new_row_form")
const updateRowForm = document.getElementById("update_row_form")

//TODO: allow null values for empty fields?
export function openUpdateDialog(rowData) {
    Object.values(schema).forEach(value => {
        const formField = document.getElementById(`update_${value.InputName}`)
        formField.value = rowData[value.Name]
    })

    idInput.value = rowData.id
    updateDialog.showModal()
}

async function submitRowHandler(e, newOrUpdate) {
    e.preventDefault()
    const rowObj = new Object()
    const data = new FormData(e.target)
    data.forEach((v, k) => {
        console.log(k)
        rowObj[schema.find(col => col.InputName === k).Name] = k[0] === "t" ? v : +v
    })
    console.log({ rowObj })

    let updateId
    if (newOrUpdate === "update") {
        updateId = document.getElementById("update_row_id").value
    }

    if (rowObj) {
        const url = newOrUpdate === "update" ? `/api/tables/${tableName}/${updateId}` : `/api/tables/${tableName}`
        const resp = await fetch(url, {
            method: newOrUpdate === "update" ? "PUT" : "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(rowObj)
        })
        if (resp.status === 200) {
            if (newOrUpdate === "update") {
                updateDialog.close()
                query()
            }
            e.target.reset()
        }
    }
}

newRowForm.addEventListener("submit", (e) => submitRowHandler(e, "new"))
updateRowForm.addEventListener("submit", (e) => submitRowHandler(e, "update"))