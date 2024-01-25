const schemaInpt = document.getElementById("schema")

const updateDialog = document.getElementById("update_dialog")
const idInput = document.getElementById("update_row_id")

export function openUpdateDialog(rowData) {

    console.log({rowData})

    // Getting the schema
    const schemaRaw = schemaInpt.value
    const regex = /\{(.*?)\}/g
    const fields = Array.from(schemaRaw.matchAll(regex)).map(match => {
        const split = match[1].split(/\s(INTEGER|TEXT|REAL)\s/g)
        return {
            name: split[0],
            id: split[2]
        }
    })
    // End getting the schema

    console.log(fields)

    fields.forEach(field => {
        const formField = document.querySelector(`#update_row_form #${field.id}`)
        formField.value = rowData[field.name]
    })

    idInput.value = rowData.id
    updateDialog.showModal()
}

//TODO: finish the form and make an endpoint in /update_row to receive the resquest