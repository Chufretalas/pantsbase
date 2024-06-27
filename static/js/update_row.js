const schema = JSON.parse(document.getElementById("schema").value)

const updateDialog = document.getElementById("update_dialog")
const idInput = document.getElementById("update_row_id")

export function openUpdateDialog(rowData) {

    Object.values(schema).forEach(value => {
        const formField = document.querySelector(`#update_row_form #update_${value.InputName}`)
        formField.value = rowData[value.ColName]
    })


    idInput.value = rowData.id
    updateDialog.showModal()
}