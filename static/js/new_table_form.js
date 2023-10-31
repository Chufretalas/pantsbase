let newIndex = 0
let columnIndexes = []

const columnIndexesInput = document.querySelector("#column_indexes_input")

const form = document.querySelector("#nt_form")
form.addEventListener("submit", (e => {
    console.log(e)
}))

const ntColumns = document.querySelector(".nt_columns")
const ncButton = document.querySelector("#nt_new_column_button")

ncButton.addEventListener("click", (e) => {
    e.preventDefault()
    columnIndexes.push(newIndex)
    columnIndexesInput.value = columnIndexes.map((number) => `${number}`).join(" ")
    const newColumnElement = document.createElement("li")
    newColumnElement.className = "nt_column"
    newColumnElement.innerHTML = `
                    <span>Column: </span>
                    <div class="nt_column_name">
                        <label for="n${newIndex}" class="nt_column_name_label">Name:</label>
                        <input type="text" name="n${newIndex}" id="n${newIndex}" required>
                    </div>
                    <div class="nt_column_type">
                        <label for="t${newIndex}" class="nt_column_type_label">Type:</label>
                        <select name="t${newIndex}" id="t${newIndex}" class="nt_column_type_input">
                            <option value="INTEGER">INTEGER</option>
                            <option value="REAL">REAL</option>
                            <option value="TEXT">TEXT</option>
                        </select>
                    </div>
                    <button class="remove_new_column" id="remove_new_column_${newIndex}" value=${newIndex} type="button">X</button>
            `
    ntColumns.appendChild(newColumnElement)
    document.querySelector(`#remove_new_column_${newIndex}`).addEventListener("click", removeNewColumn)
    newIndex++
})

function removeNewColumn(e) {
    e.preventDefault()
    console.log("input to delete:", e.target.value) // remember that this value is a string !!!!!!
    for (let child of ntColumns.children) {
        if (child.innerHTML.includes(`t${e.target.value}`)) {
            ntColumns.removeChild(child)
            columnIndexes = columnIndexes.filter(v => {
                // console.log(v, +e.target.value, v !== +e.target.value)
                return v !== +e.target.value
            })
            columnIndexesInput.value = columnIndexes.map((number) => `${number}`).join(" ")
            return
        }
    }
}