$().ready(function(){

frm = $('#login');
    frm.submit(function (ev) {
        $.ajax({
            type: frm.attr('method'),
            url: frm.attr('action'),
            data: frm.serialize(),
            success: function (data) {
                if (data["status"]=="401") {
					notie.alert('error', data["data"], 3)
				}
				if (data["status"]=="301") {
                    notie.alert('success', data["data"], 0.5)
					window.location.href = data["data"];
				}
            }
        });

        ev.preventDefault();
		// return false
    });
 

}) 
