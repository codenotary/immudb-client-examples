function removeTodo(id){
    fetch(`/delete?id=${id}`, {method: "Delete"}).then(res =>{
        if (res.status == 200){
            window.location.pathname = "/";
        }
    })
 }
 
 function updateTodo(id) {
    let input = document.getElementById(id)
    let newdescription = input.value

    fetch(`/update?id=${id}&description=${newdescription}`, {method: "PUT"}).then(res =>{
        if (res.status == 200){
            $("#success-alert").fadeTo(2000, 500).slideUp(500, function() {
                $("#success-alert").slideUp(500);
            });
        } else {
            window.location.pathname = "/"
        }
    })
 }