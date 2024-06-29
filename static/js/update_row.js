const schema = JSON.parse(document.getElementById("schema").value)

const updateDialog = document.getElementById("update_dialog")
const idInput = document.getElementById("update_row_id")

export function openUpdateDialog(rowData) {

    Object.values(schema).forEach(value => {
        const formField = document.getElementById(`update_${value.InputName}`)
        formField.value = rowData[value.Name]
    })

    idInput.value = rowData.id
    updateDialog.showModal()
}