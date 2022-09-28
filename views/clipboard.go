package views

const Clipboard = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
    <form id="form-clipboard">
        <textarea id="textarea_input" name="input" rows="10" style="width: 100%%" placeholder="Once upon a time..."></textarea>
        <input type="button" id="btn_paste" value="Paste"></input>&nbsp;&nbsp;&nbsp;
        <label for="checkbox_newline">With newline</label>
        <input type="checkbox" id="checkbox_newline">
        &nbsp;&nbsp;&nbsp;
        <input type="submit" class="button" value="Download"></input>
    </form>

    <script>
        const textareaInput = document.getElementById("textarea_input");
        const checkboxNewline = document.getElementById("checkbox_newline");
        
        document.getElementById("btn_paste").addEventListener("click", function() {
            navigator.clipboard.readText()
            .then(clipText => {
                textareaInput.value += clipText;
                
                if (checkboxNewline.checked) {
                    textareaInput.value += "\n";
                }
            });
        });
        
        const formClipboard = document.getElementById("form-clipboard");
        formClipboard.addEventListener("submit", function(e) {
            e.preventDefault();

            fetch("/clipboard", {
                method: 'POST',
                body: new FormData(formClipboard)
            })
            .then(response => {
                if (!!response.ok && response.status == 200) {
                    alert("saved");
                } else {
                    throw 'Not OK';
                }
            })
            .catch(err => {
                alert("failed: " + err);
            });
        });
        
    </script>
</body>
</html>
`
